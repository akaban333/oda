package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateToken generates a random token string
func GenerateToken(length int) (string, error) {
	if length < 1 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	// Calculate the number of bytes needed for the desired token length
	// when encoded as base64
	bytesNeeded := (length*3)/4 + 1

	tokenBytes := make([]byte, bytesNeeded)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token[:length], nil
}

// StructToMap converts a struct to a map[string]interface{}
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// IsValidEmail validates an email address format
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(pattern).MatchString(email)
}

// FormatDatetime formats a time.Time to a string
func FormatDatetime(t time.Time, format string) string {
	if format == "" {
		format = time.RFC3339
	}
	return t.Format(format)
}

// GetCurrentTimestamp returns the current Unix timestamp
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// SanitizeString removes special characters from a string
func SanitizeString(input string) string {
	// Remove any character that's not alphanumeric, space, or basic punctuation
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s.,!?-]`)
	return reg.ReplaceAllString(input, "")
}

// TruncateString truncates a string to a specified length
func TruncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	return input[:maxLength-3] + "..."
}

// RespondWithJSON sends a JSON response with the given status code and payload
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error marshalling JSON"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// RespondWithError sends an error response with the given status code and error message
func RespondWithError(w http.ResponseWriter, statusCode int, error string) {
	RespondWithJSON(w, statusCode, map[string]string{"error": error})
}

// SliceContains checks if a slice contains a specified element
func SliceContains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// GetFileExtension extracts the file extension from a filename
func GetFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

// GetContentType detects content type from filename extension
func GetContentType(filename string) string {
	ext := strings.ToLower(GetFileExtension(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".doc", ".docx":
		return "application/msword"
	case ".xls", ".xlsx":
		return "application/vnd.ms-excel"
	case ".ppt", ".pptx":
		return "application/vnd.ms-powerpoint"
	case ".mp3":
		return "audio/mpeg"
	case ".mp4":
		return "video/mp4"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}
