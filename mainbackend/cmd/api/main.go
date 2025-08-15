package main

import (
	"log"
	"mainbackend/internal/auth"
	"mainbackend/internal/platform/database"
	"mainbackend/internal/platform/email"
	"mainbackend/internal/transcription" // YENİ EKLENDİ: Transcription paketi

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Println("Uyarı: .env dosyası bulunamadı.")
	}

	// Veritabanı bağlantısını başlat
	db := database.ConnectDB()

	// E-posta servisini oluştur
	mailer := email.NewGmailMailer()

	// Katmanları (repository, service, handler) oluştur
	userRepository := auth.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository, mailer, "a-very-secret-key")
	authHandler := auth.NewAuthHandler(authService)

	// Gin router'ı kur
	router := gin.Default()

	// CORS ayarları
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Frontend adresiniz
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// API rotalarını tanımla
	api := router.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/verify-email", authHandler.VerifyEmail)
			authRoutes.POST("/resend-code", authHandler.ResendVerificationCode)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
			// DÜZELTİLDİ: Şifre sıfırlama rotası POST metodu ve doğru handler ile değiştirildi.
			authRoutes.POST("/reset-password", authHandler.ResetPassword)
		}

		// YENİ EKLENDİ: Ses metne çevirme (transcription) için yeni endpoint
		api.POST("/transcribe", transcription.TranscribeHandler)
	}

	// Sunucuyu başlat
	log.Println("Sunucu http://localhost:8080 adresinde başlatılıyor...")
	router.Run(":8080")
}
