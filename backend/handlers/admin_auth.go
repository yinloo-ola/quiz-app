package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yinloo-ola/quiz-app/backend/auth"      // Import auth package
	"github.com/yinloo-ola/quiz-app/backend/database"  // Import database package
	"github.com/yinloo-ola/quiz-app/backend/models"    // Import models package
)

// LoginRequest defines the structure for the admin login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse defines the structure for the successful login response
type LoginResponse struct {
	Token string `json:"token"`
}

// AdminLoginHandler handles the admin login request
func AdminLoginHandler(c *gin.Context) {
	var req LoginRequest
	var adminUser models.AdminUser

	// Bind JSON request body to the LoginRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Find the admin user by username
	result := database.DB.Where("username = ?", req.Username).First(&adminUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + result.Error.Error()})
		}
		return
	}

	// Check the password
	if !auth.CheckPasswordHash(req.Password, adminUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT
	tokenString, err := auth.GenerateAdminJWT(adminUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
		return
	}

	// Return the token
	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}
