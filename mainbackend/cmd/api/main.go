// mainbackend/cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"mainbackend/internal/auth"
	"mainbackend/internal/platform/ai" // ai paketini import ediyoruz
	"mainbackend/internal/platform/database"
	"mainbackend/internal/platform/email"
	"mainbackend/internal/transcription"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found.")
	}

	db := database.ConnectDB()
	mailer := email.NewGmailMailer()
	aiClient := ai.NewOllamaClient() // YENİ: AI İstemcisini oluşturuyoruz

	// Auth servisleri
	userRepository := auth.NewUserRepository(db)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "a-very-secret-key" // Fallback
	}
	authService := auth.NewAuthService(userRepository, mailer, jwtSecret)
	authHandler := auth.NewAuthHandler(authService)

	// Transcription servisleri
	transcriptionRepo := transcription.NewTranscriptionRepository(db)
	// Servisi oluştururken artık aiClient'ı da veriyoruz
	transcriptionService := transcription.NewTranscriptionService(transcriptionRepo, aiClient)
	transcriptionHandler := transcription.NewTranscribeHandler(transcriptionService)

	router := gin.Default()

	// ... (CORS ve router kodunuzun geri kalanı aynı kalacak) ...
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	api := router.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/verify-email", authHandler.VerifyEmail)
			authRoutes.POST("/resend-code", authHandler.ResendVerificationCode)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
			authRoutes.POST("/reset-password", authHandler.ResetPassword)

		}
		protectedRoutes := api.Group("/")
		protectedRoutes.Use(auth.AuthMiddleware(authService))
		{
			protectedRoutes.POST("/transcribe", transcriptionHandler.Transcribe)
			protectedRoutes.GET("/transcripts", transcriptionHandler.GetTranscripts)
			protectedRoutes.GET("/transcripts/:id/audio", transcriptionHandler.GetAudioFile)
			protectedRoutes.DELETE("/transcripts/:id", transcriptionHandler.DeleteTranscript)
			protectedRoutes.GET("/users/me", authHandler.GetUserProfile)
			protectedRoutes.PUT("/users/me/password", authHandler.ChangePassword)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Sunucu http://localhost:%s adresinde başlatıldı. İstek bekleniyor...\n", port)
	router.Run(":" + port)
}
