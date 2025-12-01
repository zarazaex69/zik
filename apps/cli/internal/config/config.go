package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents user-configurable settings (non-security critical)
type Config struct {
	// Model settings
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
	Streaming   bool    `yaml:"streaming"`

	// Commit settings
	Commit CommitConfig `yaml:"commit"`

	// Chat settings
	Chat ChatConfig `yaml:"chat"`
}

// CommitConfig holds commit message generation settings
type CommitConfig struct {
	ConventionalCommits bool   `yaml:"conventional_commits"`
	PreferredType       string `yaml:"preferred_type"`
	AutoStage           bool   `yaml:"auto_stage"`
}

// ChatConfig holds interactive chat settings
type ChatConfig struct {
	SaveHistory  bool          `yaml:"save_history"`
	HistoryLimit int           `yaml:"history_limit"`
	Timeout      time.Duration `yaml:"timeout"`
}

// Default returns a config with sensible defaults
func Default() *Config {
	return &Config{
		Model:       DefaultModel,
		Temperature: 0.7,
		MaxTokens:   2000,
		Streaming:   true,
		Commit: CommitConfig{
			ConventionalCommits: true,
			PreferredType:       "feat",
			AutoStage:           false,
		},
		Chat: ChatConfig{
			SaveHistory:  true,
			HistoryLimit: 100,
			Timeout:      30 * time.Second,
		},
	}
}

// Load reads config from ~/.config/zik/config.yaml
// Falls back to defaults if file doesn't exist
func Load() (*Config, error) {
	cfg := Default()

	configPath, err := getConfigPath()
	if err != nil {
		return cfg, nil // Return defaults if config path cannot be determined
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil // Return defaults if config doesn't exist
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save writes config to ~/.config/zik/config.yaml
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configPath, data, 0644)
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "zik", "config.yaml"), nil
}
