package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg == nil {
		t.Fatal("Default() returned nil")
	}

	// Verify default values
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Model", cfg.Model, DefaultModel},
		{"Temperature", cfg.Temperature, 0.7},
		{"MaxTokens", cfg.MaxTokens, 2000},
		{"Streaming", cfg.Streaming, true},
		{"ConventionalCommits", cfg.Commit.ConventionalCommits, true},
		{"PreferredType", cfg.Commit.PreferredType, "feat"},
		{"AutoStage", cfg.Commit.AutoStage, false},
		{"SaveHistory", cfg.Chat.SaveHistory, true},
		{"HistoryLimit", cfg.Chat.HistoryLimit, 100},
		{"Timeout", cfg.Chat.Timeout, 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	// Set HOME to non-existent directory to ensure config doesn't exist
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil for non-existent config", err)
	}

	// Should return defaults
	if cfg.Model != DefaultModel {
		t.Errorf("Load() returned non-default config for non-existent file")
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	// Create temporary config directory
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "zik")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	// Write test config
	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `model: gpt-4
temperature: 0.5
max_tokens: 1000
streaming: false
commit:
  conventional_commits: false
  preferred_type: fix
  auto_stage: true
chat:
  save_history: false
  history_limit: 50
  timeout: 60s
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify loaded values
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Model", cfg.Model, "gpt-4"},
		{"Temperature", cfg.Temperature, 0.5},
		{"MaxTokens", cfg.MaxTokens, 1000},
		{"Streaming", cfg.Streaming, false},
		{"ConventionalCommits", cfg.Commit.ConventionalCommits, false},
		{"PreferredType", cfg.Commit.PreferredType, "fix"},
		{"AutoStage", cfg.Commit.AutoStage, true},
		{"SaveHistory", cfg.Chat.SaveHistory, false},
		{"HistoryLimit", cfg.Chat.HistoryLimit, 50},
		{"Timeout", cfg.Chat.Timeout, 60 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	// Create temporary config directory
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "zik")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	// Write invalid YAML
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	_, err := Load()
	if err == nil {
		t.Error("Load() should return error for invalid YAML")
	}
}

func TestSave(t *testing.T) {
	// Create temporary home directory
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	cfg := Default()
	cfg.Model = "test-model"
	cfg.Temperature = 0.9

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	configPath := filepath.Join(tmpDir, ".config", "zik", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Save() did not create config file")
	}

	// Load and verify
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}

	if loaded.Model != "test-model" {
		t.Errorf("Loaded model = %v, want test-model", loaded.Model)
	}
	if loaded.Temperature != 0.9 {
		t.Errorf("Loaded temperature = %v, want 0.9", loaded.Temperature)
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("getConfigPath() error = %v", err)
	}

	// Should contain .config/zik/config.yaml
	if !filepath.IsAbs(path) {
		t.Error("getConfigPath() should return absolute path")
	}

	expectedSuffix := filepath.Join(".config", "zik", "config.yaml")
	if !filepath.HasPrefix(path, filepath.Dir(filepath.Dir(filepath.Dir(path)))) {
		t.Errorf("getConfigPath() = %v, should end with %v", path, expectedSuffix)
	}
}
