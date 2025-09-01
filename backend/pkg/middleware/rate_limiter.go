package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int           `json:"requests_per_minute"`
	BurstSize         int           `json:"burst_size"`
	WindowSize        time.Duration `json:"window_size"`
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		RequestsPerMinute: 300, // Temporarily increased from 60 to 300
		BurstSize:         50,  // Temporarily increased from 10 to 50
		WindowSize:        time.Minute,
	}
}

// RateLimiter implements IP-based rate limiting
type RateLimiter struct {
	config        *RateLimitConfig
	clients       map[string]*ClientTracker
	mu            sync.RWMutex
	cleanupTicker *time.Ticker
}

// ClientTracker tracks requests for a specific client
type ClientTracker struct {
	Requests    []time.Time
	LastSeen    time.Time
	Blocked     bool
	BlockedAt   time.Time
	BlockReason string
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(config *RateLimitConfig) *RateLimiter {
	if config == nil {
		config = DefaultRateLimitConfig()
	}

	rl := &RateLimiter{
		config:  config,
		clients: make(map[string]*ClientTracker),
	}

	// Start cleanup routine
	rl.cleanupTicker = time.NewTicker(5 * time.Minute)
	go rl.cleanupRoutine()

	return rl
}

// RateLimit returns a Gin middleware for rate limiting
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TEMPORARILY DISABLED - Allow all requests
		c.Next()
		return

		// Original rate limiting code (commented out)
		/*
			clientIP := rl.getClientIP(c)

			if !rl.allowRequest(clientIP) {
				logger.Warn("Rate limit exceeded",
					logger.Field("client_ip", clientIP),
					logger.Field("user_agent", c.Request.UserAgent()),
					logger.Field("path", c.Request.URL.Path),
				)

				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":       "Rate limit exceeded",
					"message":     "Too many requests. Please try again later.",
					"retry_after": rl.getRetryAfter(clientIP),
				})
				c.Abort()
				return
			}

			c.Next()
		*/
	}
}

// allowRequest checks if a request should be allowed
func (rl *RateLimiter) allowRequest(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	tracker, exists := rl.clients[clientIP]

	if !exists {
		tracker = &ClientTracker{
			Requests: make([]time.Time, 0),
			LastSeen: now,
		}
		rl.clients[clientIP] = tracker
	}

	// Check if client is blocked
	if tracker.Blocked {
		// Unblock after 5 minutes
		if time.Since(tracker.BlockedAt) > 5*time.Minute {
			tracker.Blocked = false
			tracker.BlockReason = ""
		} else {
			return false
		}
	}

	// Clean old requests outside the window
	cutoff := now.Add(-rl.config.WindowSize)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range tracker.Requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	tracker.Requests = validRequests

	// Check if adding this request would exceed the limit
	if len(tracker.Requests) >= rl.config.RequestsPerMinute {
		// Check burst allowance
		if len(tracker.Requests) >= rl.config.RequestsPerMinute+rl.config.BurstSize {
			// Block the client
			tracker.Blocked = true
			tracker.BlockedAt = now
			tracker.BlockReason = "Rate limit exceeded"
			return false
		}
	}

	// Add current request
	tracker.Requests = append(tracker.Requests, now)
	tracker.LastSeen = now

	return true
}

// getClientIP extracts the real client IP address
func (rl *RateLimiter) getClientIP(c *gin.Context) string {
	// Check for forwarded headers (common with proxies/load balancers)
	if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		return forwardedFor
	}
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}
	if cfConnectingIP := c.GetHeader("CF-Connecting-IP"); cfConnectingIP != "" {
		return cfConnectingIP
	}

	// Fall back to remote address
	return c.ClientIP()
}

// getRetryAfter returns when the client can retry
func (rl *RateLimiter) getRetryAfter(clientIP string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	tracker, exists := rl.clients[clientIP]
	if !exists || !tracker.Blocked {
		return 0
	}

	// Return seconds until unblock
	unblockTime := tracker.BlockedAt.Add(5 * time.Minute)
	remaining := time.Until(unblockTime)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// cleanupRoutine removes old client trackers to prevent memory leaks
func (rl *RateLimiter) cleanupRoutine() {
	for range rl.cleanupTicker.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)

		for ip, tracker := range rl.clients {
			if tracker.LastSeen.Before(cutoff) {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// GetRateLimitStats returns current rate limiting statistics
func (rl *RateLimiter) GetRateLimitStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	stats := map[string]interface{}{
		"total_clients": len(rl.clients),
		"config":        rl.config,
	}

	blockedClients := 0
	for _, tracker := range rl.clients {
		if tracker.Blocked {
			blockedClients++
		}
	}
	stats["blocked_clients"] = blockedClients

	return stats
}

// Close stops the cleanup routine
func (rl *RateLimiter) Close() {
	if rl.cleanupTicker != nil {
		rl.cleanupTicker.Stop()
	}
}
