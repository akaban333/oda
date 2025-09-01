package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenType defines the type of JWT token
type TokenType string

const (
	// AccessToken is a short-lived token used for authentication
	AccessToken TokenType = "access"
	// RefreshToken is a long-lived token used to obtain new access tokens
	RefreshToken TokenType = "refresh"
)

// Claims represents the JWT claims
type Claims struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// Config holds JWT configuration settings
type Config struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// Manager handles JWT token operations
type Manager struct {
	config Config
}

// NewManager creates a new JWT manager with environment variables
func NewManager() *Manager {
	accessSecret := os.Getenv("JWT_SECRET")
	if accessSecret == "" {
		accessSecret = "default_jwt_secret_key"
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = accessSecret + "_refresh"
	}

	accessExpiryStr := os.Getenv("JWT_ACCESS_EXPIRY")
	accessExpiry := time.Hour * 24 // Default to 24 hours
	if accessExpiryStr != "" {
		if duration, err := time.ParseDuration(accessExpiryStr); err == nil {
			accessExpiry = duration
		}
	}

	refreshExpiryStr := os.Getenv("JWT_REFRESH_EXPIRY")
	refreshExpiry := time.Hour * 24 * 7 // Default to 7 days
	if refreshExpiryStr != "" {
		if duration, err := time.ParseDuration(refreshExpiryStr); err == nil {
			refreshExpiry = duration
		}
	}

	return &Manager{
		config: Config{
			AccessSecret:  accessSecret,
			RefreshSecret: refreshSecret,
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
	}
}

// GenerateAccessToken creates a new access token
func (m *Manager) GenerateAccessToken(userID, username, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.AccessSecret))
}

// GenerateRefreshToken creates a new refresh token
func (m *Manager) GenerateRefreshToken(userID, username, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		TokenType: RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.RefreshSecret))
}

// ValidateAccessToken validates an access token and returns the claims
func (m *Manager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.config.AccessSecret, AccessToken)
}

// ValidateRefreshToken validates a refresh token and returns the claims
func (m *Manager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.config.RefreshSecret, RefreshToken)
}

// validateToken validates a token against its secret and expected type
func (m *Manager) validateToken(tokenString, secret string, expectedType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.TokenType != expectedType {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
