package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig Config
	}{
		{
			name: "with environment variables",
			envVars: map[string]string{
				"JWT_SECRET":         "test_secret",
				"JWT_REFRESH_SECRET": "test_refresh_secret",
				"JWT_ACCESS_EXPIRY":  "2h",
				"JWT_REFRESH_EXPIRY": "14d",
			},
			expectedConfig: Config{
				AccessSecret:  "test_secret",
				RefreshSecret: "test_refresh_secret",
				AccessExpiry:  2 * time.Hour,
				RefreshExpiry: 14 * 24 * time.Hour,
			},
		},
		{
			name:    "without environment variables",
			envVars: map[string]string{},
			expectedConfig: Config{
				AccessSecret:  "default_jwt_secret_key",
				RefreshSecret: "default_jwt_secret_key_refresh",
				AccessExpiry:  24 * time.Hour,
				RefreshExpiry: 7 * 24 * time.Hour,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			manager := NewManager()
			require.NotNil(t, manager)
			assert.Equal(t, tt.expectedConfig, manager.config)
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	manager := &Manager{
		config: Config{
			AccessSecret: "test_secret",
			AccessExpiry: time.Hour,
		},
	}

	userID := "user123"
	username := "testuser"
	email := "test@example.com"

	token, err := manager.GenerateAccessToken(userID, username, email)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify token can be parsed
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test_secret"), nil
	})
	require.NoError(t, err)
	require.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*Claims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, AccessToken, claims.TokenType)
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := &Manager{
		config: Config{
			RefreshSecret: "test_refresh_secret",
			RefreshExpiry: 7 * 24 * time.Hour,
		},
	}

	userID := "user123"
	username := "testuser"
	email := "test@example.com"

	token, err := manager.GenerateRefreshToken(userID, username, email)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify token can be parsed
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test_refresh_secret"), nil
	})
	require.NoError(t, err)
	require.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*Claims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, RefreshToken, claims.TokenType)
}

func TestValidateToken(t *testing.T) {
	manager := &Manager{
		config: Config{
			AccessSecret:  "test_secret",
			RefreshSecret: "test_refresh_secret",
		},
	}

	userID := "user123"
	username := "testuser"
	email := "test@example.com"

	// Test access token validation
	accessToken, err := manager.GenerateAccessToken(userID, username, email)
	require.NoError(t, err)

	claims, err := manager.ValidateToken(accessToken, AccessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, AccessToken, claims.TokenType)

	// Test refresh token validation
	refreshToken, err := manager.GenerateRefreshToken(userID, username, email)
	require.NoError(t, err)

	claims, err = manager.ValidateToken(refreshToken, RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, RefreshToken, claims.TokenType)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	manager := &Manager{
		config: Config{
			AccessSecret: "test_secret",
		},
	}

	// Test with invalid token
	_, err := manager.ValidateToken("invalid_token", AccessToken)
	assert.Error(t, err)

	// Test with wrong secret
	wrongManager := &Manager{
		config: Config{
			AccessSecret: "wrong_secret",
		},
	}

	// Generate token with one secret, validate with another
	token, err := manager.GenerateAccessToken("user123", "testuser", "test@example.com")
	require.NoError(t, err)

	_, err = wrongManager.ValidateToken(token, AccessToken)
	assert.Error(t, err)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	manager := &Manager{
		config: Config{
			AccessSecret: "test_secret",
			AccessExpiry: -time.Hour, // Expired in the past
		},
	}

	token, err := manager.GenerateAccessToken("user123", "testuser", "test@example.com")
	require.NoError(t, err)

	_, err = manager.ValidateToken(token, AccessToken)
	assert.Error(t, err)
}

func TestClaims_Valid(t *testing.T) {
	claims := &Claims{
		UserID:    "user123",
		Username:  "testuser",
		Email:     "test@example.com",
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	err := claims.Valid()
	assert.NoError(t, err)
}

func TestClaims_Expired(t *testing.T) {
	claims := &Claims{
		UserID:    "user123",
		Username:  "testuser",
		Email:     "test@example.com",
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	err := claims.Valid()
	assert.Error(t, err)
}
