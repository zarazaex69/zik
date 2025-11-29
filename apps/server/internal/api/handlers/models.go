package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
)

// Models handles /v1/models requests
func Models(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := domain.ModelsResponse{
			Object: "list",
			Data: []domain.Model{
				{
					ID:      cfg.Model.Default,
					Object:  "model",
					Created: time.Now().Unix(),
					OwnedBy: "zik",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
