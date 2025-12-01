package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/zarazaex/zik/apps/server/internal/api/handlers"
	"github.com/zarazaex/zik/apps/server/internal/api/middleware"
	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/pkg/utils"
	"github.com/zarazaex/zik/apps/server/internal/service/ai"
)

// NewRouter creates a new HTTP router with all routes and middleware
func NewRouter(cfg *config.Config, aiClient ai.AIClienter, tokenizer utils.Tokener) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recovery)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimit(30)) // 30 requests per minute per IP

	// Routes
	r.Get("/health", handlers.Health(cfg))
	r.Get("/v1/models", handlers.Models(cfg))
	r.Post("/v1/chat/completions", handlers.ChatCompletions(cfg, aiClient, tokenizer))

	return r
}
