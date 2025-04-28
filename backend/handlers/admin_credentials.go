package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yinloo-ola/quiz-app/backend/auth"
	"github.com/yinloo-ola/quiz-app/backend/database"
	"github.com/yinloo-ola/quiz-app/backend/middleware"
	"github.com/yinloo-ola/quiz-app/backend/models"
)

// --- Structs ---

// GenerateCredentialsRequest defines the (currently empty) request body for generating credentials.
// We might add options like specifying expiry duration or number of credentials later.
type GenerateCredentialsRequest struct {
	// Count int `json:"count" binding:"required,min=1"` // Example: If we want to generate multiple at once
	// DurationHours int `json:"duration_hours" binding:"required,min=1"` // Example: Custom duration
}

// GenerateCredentialsResponse defines the response containing the generated plain-text credentials.
type GenerateCredentialsResponse struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	ExpiresAt time.Time `json:"expires_at"`
	CredentialID uint   `json:"credential_id"`
}

// CredentialView represents the data returned for a single credential in the ViewCredentials list.
// Note: It omits the Token (password hash) for security.
type CredentialView struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	Used      bool       `json:"used"`
}

// --- Utility ---

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// --- Handler ---

// GenerateCredentialsHandler handles generating a new set of responder credentials for a quiz.
func GenerateCredentialsHandler(c *gin.Context) {
	// 1. Get Quiz ID
	quizIDStr := c.Param("quiz_id")
	var quizID uint
	_, err := fmt.Sscan(quizIDStr, &quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// 2. Get Admin User ID
	adminUserIDAny, exists := c.Get(middleware.AuthPayloadKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID not found in context"})
		return
	}
	adminUserID, ok := adminUserIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID has incorrect type in context"})
		return
	}

	// 3. Bind Request (currently empty, but good practice)
	var req GenerateCredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Handle potential future binding errors even if empty now
		if err.Error() != "EOF" { // Ignore EOF for empty body
 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
 			return
 		}
	}

	// --- Transaction Start ---
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction: " + tx.Error.Error()})
		return
	}

	// 4. Verify Quiz exists and Admin owns it
	var quiz models.Quiz
	if err := tx.First(&quiz, quizID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding quiz: " + err.Error()})
		}
		return
	}
	if quiz.AdminUserID != adminUserID {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"}) // 404, not 403
		return
	}

	// 5. Generate Credentials
	// Simple generation logic for now
	username := fmt.Sprintf("quiz%d_user%d", quizID, time.Now().UnixNano()%10000)
	plaintextPassword := generateRandomString(12) // Plain-text password
	expiresAt := time.Now().Add(24 * time.Hour) // Default 24-hour expiry

	// Hash the password
	hashedPassword, err := auth.HashPassword(plaintextPassword)
	if err != nil {
		tx.Rollback()
		log.Printf("Error hashing password for quiz %d: %v", quizID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 6. Create ResponderCredential record
	newCredential := models.ResponderCredential{
		QuizID:       quizID,
		Username:     username,
		PasswordHash: hashedPassword, // Store the hash
		ExpiresAt:    &expiresAt, // Store expiry time
	}
	if err := tx.Create(&newCredential).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating credential for quiz %d: %v", quizID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credential"})
		return
	}

	// --- Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 7. Respond with plain-text credentials
	resp := GenerateCredentialsResponse{
		Username:    username,
		Password:    plaintextPassword, // Send plain text back to admin
		ExpiresAt:   expiresAt,
		CredentialID: newCredential.ID,
	}

	c.JSON(http.StatusCreated, resp)
}

// ViewCredentialsHandler handles fetching active credentials for a specific quiz.
func ViewCredentialsHandler(c *gin.Context) {
	// 1. Get Quiz ID
	quizIDStr := c.Param("quiz_id")
	var quizID uint
	_, err := fmt.Sscan(quizIDStr, &quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// 2. Get Admin User ID
	adminUserIDAny, exists := c.Get(middleware.AuthPayloadKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID not found in context"})
		return
	}
	adminUserID, ok := adminUserIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID has incorrect type in context"})
		return
	}

	// --- Transaction Start ---
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction: " + tx.Error.Error()})
		return
	}

	// 3. Verify Quiz exists and Admin owns it
	var quiz models.Quiz
	if err := tx.First(&quiz, quizID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding quiz: " + err.Error()})
		}
		return
	}
	if quiz.AdminUserID != adminUserID {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"}) // 404, not 403
		return
	}

	// 4. Fetch Active Credentials
	var credentials []models.ResponderCredential
	now := time.Now()
	query := tx.Where("quiz_id = ? AND used = ? AND (expires_at IS NULL OR expires_at > ?)", quizID, false, now).Order("created_at desc").Find(&credentials)
	if query.Error != nil {
		tx.Rollback()
		log.Printf("Error fetching credentials for quiz %d: %v", quizID, query.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch credentials"})
		return
	}

	// --- Commit Transaction (Read-only, but good practice to manage tx lifecycle) ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 5. Prepare Response View (excluding token)
	responseView := make([]CredentialView, len(credentials))
	for i, cred := range credentials {
		responseView[i] = CredentialView{
			ID:        cred.ID,
			Username:  cred.Username,
			ExpiresAt: cred.ExpiresAt,
			CreatedAt: cred.CreatedAt,
			Used:      cred.Used,
		}
	}

	c.JSON(http.StatusOK, responseView)
}

// RevokeCredentialHandler handles deleting (soft deleting) a specific responder credential.
func RevokeCredentialHandler(c *gin.Context) {
	// 1. Get Credential ID
	credentialIDStr := c.Param("credential_id")
	var credentialID uint
	_, err := fmt.Sscan(credentialIDStr, &credentialID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credential ID format"})
		return
	}

	// 2. Get Admin User ID
	adminUserIDAny, exists := c.Get(middleware.AuthPayloadKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID not found in context"})
		return
	}
	adminUserID, ok := adminUserIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID has incorrect type in context"})
		return
	}

	// --- Transaction Start ---
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction: " + tx.Error.Error()})
		return
	}

	// 3. Find Credential and Verify Ownership via Quiz
	var credential models.ResponderCredential
	// Preload Quiz to get AdminUserID for ownership check
	if err := tx.Preload("Quiz").First(&credential, credentialID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Credential not found"})
		} else {
			log.Printf("Error finding credential %d: %v", credentialID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding credential"})
		}
		return
	}

	// Check if the AdminUser associated with the Quiz owns this credential
	if credential.Quiz.AdminUserID != adminUserID {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Credential not found"}) // 404, not 403, to avoid revealing existence
		return
	}

	// 4. Delete (Soft Delete) the Credential
	if err := tx.Delete(&credential).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting credential %d: %v", credentialID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete credential"})
		return
	}

	// --- Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 5. Respond with No Content
	c.Status(http.StatusNoContent)
}
