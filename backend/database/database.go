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

	log.Println("Running database migrations...")

	// --- Drop responder_credentials table to ensure clean schema (Development Only!) ---
	log.Println("Attempting to drop responder_credentials table before migration...")
	if err := DB.Migrator().DropTable(&models.ResponderCredential{}); err != nil {
		log.Printf("Warning: Failed to drop responder_credentials table (may not exist yet): %v", err)
	} else {
		log.Println("Successfully dropped responder_credentials table.")
	}
	// --- End Drop Table ---

	// AutoMigrate will create or update tables based on the struct definitions.
	// It will ONLY add missing fields, WON'T delete/change existing ones.
	err := DB.AutoMigrate(
		&models.AdminUser{},
		&models.Quiz{},
		&models.Question{},
		&models.Choice{},
		&models.ResponderCredential{},
		// models.Response{} was removed and replaced by QuizResponse
		&models.QuizResponse{},
		&models.Answer{},
	)

	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	fmt.Println("Database migrated successfully")
}
