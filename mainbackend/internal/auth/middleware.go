package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware, AuthService'i kullanarak bir Gin middleware fonksiyonu oluşturur.
func AuthMiddleware(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header'ı eksik"})
			return
		}

		// Header "Bearer <token>" formatında olmalıdır.
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header'ı formatı hatalı"})
			return
		}
		tokenString := parts[1]

		// authService'in içindeki gizli anahtarı kullanarak token'ı doğrula.
		// Bu kısmı authService'e taşıyarak mantığı merkezileştirebiliriz.
		// Şimdilik basit haliyle burada yapalım.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// İmzalama metodunu kontrol et
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("beklenmedik imzalama metodu: %v", token.Header["alg"])
			}
			// Gizli anahtarı authService'ten almamız lazım, bu yüzden servisi public yapmalıyız
			// veya JWT doğrulamayı servise taşımalıyız. Şimdilik servisi güncelleyelim.
			// TODO: Bu mantığı AuthService'e taşı.
			return authService.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz veya süresi dolmuş token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token claims"})
			return
		}

		userID, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz kullanıcı ID'si"})
			return
		}

		// Kullanıcı ID'sini isteğin context'ine ekle.
		// Böylece handler fonksiyonları bu bilgiye erişebilir.
		c.Set("userID", uint(userID))

		// Bir sonraki adıma geç
		c.Next()
	}
}

// Bu middleware'in çalışması için auth_service.go'da küçük bir ekleme yapmamız gerekiyor:
// AuthService interface'ine ve authService struct'ına GetJWTSecret() []byte metodu eklenmeli.
