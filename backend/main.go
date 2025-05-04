package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yinloo-ola/quiz-app/backend/database" // Import the database package
	"github.com/gin-contrib/cors"                      // Import CORS middleware
	"github.com/yinloo-ola/quiz-app/backend/auth"      // Import auth package
	"github.com/yinloo-ola/quiz-app/backend/handlers"  // Import handlers package
	"github.com/yinloo-ola/quiz-app/backend/models"    // Import models package
	"github.com/yinloo-ola/quiz-app/backend/middleware" // Import middleware package
	"gorm.io/gorm"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
}

func main() {
	loadEnv()

	// Load JWT secret keys
	if err := auth.LoadJWTSecret(); err != nil {
		log.Fatalf("FATAL: Failed to load JWT secrets: %v", err)
	}

	// Initialize Database
	database.ConnectDatabase()
	database.MigrateDatabase()

	// Create default admin user if not exists (for development/testing)
	createDefaultAdminUserIfNotExists()

	// Initialize Gin router
	router := gin.Default()

	// Apply CORS middleware
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // Allow any origin for development
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// MaxAge: 12 * time.Hour, // Optional: You can set MaxAge
	}
	router.Use(cors.New(config))

	// Simple root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Quiz App Backend is running!",
		})
	})

	// --- Public Routes ---
	router.POST("/admin/login", handlers.AdminLoginHandler)
	router.POST("/responder/login", handlers.ResponderLoginHandler) // Responder login route

	// --- Responder Routes (Protected by Responder Auth) ---
	responderRoutes := router.Group("/quizzes")
	{
		// Apply Responder Auth Middleware to subsequent responder routes
		responderRoutes.Use(middleware.ResponderAuthMiddleware())
		
		// Get quiz details for responder (without correct answers)
		responderRoutes.GET("/:quiz_id", handlers.GetQuizForResponderHandler)
		
		// Submit quiz answers
		responderRoutes.POST("/:quiz_id/submit", handlers.SubmitQuizHandler)
	}

	// --- Admin Routes (Protected by Admin Auth) ---
	adminRoutes := router.Group("/admin")
	{
		// Apply Auth Middleware to subsequent admin routes
		adminRoutes.Use(middleware.AdminAuthMiddleware())

		// Example protected route (add actual routes later)
		adminRoutes.GET("/me", func(c *gin.Context) {
			adminUserID, exists := c.Get(middleware.AuthPayloadKey)
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID not found in context"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"admin_user_id": adminUserID})
		})

		// --- Quiz Management Routes ---
		adminRoutes.POST("/quizzes", handlers.CreateQuizHandler)
		adminRoutes.GET("/quizzes", handlers.GetQuizzesHandler)
		adminRoutes.GET("/quizzes/:quiz_id", handlers.GetQuizDetailsHandler)
		adminRoutes.PUT("/quizzes/:quiz_id", handlers.UpdateQuizHandler)
		adminRoutes.DELETE("/quizzes/:quiz_id", handlers.DeleteQuizHandler)

		// --- Question Management Routes (within a Quiz) ---
		adminRoutes.POST("/quizzes/:quiz_id/questions", handlers.AddQuestionHandler)

		// --- Responder Credential Management Routes ---
		adminRoutes.POST("/quizzes/:quiz_id/credentials", handlers.GenerateCredentialsHandler)
		adminRoutes.GET("/quizzes/:quiz_id/credentials", handlers.ViewCredentialsHandler)

		// --- Credential Management Routes (Direct) ---
		adminRoutes.DELETE("/credentials/:credential_id", handlers.RevokeCredentialHandler)

		// --- Response Management Routes ---
		adminRoutes.GET("/quizzes/:quiz_id/responses", handlers.ViewResponsesHandler)
		adminRoutes.GET("/responses/:response_id", handlers.ViewResponseDetailsHandler)

		// --- Question Management Routes (Direct) ---
		adminRoutes.PUT("/questions/:question_id", handlers.UpdateQuestionHandler)
		adminRoutes.DELETE("/questions/:question_id", handlers.DeleteQuestionHandler)

		// Other admin routes will go here, potentially under auth middleware
	}

	// Get port from environment variable, default to 8082
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Server starting on port %s", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// createDefaultAdminUserIfNotExists checks if the default 'admin' user exists
// and creates one with a default password if it doesn't.
// This is intended for development/testing purposes only.
func createDefaultAdminUserIfNotExists() {
	var adminUser models.AdminUser
	defaultUsername := "admin"
	defaultPassword := "password" // Use a more secure default if needed, or manage externally

	// Check if the user already exists
	result := database.DB.Where("username = ?", defaultUsername).First(&adminUser)

	if result.Error == gorm.ErrRecordNotFound {
		// User not found, create one
		hashedPassword, err := auth.HashPassword(defaultPassword)
		if err != nil {
			log.Printf("ERROR: Failed to hash default admin password: %v", err)
			return // Don't proceed if hashing fails
		}

		newUser := models.AdminUser{
			Username: defaultUsername,
			Password: hashedPassword,
		}

		createResult := database.DB.Create(&newUser)
		if createResult.Error != nil {
			log.Printf("ERROR: Failed to create default admin user: %v", createResult.Error)
		} else {
			log.Printf("Default admin user '%s' created successfully.", defaultUsername)
		}
	} else if result.Error != nil {
		// Some other database error occurred during the check
		log.Printf("ERROR: Failed to check for default admin user: %v", result.Error)
	} else {
		// User already exists
		log.Printf("Default admin user '%s' already exists.", defaultUsername)
	}
}
