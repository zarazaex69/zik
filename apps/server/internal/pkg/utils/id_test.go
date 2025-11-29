package utils

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	id := GenerateID()

	if id == "" {
		t.Error("GenerateID() returned empty string")
	}

	// UUID v4 format: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
	if len(id) != 36 {
		t.Errorf("GenerateID() returned unexpected length: got %d, want 36", len(id))
	}

	// Check for dashes at correct positions
	if id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		t.Errorf("GenerateID() returned invalid UUID format: %s", id)
	}
}

func TestGenerateID_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		id := GenerateID()
		if ids[id] {
			t.Errorf("GenerateID() generated duplicate ID: %s", id)
		}
		ids[id] = true
	}

	if len(ids) != iterations {
		t.Errorf("GenerateID() uniqueness check failed: got %d unique IDs, want %d", len(ids), iterations)
	}
}

func TestGenerateChatCompletionID(t *testing.T) {
	id := GenerateChatCompletionID()

	if id == "" {
		t.Error("GenerateChatCompletionID() returned empty string")
	}

	// Should have "chatcmpl-" prefix
	prefix := "chatcmpl-"
	if len(id) < len(prefix) {
		t.Errorf("GenerateChatCompletionID() too short: %s", id)
	}

	if id[:len(prefix)] != prefix {
		t.Errorf("GenerateChatCompletionID() missing prefix: got %s, want prefix %s", id, prefix)
	}

	// Rest should be UUID
	uuidPart := id[len(prefix):]
	if len(uuidPart) != 36 {
		t.Errorf("GenerateChatCompletionID() UUID part has wrong length: got %d, want 36", len(uuidPart))
	}
}

func TestGenerateRequestID(t *testing.T) {
	id := GenerateRequestID()

	if id == "" {
		t.Error("GenerateRequestID() returned empty string")
	}

	// Should be a valid UUID
	if len(id) != 36 {
		t.Errorf("GenerateRequestID() returned unexpected length: got %d, want 36", len(id))
	}
}

func TestGenerateRequestID_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		id := GenerateRequestID()
		if ids[id] {
			t.Errorf("GenerateRequestID() generated duplicate ID: %s", id)
		}
		ids[id] = true
	}
}

func TestGenerateChatCompletionID_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		id := GenerateChatCompletionID()
		if ids[id] {
			t.Errorf("GenerateChatCompletionID() generated duplicate ID: %s", id)
		}
		ids[id] = true
	}
}
