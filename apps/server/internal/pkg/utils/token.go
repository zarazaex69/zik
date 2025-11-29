package utils

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/pkoukk/tiktoken-go"
	"github.com/zarazaex/zik/apps/server/internal/pkg/logger"
)

var (
	encoder *tiktoken.Tiktoken
	once    sync.Once
	initErr error
)

// InitTokenizer initializes the tiktoken encoder for token counting
func InitTokenizer() error {
	once.Do(func() {
		// Set cache directory to avoid downloading on every run
		cacheDir := filepath.Join(".", "tiktoken")
		os.Setenv("TIKTOKEN_CACHE_DIR", cacheDir)

		// Ensure cache directory exists
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			initErr = err
			logger.Warn().Err(err).Msg("Failed to create tiktoken cache directory")
			return
		}

		// Initialize encoder with cl100k_base (used by GPT-4 and compatible models)
		var err error
		encoder, err = tiktoken.GetEncoding("cl100k_base")
		if err != nil {
			initErr = err
			logger.Warn().Err(err).Msg("Failed to initialize tiktoken encoder")
			return
		}

		logger.Info().Msg("Tokenizer initialized successfully")
	})

	return initErr
}

// CountTokens counts the number of tokens in the given text
// Returns 0 if tokenizer is not initialized
func CountTokens(text string) int {
	if encoder == nil {
		if err := InitTokenizer(); err != nil {
			logger.Debug().Err(err).Msg("Tokenizer not available, returning 0 tokens")
			return 0
		}
	}

	if encoder == nil {
		return 0
	}

	tokens := encoder.Encode(text, nil, nil)
	return len(tokens)
}
