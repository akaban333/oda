package api

import "time"

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Error represents an API error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta contains metadata for list responses
type Meta struct {
	Total       int       `json:"total"`
	Page        int       `json:"page"`
	PerPage     int       `json:"perPage"`
	TotalPages  int       `json:"totalPages"`
	Timestamp   time.Time `json:"timestamp"`
	ProcessedIn string    `json:"processedIn,omitempty"`
}

// PaginationParams represents pagination parameters for list requests
type PaginationParams struct {
	Page    int `form:"page" binding:"min=1"`
	PerPage int `form:"perPage" binding:"min=1,max=100"`
}

// SortParams represents sorting parameters for list requests
type SortParams struct {
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder" binding:"oneof=asc desc"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Timestamp: time.Now(),
		},
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code int, message string, details string) Response {
	return Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: &Meta{
			Timestamp: time.Now(),
		},
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, total, page, perPage int, processingTime time.Duration) Response {
	totalPages := 0
	if perPage > 0 {
		totalPages = (total + perPage - 1) / perPage
	}

	return Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Total:       total,
			Page:        page,
			PerPage:     perPage,
			TotalPages:  totalPages,
			Timestamp:   time.Now(),
			ProcessedIn: processingTime.String(),
		},
	}
}

// Common error codes
const (
	ErrBadRequest          = 400
	ErrUnauthorized        = 401
	ErrForbidden           = 403
	ErrNotFound            = 404
	ErrMethodNotAllowed    = 405
	ErrConflict            = 409
	ErrInternalServerError = 500
)

// Common error messages
const (
	MsgBadRequest          = "Bad Request"
	MsgUnauthorized        = "Unauthorized"
	MsgForbidden           = "Forbidden"
	MsgNotFound            = "Not Found"
	MsgMethodNotAllowed    = "Method Not Allowed"
	MsgInternalServerError = "Internal Server Error"
	MsgConflict            = "Conflict"
)

// Common validation errors
const (
	MsgInvalidJSON          = "Invalid JSON provided"
	MsgMissingRequiredField = "Missing required field"
	MsgInvalidID            = "Invalid ID format"
	MsgInvalidEmail         = "Invalid email format"
	MsgInvalidPassword      = "Invalid password format"
	MsgInvalidDate          = "Invalid date format"
)
