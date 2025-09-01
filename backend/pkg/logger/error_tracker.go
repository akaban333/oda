package logger

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity string

const (
	// SeverityLow represents a low severity error
	SeverityLow ErrorSeverity = "low"
	// SeverityMedium represents a medium severity error
	SeverityMedium ErrorSeverity = "medium"
	// SeverityHigh represents a high severity error
	SeverityHigh ErrorSeverity = "high"
	// SeverityCritical represents a critical severity error
	SeverityCritical ErrorSeverity = "critical"
)

// ErrorCategory represents the category of an error
type ErrorCategory string

const (
	// CategoryDatabase represents database-related errors
	CategoryDatabase ErrorCategory = "database"
	// CategoryAuthentication represents authentication-related errors
	CategoryAuthentication ErrorCategory = "authentication"
	// CategoryValidation represents validation-related errors
	CategoryValidation ErrorCategory = "validation"
	// CategoryNetwork represents network-related errors
	CategoryNetwork ErrorCategory = "network"
	// CategorySystem represents system-related errors
	CategorySystem ErrorCategory = "system"
	// CategoryBusiness represents business logic errors
	CategoryBusiness ErrorCategory = "business"
	// CategoryUnknown represents unknown error categories
	CategoryUnknown ErrorCategory = "unknown"
)

// TrackedError represents a tracked error with metadata
type TrackedError struct {
	ID          string                 `json:"id"`
	Message     string                 `json:"message"`
	Error       string                 `json:"error"`
	Severity    ErrorSeverity          `json:"severity"`
	Category    ErrorCategory          `json:"category"`
	Stack       string                 `json:"stack,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	Path        string                 `json:"path,omitempty"`
	Method      string                 `json:"method,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	ClientIP    string                 `json:"client_ip,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Occurrences int                    `json:"occurrences"`
	FirstSeen   time.Time              `json:"first_seen"`
	LastSeen    time.Time              `json:"last_seen"`
	Resolved    bool                   `json:"resolved"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
}

// ErrorTracker manages error tracking and aggregation
type ErrorTracker struct {
	errors    map[string]*TrackedError
	mu        sync.RWMutex
	stats     ErrorStats
	callbacks []ErrorCallback
}

// ErrorStats represents error tracking statistics
type ErrorStats struct {
	TotalErrors      int                   `json:"total_errors"`
	ErrorsBySeverity map[ErrorSeverity]int `json:"errors_by_severity"`
	ErrorsByCategory map[ErrorCategory]int `json:"errors_by_category"`
	ErrorsByHour     map[int]int           `json:"errors_by_hour"`
	RecentErrors     []*TrackedError       `json:"recent_errors"`
}

// ErrorCallback is a function called when errors occur
type ErrorCallback func(*TrackedError)

// NewErrorTracker creates a new error tracker
func NewErrorTracker() *ErrorTracker {
	return &ErrorTracker{
		errors:    make(map[string]*TrackedError),
		callbacks: make([]ErrorCallback, 0),
		stats: ErrorStats{
			ErrorsBySeverity: make(map[ErrorSeverity]int),
			ErrorsByCategory: make(map[ErrorCategory]int),
			ErrorsByHour:     make(map[int]int),
			RecentErrors:     make([]*TrackedError, 0),
		},
	}
}

// TrackError tracks a new error occurrence
func (et *ErrorTracker) TrackError(err error, severity ErrorSeverity, category ErrorCategory, fields ...zapcore.Field) *TrackedError {
	et.mu.Lock()
	defer et.mu.Unlock()

	// Create error fingerprint
	fingerprint := et.createFingerprint(err, severity, category, fields...)

	// Check if we've seen this error before
	trackedError, exists := et.errors[fingerprint]
	if exists {
		// Update existing error
		trackedError.Occurrences++
		trackedError.LastSeen = time.Now()

		// Update metadata if new fields provided
		if len(fields) > 0 {
			for _, field := range fields {
				if field.Key != "" {
					// Extract value from zapcore.Field based on its type
					var value interface{}
					switch field.Type {
					case zapcore.StringType:
						value = field.String
					case zapcore.Int64Type:
						value = field.Integer
					case zapcore.Float64Type:
						value = field.Interface
					case zapcore.BoolType:
						value = field.Integer == 1
					default:
						value = field.String // fallback to string representation
					}
					trackedError.Metadata[field.Key] = value
				}
			}
		}
	} else {
		// Create new tracked error
		now := time.Now()
		trackedError = &TrackedError{
			ID:          fingerprint,
			Message:     err.Error(),
			Error:       err.Error(),
			Severity:    severity,
			Category:    category,
			Occurrences: 1,
			FirstSeen:   now,
			LastSeen:    now,
			Metadata:    make(map[string]interface{}),
			Resolved:    false,
		}

		// Extract metadata from fields
		for _, field := range fields {
			if field.Key != "" {
				// Extract value from zapcore.Field based on its type
				var value interface{}
				switch field.Type {
				case zapcore.StringType:
					value = field.String
				case zapcore.Int64Type:
					value = field.Integer
				case zapcore.Float64Type:
					value = field.Interface
				case zapcore.BoolType:
					value = field.Integer == 1
				default:
					value = field.Interface // fallback to interface
				}
				trackedError.Metadata[field.Key] = value
			}
		}

		et.errors[fingerprint] = trackedError
	}

	// Update statistics
	et.updateStats(trackedError)

	// Call callbacks
	et.notifyCallbacks(trackedError)

	return trackedError
}

// createFingerprint creates a unique fingerprint for an error
func (et *ErrorTracker) createFingerprint(err error, severity ErrorSeverity, category ErrorCategory, fields ...zapcore.Field) string {
	// Create a simple fingerprint based on error message, severity, and category
	// In production, you might want to use a more sophisticated algorithm
	return fmt.Sprintf("%s:%s:%s", severity, category, err.Error())
}

// updateStats updates error tracking statistics
func (et *ErrorTracker) updateStats(error *TrackedError) {
	et.stats.TotalErrors++
	et.stats.ErrorsBySeverity[error.Severity]++
	et.stats.ErrorsByCategory[error.Category]++

	// Track errors by hour
	hour := error.LastSeen.Hour()
	et.stats.ErrorsByHour[hour]++

	// Update recent errors (keep last 100)
	et.stats.RecentErrors = append(et.stats.RecentErrors, error)
	if len(et.stats.RecentErrors) > 100 {
		et.stats.RecentErrors = et.stats.RecentErrors[1:]
	}
}

// notifyCallbacks calls all registered error callbacks
func (et *ErrorTracker) notifyCallbacks(error *TrackedError) {
	for _, callback := range et.callbacks {
		go func(cb ErrorCallback, err *TrackedError) {
			defer func() {
				if r := recover(); r != nil {
					// Prevent callback panics from crashing the tracker
					Error("Error callback panicked", zapcore.Field{Key: "panic", Type: zapcore.StringType, String: fmt.Sprintf("%v", r)})
				}
			}()
			cb(err)
		}(callback, error)
	}
}

// AddCallback registers a new error callback
func (et *ErrorTracker) AddCallback(callback ErrorCallback) {
	et.mu.Lock()
	defer et.mu.Unlock()
	et.callbacks = append(et.callbacks, callback)
}

// GetErrorStats returns current error statistics
func (et *ErrorTracker) GetErrorStats() ErrorStats {
	et.mu.RLock()
	defer et.mu.RUnlock()
	return et.stats
}

// GetErrorsBySeverity returns errors filtered by severity
func (et *ErrorTracker) GetErrorsBySeverity(severity ErrorSeverity) []*TrackedError {
	et.mu.RLock()
	defer et.mu.RUnlock()

	var errors []*TrackedError
	for _, err := range et.errors {
		if err.Severity == severity {
			errors = append(errors, err)
		}
	}
	return errors
}

// GetErrorsByCategory returns errors filtered by category
func (et *ErrorTracker) GetErrorsByCategory(category ErrorCategory) []*TrackedError {
	et.mu.RLock()
	defer et.mu.RUnlock()

	var errors []*TrackedError
	for _, err := range et.errors {
		if err.Category == category {
			errors = append(errors, err)
		}
	}
	return errors
}

// GetRecentErrors returns the most recent errors
func (et *ErrorTracker) GetRecentErrors(limit int) []*TrackedError {
	et.mu.RLock()
	defer et.mu.RUnlock()

	if limit <= 0 || limit > len(et.stats.RecentErrors) {
		limit = len(et.stats.RecentErrors)
	}

	recent := make([]*TrackedError, limit)
	copy(recent, et.stats.RecentErrors[len(et.stats.RecentErrors)-limit:])
	return recent
}

// MarkErrorResolved marks an error as resolved
func (et *ErrorTracker) MarkErrorResolved(errorID string) bool {
	et.mu.Lock()
	defer et.mu.Unlock()

	if error, exists := et.errors[errorID]; exists {
		now := time.Now()
		error.Resolved = true
		error.ResolvedAt = &now
		return true
	}
	return false
}

// GetUnresolvedErrors returns all unresolved errors
func (et *ErrorTracker) GetUnresolvedErrors() []*TrackedError {
	et.mu.RLock()
	defer et.mu.RUnlock()

	var unresolved []*TrackedError
	for _, err := range et.errors {
		if !err.Resolved {
			unresolved = append(unresolved, err)
		}
	}
	return unresolved
}

// ExportErrors exports errors to JSON format
func (et *ErrorTracker) ExportErrors() ([]byte, error) {
	et.mu.RLock()
	defer et.mu.RUnlock()

	export := struct {
		Timestamp time.Time                `json:"timestamp"`
		Stats     ErrorStats               `json:"stats"`
		Errors    map[string]*TrackedError `json:"errors"`
	}{
		Timestamp: time.Now(),
		Stats:     et.stats,
		Errors:    et.errors,
	}

	return json.MarshalIndent(export, "", "  ")
}

// ClearOldErrors removes errors older than the specified duration
func (et *ErrorTracker) ClearOldErrors(olderThan time.Duration) int {
	et.mu.Lock()
	defer et.mu.Unlock()

	cutoff := time.Now().Add(-olderThan)
	removed := 0

	for id, err := range et.errors {
		if err.LastSeen.Before(cutoff) {
			delete(et.errors, id)
			removed++
		}
	}

	return removed
}
