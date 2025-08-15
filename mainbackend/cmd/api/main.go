package main

import (
	"fmt" // Bu satırı ekledim, sunucu başlangıç mesajı için.
	"github.com/gin-contrib/cors"
	"log"
	"mainbackend/internal/auth"
	"mainbackend/internal/platform/database"
	"mainbackend/internal/platform/email"
	"mainbackend/internal/transcription" // YENİ: Transcription paketini import et

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
	userRepository := auth.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository, mailer, "a-very-secret-key")
	authHandler := auth.NewAuthHandler(authService)

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
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
		protectedRoutes.Use(auth.AuthMiddleware(authService)) // Güvenlik görevlisini (middleware) bu gruba atıyoruz
		{
			// /transcribe rotasını bu korumalı grubun içine taşıyoruz.
			protectedRoutes.POST("/transcribe", transcription.TranscribeHandler)
			// Gelecekte "dosyalarımı listele" gibi diğer korumalı rotalar da buraya eklenecek.
		}
	}

	fmt.Println("Sunucu http://localhost:8080 adresinde başlatıldı. İstek bekleniyor...")
	router.Run(":8080")
}
