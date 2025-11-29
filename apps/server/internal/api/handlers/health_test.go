package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
)

func TestHealth(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Version: "1.0.0",
		},
	}

	handler := Health(cfg)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response domain.HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Health() failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Health() status = %s, want ok", response.Status)
	}

	if response.Version != "1.0.0" {
		t.Errorf("Health() version = %s, want 1.0.0", response.Version)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Health() Content-Type = %s, want application/json", contentType)
	}
}
