package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"log"
	"mainbackend/internal/auth"
	"mainbackend/internal/platform/ai"
	"mainbackend/internal/platform/database"
	"mainbackend/internal/platform/email"
	"mainbackend/internal/transcription"

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

	aiClient := ai.NewOllamaClient()

	// Auth (Kimlik Doğrulama) servisleri
	userRepository := auth.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository, mailer, "a-very-secret-key")
	authHandler := auth.NewAuthHandler(authService)

	// --- YENİ EKLENEN KISIM ---
	// Transcription (Metne Dönüştürme) servisleri
	transcriptionRepo := transcription.NewTranscriptionRepository(db)
	transcriptionService := transcription.NewTranscriptionService(transcriptionRepo, aiClient)
	transcriptionHandler := transcription.NewTranscribeHandler(transcriptionService) // Handler'ı servisle başlat
	// --- BİTİŞ ---

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"} // Authorization başlığına izin ver
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
			// Handler'ı doğrudan kullanmak yerine oluşturduğumuz handler örneğini kullanıyoruz
			protectedRoutes.POST("/transcribe", transcriptionHandler.Transcribe)
			protectedRoutes.GET("/transcripts", transcriptionHandler.GetTranscripts)
			protectedRoutes.GET("/transcripts/:id/audio", transcriptionHandler.GetAudioFile)
			protectedRoutes.DELETE("/transcripts/:id", transcriptionHandler.DeleteTranscript)
			protectedRoutes.GET("/users/me", authHandler.GetUserProfile)
			protectedRoutes.PUT("/users/me/password", authHandler.ChangePassword)
		}
	}

	fmt.Println("Sunucu http://localhost:8080 adresinde başlatıldı. İstek bekleniyor...")
	router.Run(":8080")
}
