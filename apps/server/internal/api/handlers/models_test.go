package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zarazaex69/zik/apps/server/internal/config"
	"github.com/zarazaex69/zik/apps/server/internal/domain"
)

func TestModels(t *testing.T) {
	cfg := &config.Config{
		Model: config.ModelConfig{
			Default: "test-model-v1",
		},
	}

	handler := Models(cfg)

	req := httptest.NewRequest("GET", "/v1/models", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Models() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response domain.ModelsResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Models() failed to decode response: %v", err)
	}

	if response.Object != "list" {
		t.Errorf("Models() object = %s, want list", response.Object)
	}

	if len(response.Data) == 0 {
		t.Error("Models() returned empty data array")
	}

	model := response.Data[0]
	if model.ID != "test-model-v1" {
		t.Errorf("Models() model ID = %s, want test-model-v1", model.ID)
	}

	if model.Object != "model" {
		t.Errorf("Models() model object = %s, want model", model.Object)
	}

	if model.OwnedBy != "zik" {
		t.Errorf("Models() model owned_by = %s, want zik", model.OwnedBy)
	}

	if model.Created == 0 {
		t.Error("Models() model created timestamp is zero")
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Models() Content-Type = %s, want application/json", contentType)
	}
}
