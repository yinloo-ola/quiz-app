package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yinloo-ola/quiz-app/backend/auth"
	"github.com/yinloo-ola/quiz-app/backend/database"
	"github.com/yinloo-ola/quiz-app/backend/middleware"
	"github.com/yinloo-ola/quiz-app/backend/models"
)

// --- Structs ---

// GenerateCredentialsRequest defines the request body for generating credentials.
// Count is removed as we always generate one.
type GenerateCredentialsRequest struct {
	Username    *string `json:"username"`     // Optional username
	ExpiryHours *int    `json:"expiry_hours"` // Optional expiry duration in hours
}

// GenerateCredentialsResponse defines the response containing the generated plain-text credentials.
type GenerateCredentialsResponse struct {
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	ExpiresAt    time.Time `json:"expires_at"`
	CredentialID uint      `json:"credential_id"`
}

// CredentialView represents the data returned for a single credential in the ViewCredentials list.
// Note: It omits the Token (password hash) for security.
type CredentialView struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	ExpiresAt time.Time  `json:"expiresAt"` // Updated to non-pointer to match model
	CreatedAt time.Time  `json:"createdAt"` // Match frontend type
	Used      bool       `json:"used"`
	UsedAt    *time.Time `json:"usedAt"` // Add UsedAt, match frontend type
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

	// 3. Bind Request
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

	// 5. Check for existing credentials for this quiz (enforcing 1-to-1 relationship)
	var existingCredential models.ResponderCredential
	existingQuery := tx.Where("quiz_id = ?", quizID).First(&existingCredential)

	// Generate new credential details
	username := ""
	if req.Username != nil && *req.Username != "" {
		username = *req.Username
	} else {
		// Generate random if not provided
		username = fmt.Sprintf("quiz%d_user%d", quizID, time.Now().UnixNano()%10000)
	}

	plainTextPassword := generateRandomString(12) // Plain-text password

	// Set an expiration date, using a default of 24 hours if none is provided
	expiry := time.Now().Add(24 * time.Hour) // Default: 24 hours from now
	if req.ExpiryHours != nil && *req.ExpiryHours > 0 {
		expiry = time.Now().Add(time.Duration(*req.ExpiryHours) * time.Hour)
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(plainTextPassword)
	if err != nil {
		tx.Rollback()
		log.Printf("Error hashing password for quiz %d: %v", quizID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var newCredential models.ResponderCredential

	// 6. Create or update ResponderCredential record
	if existingQuery.Error == nil {
		// Credential exists for this quiz - update it
		existingCredential.Username = username
		existingCredential.PasswordHash = hashedPassword
		existingCredential.ExpiresAt = expiry
		existingCredential.Used = false // Reset usage status
		existingCredential.UsedAt = nil // Clear usage timestamp

		if err := tx.Save(&existingCredential).Error; err != nil {
			tx.Rollback()
			log.Printf("Error updating credential for quiz %d: %v", quizID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update credential"})
			return
		}
		newCredential = existingCredential
	} else if existingQuery.Error == gorm.ErrRecordNotFound {
		// No credential exists - create new one
		newCredential = models.ResponderCredential{
			QuizID:       quizID,
			Username:     username,
			PasswordHash: hashedPassword, // Store the hash
			ExpiresAt:    expiry,         // Store expiry time
		}
		if err := tx.Create(&newCredential).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating credential for quiz %d: %v", quizID, err)
			
			// Check if this is a uniqueness constraint error
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				if strings.Contains(err.Error(), "username") {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists. Please choose a different username."})
					return
				}
			}
			
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credential"})
			return
		}
	} else {
		// Some other database error
		tx.Rollback()
		log.Printf("Error checking for existing credentials for quiz %d: %v", quizID, existingQuery.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking existing credentials"})
		return
	}

	// --- Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 7. Respond with plain-text credentials
	resp := GenerateCredentialsResponse{
		Username:     username,
		Password:     plainTextPassword, // Send plain text back to admin
		ExpiresAt:    expiry,            // Send back the calculated expiry time
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

	// 4. Fetch Credentials from DB
	var credentials []models.ResponderCredential
	if err := tx.Where("quiz_id = ?", quizID).Order("created_at desc").Find(&credentials).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch credentials: " + err.Error()})
		return
	}

	// --- Commit Transaction (Read-only operation, but good practice for consistency) ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 5. Map to Response View Model (omitting password hash)
	responseCredentials := make([]CredentialView, len(credentials))
	for i, cred := range credentials {
		responseCredentials[i] = CredentialView{
			ID:        cred.ID,
			Username:  cred.Username,
			ExpiresAt: cred.ExpiresAt, // Already a pointer
			CreatedAt: cred.CreatedAt,
			Used:      cred.Used,
			UsedAt:    cred.UsedAt, // Add UsedAt
		}
	}

	c.JSON(http.StatusOK, responseCredentials)
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

	// 4. Hard Delete the Credential (instead of soft delete)
	// Use Unscoped() to bypass GORM's soft delete and perform a real DELETE
	if err := tx.Unscoped().Delete(&credential).Error; err != nil {
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
