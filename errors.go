package ujeebu

import (
	"fmt"
	"net/http"
)

// APIError represents an error response from the Ujeebu API
type APIError struct {
	// URL is the URL that was requested
	URL string `json:"url,omitempty"`
	// Message is the human-readable error message
	Message string `json:"message"`
	// ErrorCode is the API-specific error code
	ErrorCode string `json:"error_code,omitempty"`
	// Errors is a list of detailed error messages
	Errors []string `json:"errors,omitempty"`
	// StatusCode is the HTTP status code
	StatusCode int `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("ujeebu API error (code: %s, status: %d): %s", e.ErrorCode, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("ujeebu API error (status: %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found error
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsTimeout returns true if the error is a 408 Request Timeout error
func (e *APIError) IsTimeout() bool {
	return e.StatusCode == http.StatusRequestTimeout
}

// IsRateLimited returns true if the error is a 429 Too Many Requests error
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// ValidationError represents a client-side validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NetworkError represents a network-related error
type NetworkError struct {
	Err error
}

// Error implements the error interface
func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

// Unwrap returns the underlying error
func (e *NetworkError) Unwrap() error {
	return e.Err
}
