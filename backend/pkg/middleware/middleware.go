package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/studyplatform/backend/pkg/auth"
	"github.com/studyplatform/backend/pkg/logger"
)

// Middleware holds all middleware handlers
type Middleware struct {
	jwtManager *auth.Manager
}

// NewMiddleware creates a new middleware instance
func NewMiddleware() *Middleware {
	return &Middleware{
		jwtManager: auth.NewManager(),
	}
}

// CORS returns a CORS middleware configuration
func CORS() gin.HandlerFunc {
	// Get allowed origins from environment
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "") {
		allowedOrigins = []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:3000", "http://127.0.0.1:5173"} // Default values
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// Logger logs request and response details
func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		logger.Info(
			"API Request",
			logger.Field("method", c.Request.Method),
			logger.Field("path", path),
			logger.Field("query", raw),
			logger.Field("status", c.Writer.Status()),
			logger.Field("latency", latency.String()),
			logger.Field("client_ip", c.ClientIP()),
			logger.Field("user_agent", c.Request.UserAgent()),
		)
	}
}

// Auth verifies the JWT token in the request
func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		// Extract the token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			logger.Warn("JWT validation failed", logger.Field("error", err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// Recovery recovers from any panics and logs the error
func (m *Middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(
					"Panic recovered",
					logger.Field("error", err),
					logger.Field("path", c.Request.URL.Path),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// RateLimiter limits request rates based on client IP
func (m *Middleware) RateLimiter() gin.HandlerFunc {
	// This is a simplified version. A production version would use Redis or another distributed store
	// to track request counts across multiple instances.
	return func(c *gin.Context) {
		// Implementation would go here
		c.Next()
	}
}

// RequestID adds a unique ID to each request
func (m *Middleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a unique request ID
		requestID := generateRequestID()
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// generateRequestID creates a unique ID for a request
func generateRequestID() string {
	// Simple implementation - in production, use a more robust method
	return time.Now().Format("20060102150405") + "-" + strings.ReplaceAll(time.Now().String(), " ", "-")
}
