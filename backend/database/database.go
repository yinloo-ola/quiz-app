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

	// AutoMigrate will create or update tables based on the struct definitions.
	// It will ONLY add missing fields, WON'T delete/change existing ones.
	// --- Manual index migration for ResponderCredential --- 
	// GORM AutoMigrate doesn't reliably handle changing unique indexes.
	// We need to manually drop the old index if it exists.
	migrator := DB.Migrator()
	model := &models.ResponderCredential{}
	
	// Check for common default names GORM might have used for the old index
	oldIndexNames := []string{"idx_responder_credentials_username", "uix_responder_credentials_username"}
	for _, indexName := range oldIndexNames {
		if migrator.HasIndex(model, indexName) {
			log.Printf("Dropping old index '%s' on responder_credentials...", indexName)
			err := migrator.DropIndex(model, indexName)
			if err != nil {
				log.Fatalf("Failed to drop old index %s: %v", indexName, err)
			}
			log.Printf("Successfully dropped old index '%s'.", indexName)
			// Assuming only one of the old names exists, break after dropping
			break 
		}
	}
	// ------------------------------------------------------

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
