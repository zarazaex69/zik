package domain

import (
	"testing"
)

func TestNewAPIError(t *testing.T) {
	code := 400
	message := "Test error"

	err := NewAPIError(code, message)

	if err.Code != code {
		t.Errorf("NewAPIError() code = %d, want %d", err.Code, code)
	}

	if err.Message != message {
		t.Errorf("NewAPIError() message = %s, want %s", err.Message, message)
	}

	if err.Type != "api_error" {
		t.Errorf("NewAPIError() type = %s, want api_error", err.Type)
	}
}

func TestNewValidationError(t *testing.T) {
	message := "Validation failed"

	err := NewValidationError(message)

	if err.Code != 400 {
		t.Errorf("NewValidationError() code = %d, want 400", err.Code)
	}

	if err.Message != message {
		t.Errorf("NewValidationError() message = %s, want %s", err.Message, message)
	}

	if err.Type != "validation_error" {
		t.Errorf("NewValidationError() type = %s, want validation_error", err.Type)
	}
}

func TestNewUpstreamError(t *testing.T) {
	statusCode := 502
	message := "Upstream failed"

	err := NewUpstreamError(statusCode, message)

	if err.Code != statusCode {
		t.Errorf("NewUpstreamError() code = %d, want %d", err.Code, statusCode)
	}

	expectedMsg := "Upstream API error: " + message
	if err.Message != expectedMsg {
		t.Errorf("NewUpstreamError() message = %s, want %s", err.Message, expectedMsg)
	}

	if err.Type != "upstream_error" {
		t.Errorf("NewUpstreamError() type = %s, want upstream_error", err.Type)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		Code:    404,
		Message: "Not found",
		Type:    "api_error",
	}

	expected := "API error 404: Not found"
	if err.Error() != expected {
		t.Errorf("APIError.Error() = %s, want %s", err.Error(), expected)
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that error constants are defined
	errors := []error{
		ErrInvalidRequest,
		ErrUpstreamAPI,
		ErrUnauthorized,
		ErrRateLimited,
		ErrInternalServer,
	}

	for _, err := range errors {
		if err == nil {
			t.Error("Error constant is nil")
		}
		if err.Error() == "" {
			t.Error("Error constant has empty message")
		}
	}
}
