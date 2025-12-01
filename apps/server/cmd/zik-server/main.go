package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/zarazaex69/zik/apps/server/internal/api"
	"github.com/zarazaex69/zik/apps/server/internal/config"
	"github.com/zarazaex69/zik/apps/server/internal/pkg/crypto"
	"github.com/zarazaex69/zik/apps/server/internal/pkg/logger"
	"github.com/zarazaex69/zik/apps/server/internal/pkg/utils"
	"github.com/zarazaex69/zik/apps/server/internal/service/ai"
	"github.com/zarazaex69/zik/apps/server/internal/service/auth"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.Server.Debug)

	logger.Info().Msg("Starting Zik AI Server")
	logger.Info().
		Str("version", cfg.Server.Version).
		Int("port", cfg.Server.Port).
		Bool("debug", cfg.Server.Debug).
		Msg("Configuration loaded")

	// Initialize tokenizer
	tokenizer := utils.NewTokenizer()
	if err := tokenizer.Init(); err != nil {
		logger.Warn().Err(err).Msg("Failed to initialize tokenizer, token counting will be unavailable")
	}

	// Initialize services
	authService := auth.NewService()
	sigGen := crypto.NewSignatureGenerator()
	aiClient := ai.NewClient(cfg, authService, sigGen)

	// Create router
	router := api.NewRouter(cfg, aiClient, tokenizer)

	// Print startup info
	logger.Info().Msg("Available endpoints:")
	logger.Info().Msg("  GET  /health                  - Health check")
	logger.Info().Msg("  GET  /v1/models               - List available models")
	logger.Info().Msg("  POST /v1/chat/completions     - OpenAI-compatible chat completions")

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info().Str("addr", addr).Msg("Server starting")

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
