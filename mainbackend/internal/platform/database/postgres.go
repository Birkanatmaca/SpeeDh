package database

import (
	"fmt"
	"log"
	"mainbackend/internal/model"
	"os" // We need the 'os' package to read environment variables

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// The hardcoded DSN constant has been removed.

// ConnectDB connects to the PostgreSQL database using credentials from environment variables.
func ConnectDB() *gorm.DB {
	// 1. Read the database configuration from environment variables.
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// 2. Build the DSN (Data Source Name) string dynamically.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	// 3. Open the database connection with the dynamically created DSN.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	fmt.Println("Successfully connected to the database")

	// 4. AutoMigrate the database schema.
	fmt.Println("Database is being migrated")
	err = db.AutoMigrate(&model.User{}, &model.Transcription{}) // Add other models here later, e.g., &model.Transcription{}
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	fmt.Println("The database has been successfully migrated.")

	return db
}
