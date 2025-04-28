package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/quiz-app/backend/auth" // Import our auth package
)

const (
	authorizationHeaderKey = "Authorization"
	authorizationTypeBearer = "Bearer"
	AuthPayloadKey        = "adminUserID" // Key to store user ID in context
)

// AdminAuthMiddleware creates a gin middleware for admin JWT authentication.
func AdminAuthMiddleware() gin.HandlerFunc {
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
		claims, err := auth.ValidateAdminJWT(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token: " + err.Error()})
			return
		}

		// Set the admin user ID in the context for downstream handlers
		c.Set(AuthPayloadKey, claims.AdminUserID)

		// Continue to the next handler
		c.Next()
	}
}
