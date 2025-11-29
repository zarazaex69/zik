package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidRequest indicates that the request is malformed or invalid
	ErrInvalidRequest = errors.New("invalid request")

	// ErrUpstreamAPI indicates an error from the upstream Z.AI API
	ErrUpstreamAPI = errors.New("upstream API error")

	// ErrUnauthorized indicates authentication failure
	ErrUnauthorized = errors.New("unauthorized")

	// ErrRateLimited indicates rate limit exceeded
	ErrRateLimited = errors.New("rate limit exceeded")

	// ErrInternalServer indicates an internal server error
	ErrInternalServer = errors.New("internal server error")
)

// APIError represents a structured API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// NewAPIError creates a new APIError
func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Type:    "api_error",
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *APIError {
	return &APIError{
		Code:    400,
		Message: message,
		Type:    "validation_error",
	}
}

// NewUpstreamError creates an upstream API error
func NewUpstreamError(statusCode int, message string) *APIError {
	return &APIError{
		Code:    statusCode,
		Message: fmt.Sprintf("Upstream API error: %s", message),
		Type:    "upstream_error",
	}
}
