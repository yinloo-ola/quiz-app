package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/quiz-app/backend/auth"
	"github.com/yinloo-ola/quiz-app/backend/database"
	"github.com/yinloo-ola/quiz-app/backend/models"
)

const (
	ResponderAuthPayloadKey = "responderCredentialID" // Key to store credential ID in context
	ResponderQuizIDKey      = "responderQuizID"       // Key to store quiz ID in context
)

// ResponderAuthMiddleware creates a gin middleware for responder JWT authentication.
func ResponderAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeaderKey)

		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is not provided"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != strings.ToLower(authorizationTypeBearer) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unsupported authorization type: " + authType})
			return
		}

		accessToken := fields[1]
		claims, err := auth.ValidateResponderJWT(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token: " + err.Error()})
			return
		}

		// Verify that the credential still exists and is valid
		var credential models.ResponderCredential
		result := database.DB.First(&credential, claims.ResponderCredentialID)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Credential no longer valid"})
			return
		}

		// Check if the credential has been marked as used
		if credential.Used {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "This quiz has already been submitted"})
			return
		}

		// Set the responder credential ID and quiz ID in the context for downstream handlers
		c.Set(ResponderAuthPayloadKey, claims.ResponderCredentialID)
		c.Set(ResponderQuizIDKey, claims.QuizID)

		// Continue to the next handler
		c.Next()
	}
}
