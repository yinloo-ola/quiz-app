package database

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/yinloo-ola/quiz-app/backend/models" // Import models from the models package
)

// DB is the global database connection pool
var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	var err error
	// Using a simple SQLite file named quiz_app.db in the current directory
	// Ensure the backend executable is run from the 'backend' directory
	dbPath := "quiz_app.db"
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successfully opened")
}

// MigrateDatabase runs the auto-migration for the defined models
func MigrateDatabase() {
	if DB == nil {
		log.Fatal("Database connection is not initialized. Call ConnectDatabase first.")
	}

	err := DB.AutoMigrate(
		&models.AdminUser{},
		&models.Quiz{},
		&models.Question{},
		&models.Choice{},
		&models.ResponderCredential{},
		&models.QuizResponse{},
		&models.Answer{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	fmt.Println("Database migrated successfully")
}
