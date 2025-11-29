package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
)

// Health handles health check requests
func Health(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := domain.HealthResponse{
			Status:  "ok",
			Version: cfg.Server.Version,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
