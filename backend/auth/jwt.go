package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecretKey        []byte
	responderJwtSecretKey []byte
	jwtIssuer           = "quiz-app-backend"
	jwtAdminDuration    = time.Hour * 24 * 7 // 1 week for admins
	jwtResponderDuration = time.Hour * 2      // 2 hours for responders
)

// LoadJWTSecret loads the JWT secret key from environment variables.
func LoadJWTSecret() error {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return errors.New("JWT_SECRET_KEY environment variable not set")
	}
	jwtSecretKey = []byte(secret)

	responderSecret := os.Getenv("RESPONDER_JWT_SECRET_KEY")
	if responderSecret == "" {
		return errors.New("RESPONDER_JWT_SECRET_KEY environment variable not set")
	}
	responderJwtSecretKey = []byte(responderSecret)

	return nil
}

// AdminClaims defines the structure for admin JWT claims.
type AdminClaims struct {
	AdminUserID uint `json:"admin_user_id"`
	jwt.RegisteredClaims
}

// ResponderClaims defines the structure for responder JWT claims.
type ResponderClaims struct {
	ResponderCredentialID uint `json:"responder_credential_id"`
	QuizID                uint `json:"quiz_id"`
	jwt.RegisteredClaims
}

// GenerateAdminJWT generates a JWT for an admin user.
func GenerateAdminJWT(adminUserID uint) (string, error) {
	if len(jwtSecretKey) == 0 {
		return "", fmt.Errorf("JWT secret key not loaded")
	}

	expirationTime := time.Now().Add(jwtAdminDuration)

	claims := &AdminClaims{
		AdminUserID: adminUserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)

	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

// ValidateAdminJWT validates an admin JWT string.
func ValidateAdminJWT(tokenString string) (*AdminClaims, error) {
	if len(jwtSecretKey) == 0 {
		return nil, fmt.Errorf("JWT secret key not loaded")
	}

	claims := &AdminClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse or validate JWT: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid JWT")
	}

	return claims, nil
}

// GenerateResponderJWT generates a JWT for a quiz responder.
func GenerateResponderJWT(credentialID uint, quizID uint) (string, error) {
	expirationTime := time.Now().Add(jwtResponderDuration)
	claims := &ResponderClaims{
		ResponderCredentialID: credentialID,
		QuizID:                quizID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(responderJwtSecretKey)
	if err != nil {
		log.Printf("Error signing responder JWT: %v", err)
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateResponderJWT validates a responder JWT string.
func ValidateResponderJWT(tokenString string) (*ResponderClaims, error) {
	claims := &ResponderClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return responderJwtSecretKey, nil
	})

	if err != nil {
		log.Printf("Error parsing/validating responder JWT: %v", err)
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
