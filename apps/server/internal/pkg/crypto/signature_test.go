package crypto

import (
	"os"
	"strings"
	"testing"
)

func TestNewSignatureGenerator(t *testing.T) {
	generator := NewSignatureGenerator()
	if generator == nil {
		t.Fatal("NewSignatureGenerator() returned nil")
	}
}

func TestGenerateSignature_Success(t *testing.T) {
	// Set up environment variable for testing
	originalKey := os.Getenv("ZAI_SECRET_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("ZAI_SECRET_KEY", originalKey)
		} else {
			os.Unsetenv("ZAI_SECRET_KEY")
		}
	}()

	testKey := "test-secret-key-12345"
	os.Setenv("ZAI_SECRET_KEY", testKey)

	generator := NewSignatureGenerator()

	testCases := []struct {
		name            string
		params          map[string]string
		lastUserMessage string
	}{
		{
			name: "basic parameters",
			params: map[string]string{
				"requestId": "req-123",
				"timestamp": "1732912800000",
				"user_id":   "user-456",
			},
			lastUserMessage: "Hello, AI!",
		},
		{
			name: "extra parameters",
			params: map[string]string{
				"requestId": "req-789",
				"timestamp": "1732912800000",
				"user_id":   "user-999",
				"extra":     "value",
			},
			lastUserMessage: "Complex test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := generator.GenerateSignature(tc.params, tc.lastUserMessage)
			if err != nil {
				t.Fatalf("GenerateSignature() error = %v", err)
			}

			if result == nil {
				t.Fatal("GenerateSignature() returned nil result")
			}

			if result.Signature == "" {
				t.Error("Signature should not be empty")
			}

			if result.Timestamp <= 0 {
				t.Error("Timestamp should be positive")
			}

			// Verify signature is deterministic
			result2, err := generator.GenerateSignature(tc.params, tc.lastUserMessage)
			if err != nil {
				t.Fatalf("GenerateSignature() second call error = %v", err)
			}

			// Signatures should be valid hex strings (64 chars for SHA256)
			if len(result.Signature) != 64 {
				t.Errorf("Signature length = %d, want 64", len(result.Signature))
			}

			// Different timestamps but same params should produce same signature
			if result.Signature != result2.Signature {
				t.Error("Same parameters should produce same signature")
			}
		})
	}
}

func TestGenerateSignature_MissingParams(t *testing.T) {
	generator := NewSignatureGenerator()

	testCases := []struct {
		name   string
		params map[string]string
	}{
		{
			name:   "empty params",
			params: map[string]string{},
		},
		{
			name: "missing requestId",
			params: map[string]string{
				"timestamp": "1234567890",
				"user_id":   "user-1",
			},
		},
		{
			name: "missing timestamp",
			params: map[string]string{
				"requestId": "req-1",
				"user_id":   "user-1",
			},
		},
		{
			name: "missing user_id",
			params: map[string]string{
				"requestId": "req-1",
				"timestamp": "1234567890",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := generator.GenerateSignature(tc.params, "test")
			if err == nil {
				t.Error("Expected error when required parameters are missing")
			}
		})
	}
}

func TestGenerateSignature_FallbackKey(t *testing.T) {
	// Ensure secret key is not set
	originalKey := os.Getenv("ZAI_SECRET_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("ZAI_SECRET_KEY", originalKey)
		}
	}()

	os.Unsetenv("ZAI_SECRET_KEY")

	generator := NewSignatureGenerator()
	params := map[string]string{
		"requestId": "req-1",
		"timestamp": "1732912800000",
		"user_id":   "user-1",
	}

	result, err := generator.GenerateSignature(params, "test message")
	if err != nil {
		t.Fatalf("Expected success with fallback key, got error: %v", err)
	}

	if result.Signature == "" {
		t.Error("Signature should not be empty with fallback key")
	}
}

func TestGenerateSignature_DifferentInputs(t *testing.T) {
	generator := NewSignatureGenerator()

	baseParams := map[string]string{
		"requestId": "req-1",
		"timestamp": "1732912800000",
		"user_id":   "user-1",
	}

	result1, _ := generator.GenerateSignature(baseParams, "message1")
	result2, _ := generator.GenerateSignature(baseParams, "message2")

	// Different messages should produce different signatures
	if result1.Signature == result2.Signature {
		t.Error("Different messages should produce different signatures")
	}

	params2 := map[string]string{
		"requestId": "req-2",
		"timestamp": "1732912800000",
		"user_id":   "user-1",
	}
	result3, _ := generator.GenerateSignature(params2, "message1")

	// Different parameters should produce different signatures
	if result1.Signature == result3.Signature {
		t.Error("Different parameters should produce different signatures")
	}
}

func TestGenerateSignature_Deterministic(t *testing.T) {
	generator := NewSignatureGenerator()

	params := map[string]string{
		"timestamp": "1732912800000",
		"requestId": "test-request-id",
		"user_id":   "test-user-id",
	}
	content := "Test content"

	// Generate signature twice with same inputs
	result1, err1 := generator.GenerateSignature(params, content)
	if err1 != nil {
		t.Fatalf("First GenerateSignature() failed: %v", err1)
	}

	result2, err2 := generator.GenerateSignature(params, content)
	if err2 != nil {
		t.Fatalf("Second GenerateSignature() failed: %v", err2)
	}

	// Signatures should be identical for same inputs
	if result1.Signature != result2.Signature {
		t.Errorf("GenerateSignature() not deterministic: got %v and %v", result1.Signature, result2.Signature)
	}
}

func TestGenerateSignature_ValidHexFormat(t *testing.T) {
	generator := NewSignatureGenerator()
	params := map[string]string{
		"requestId": "req-1",
		"timestamp": "1732912800000",
		"user_id":   "user-1",
	}

	result, err := generator.GenerateSignature(params, "message")
	if err != nil {
		t.Fatalf("GenerateSignature() error = %v", err)
	}

	// Verify signature is hex-encoded (only contains 0-9, a-f)
	for _, char := range result.Signature {
		if !strings.ContainsRune("0123456789abcdef", char) {
			t.Errorf("Signature contains invalid hex character: %c", char)
		}
	}

	// SHA256 produces 32 bytes = 64 hex chars
	if len(result.Signature) != 64 {
		t.Errorf("Signature length = %d, want 64", len(result.Signature))
	}
}
