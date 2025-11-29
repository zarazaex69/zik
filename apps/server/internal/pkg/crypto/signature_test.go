package crypto

import (
	"testing"
)

func TestGenerateSignature(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid signature generation",
			params: map[string]string{
				"timestamp": "1732912800000",
				"requestId": "test-request-id",
				"user_id":   "test-user-id",
			},
			content: "Hello, world!",
			wantErr: false,
		},
		{
			name: "missing timestamp",
			params: map[string]string{
				"requestId": "test-request-id",
				"user_id":   "test-user-id",
			},
			content:     "Hello, world!",
			wantErr:     true,
			errContains: "missing required parameter: timestamp",
		},
		{
			name: "missing requestId",
			params: map[string]string{
				"timestamp": "1732912800000",
				"user_id":   "test-user-id",
			},
			content:     "Hello, world!",
			wantErr:     true,
			errContains: "missing required parameter: requestId",
		},
		{
			name: "missing user_id",
			params: map[string]string{
				"timestamp": "1732912800000",
				"requestId": "test-request-id",
			},
			content:     "Hello, world!",
			wantErr:     true,
			errContains: "missing required parameter: user_id",
		},
		{
			name: "invalid timestamp format",
			params: map[string]string{
				"timestamp": "invalid",
				"requestId": "test-request-id",
				"user_id":   "test-user-id",
			},
			content:     "Hello, world!",
			wantErr:     true,
			errContains: "invalid timestamp",
		},
		{
			name: "empty content",
			params: map[string]string{
				"timestamp": "1732912800000",
				"requestId": "test-request-id",
				"user_id":   "test-user-id",
			},
			content: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateSignature(tt.params, tt.content)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateSignature() expected error but got nil")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GenerateSignature() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateSignature() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Error("GenerateSignature() returned nil result")
				return
			}

			if result.Signature == "" {
				t.Error("GenerateSignature() returned empty signature")
			}

			if result.Timestamp == 0 {
				t.Error("GenerateSignature() returned zero timestamp")
			}
		})
	}
}

func TestGenerateSignature_Deterministic(t *testing.T) {
	params := map[string]string{
		"timestamp": "1732912800000",
		"requestId": "test-request-id",
		"user_id":   "test-user-id",
	}
	content := "Test content"

	// Generate signature twice with same inputs
	result1, err1 := GenerateSignature(params, content)
	if err1 != nil {
		t.Fatalf("First GenerateSignature() failed: %v", err1)
	}

	result2, err2 := GenerateSignature(params, content)
	if err2 != nil {
		t.Fatalf("Second GenerateSignature() failed: %v", err2)
	}

	// Signatures should be identical for same inputs
	if result1.Signature != result2.Signature {
		t.Errorf("GenerateSignature() not deterministic: got %v and %v", result1.Signature, result2.Signature)
	}
}

func TestGenerateSignature_DifferentContent(t *testing.T) {
	params := map[string]string{
		"timestamp": "1732912800000",
		"requestId": "test-request-id",
		"user_id":   "test-user-id",
	}

	result1, _ := GenerateSignature(params, "content1")
	result2, _ := GenerateSignature(params, "content2")

	// Different content should produce different signatures
	if result1.Signature == result2.Signature {
		t.Error("GenerateSignature() produced same signature for different content")
	}
}

func TestHmacSHA256(t *testing.T) {
	key := []byte("test-key")
	message := []byte("test-message")

	result := hmacSHA256(key, message)

	if result == "" {
		t.Error("hmacSHA256() returned empty string")
	}

	// Should return hex-encoded string
	if len(result) != 64 { // SHA256 produces 32 bytes = 64 hex chars
		t.Errorf("hmacSHA256() returned unexpected length: got %d, want 64", len(result))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
