package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/quiz-app/backend/auth"
	"github.com/yinloo-ola/quiz-app/backend/database"
	"github.com/yinloo-ola/quiz-app/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ResponderLoginInput defines the structure for the responder login request body.
type ResponderLoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ResponderLoginResponse defines the structure for the responder login response body.
type ResponderLoginResponse struct {
	Token string `json:"token"`
}

// ResponderLoginHandler handles the login process for quiz responders.
func ResponderLoginHandler(c *gin.Context) {
	var input ResponderLoginInput

	// 1. Bind and validate input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// --- Transaction Start ---
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction: " + tx.Error.Error()})
		return
	}

	// 2. Find the responder credential by username
	// Since we have a 1-to-1 relationship between quiz and credential,
	// we can uniquely identify a credential by username alone
	var credential models.ResponderCredential
	query := tx.Where("username = ?", input.Username).First(&credential)

	if query.Error != nil {
		tx.Rollback()
		if query.Error == gorm.ErrRecordNotFound {
			// Generic error to avoid revealing if username exists
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			log.Printf("Database error finding responder credential for %s: %v", input.Username, query.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during login"})
		}
		return
	}

	// 3. Compare the provided password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(credential.PasswordHash), []byte(input.Password))
	if err != nil {
		// Passwords don't match
		tx.Rollback()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 4. Check if the credential has already been used
	if credential.Used {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{"error": "Credential has already been used"})
		return
	}

	// 5. Check if the credential has expired
	if credential.ExpiresAt.Before(time.Now()) {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{"error": "Credential has expired"})
		return
	}

	// --- Commit Transaction (Read-only, but good practice before generating JWT) ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 6. Generate Responder JWT
	tokenString, err := auth.GenerateResponderJWT(credential.ID, credential.QuizID)
	if err != nil {
		log.Printf("Error generating JWT for responder credential %d: %v", credential.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	// 7. Return the token
	c.JSON(http.StatusOK, ResponderLoginResponse{Token: tokenString})
}
