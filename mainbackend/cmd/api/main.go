package main

import (
	"github.com/gin-contrib/cors"
	"log"
	"mainbackend/internal/auth"
	"mainbackend/internal/platform/database"
	"mainbackend/internal/platform/email"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found.")
	}

	// Initialize the database connection
	db := database.ConnectDB()

	// Create the email service instance
	mailer := email.NewGmailMailer()

	// Create the instances of our layers, passing the dependencies
	userRepository := auth.NewUserRepository(db)

	// CORRECTION: We pass the 'mailer' instance to the auth service here.
	authService := auth.NewAuthService(userRepository, mailer, "a-very-secret-key")

	authHandler := auth.NewAuthHandler(authService)

	// Set up the Gin router and define the routes
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Vite'ın verdiği adresi buraya yaz.
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	api := router.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/verify-email", authHandler.VerifyEmail) // YENİ ROTA
			authRoutes.POST("/resend-code", authHandler.ResendVerificationCode)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
			authRoutes.GET("/confirm-reset", authHandler.ConfirmPasswordReset) // DİKKAT: Metot GET OLDU
		}
	}

	// Start the server
	router.Run(":8080")
}
