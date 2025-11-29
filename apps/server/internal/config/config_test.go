package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	if cfg.Server.Port != 8080 {
		t.Errorf("defaultConfig() port = %d, want 8080", cfg.Server.Port)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("defaultConfig() host = %s, want 0.0.0.0", cfg.Server.Host)
	}

	if cfg.Server.Version != "0.1.0" {
		t.Errorf("defaultConfig() version = %s, want 0.1.0", cfg.Server.Version)
	}

	if cfg.Upstream.Protocol != "https:" {
		t.Errorf("defaultConfig() protocol = %s, want https:", cfg.Upstream.Protocol)
	}

	if cfg.Upstream.Host != "chat.z.ai" {
		t.Errorf("defaultConfig() upstream host = %s, want chat.z.ai", cfg.Upstream.Host)
	}

	if cfg.Model.Default != "GLM-4-6-API-V1" {
		t.Errorf("defaultConfig() model = %s, want GLM-4-6-API-V1", cfg.Model.Default)
	}

	if cfg.Model.ThinkMode != "reasoning" {
		t.Errorf("defaultConfig() think mode = %s, want reasoning", cfg.Model.ThinkMode)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*Config)
		wantErr bool
	}{
		{
			name:    "valid config",
			modify:  func(c *Config) {},
			wantErr: false,
		},
		{
			name: "invalid port too low",
			modify: func(c *Config) {
				c.Server.Port = 0
			},
			wantErr: true,
		},
		{
			name: "invalid port too high",
			modify: func(c *Config) {
				c.Server.Port = 70000
			},
			wantErr: true,
		},
		{
			name: "invalid think mode",
			modify: func(c *Config) {
				c.Model.ThinkMode = "invalid"
			},
			wantErr: true,
		},
		{
			name: "valid think mode - reasoning",
			modify: func(c *Config) {
				c.Model.ThinkMode = "reasoning"
			},
			wantErr: false,
		},
		{
			name: "valid think mode - think",
			modify: func(c *Config) {
				c.Model.ThinkMode = "think"
			},
			wantErr: false,
		},
		{
			name: "valid think mode - strip",
			modify: func(c *Config) {
				c.Model.ThinkMode = "strip"
			},
			wantErr: false,
		},
		{
			name: "valid think mode - details",
			modify: func(c *Config) {
				c.Model.ThinkMode = "details"
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := defaultConfig()
			tt.modify(cfg)

			err := cfg.validate()

			if tt.wantErr && err == nil {
				t.Error("validate() expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("validate() unexpected error: %v", err)
			}
		})
	}
}

func TestConfig_ApplyEnvOverrides(t *testing.T) {
	// Save original env
	origPort := os.Getenv("PORT")
	origHost := os.Getenv("HOST")
	origDebug := os.Getenv("DEBUG")
	origToken := os.Getenv("ZAI_TOKEN")
	origModel := os.Getenv("MODEL")
	origThink := os.Getenv("THINK_MODE")

	defer func() {
		// Restore original env
		os.Setenv("PORT", origPort)
		os.Setenv("HOST", origHost)
		os.Setenv("DEBUG", origDebug)
		os.Setenv("ZAI_TOKEN", origToken)
		os.Setenv("MODEL", origModel)
		os.Setenv("THINK_MODE", origThink)
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		checkFn  func(*Config) bool
		expected bool
	}{
		{
			name: "override port",
			envVars: map[string]string{
				"PORT": "9000",
			},
			checkFn: func(c *Config) bool {
				return c.Server.Port == 9000
			},
			expected: true,
		},
		{
			name: "override host",
			envVars: map[string]string{
				"HOST": "127.0.0.1",
			},
			checkFn: func(c *Config) bool {
				return c.Server.Host == "127.0.0.1"
			},
			expected: true,
		},
		{
			name: "override debug",
			envVars: map[string]string{
				"DEBUG": "true",
			},
			checkFn: func(c *Config) bool {
				return c.Server.Debug == true
			},
			expected: true,
		},
		{
			name: "override token",
			envVars: map[string]string{
				"ZAI_TOKEN": "test-token",
			},
			checkFn: func(c *Config) bool {
				return c.Upstream.Token == "test-token" && !c.Upstream.Anonymous
			},
			expected: true,
		},
		{
			name: "override model",
			envVars: map[string]string{
				"MODEL": "custom-model",
			},
			checkFn: func(c *Config) bool {
				return c.Model.Default == "custom-model"
			},
			expected: true,
		},
		{
			name: "override think mode",
			envVars: map[string]string{
				"THINK_MODE": "strip",
			},
			checkFn: func(c *Config) bool {
				return c.Model.ThinkMode == "strip"
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars
			os.Unsetenv("PORT")
			os.Unsetenv("HOST")
			os.Unsetenv("DEBUG")
			os.Unsetenv("ZAI_TOKEN")
			os.Unsetenv("MODEL")
			os.Unsetenv("THINK_MODE")

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg := defaultConfig()
			cfg.applyEnvOverrides()

			result := tt.checkFn(cfg)
			if result != tt.expected {
				t.Errorf("applyEnvOverrides() check failed: got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConfig_GetUpstreamHeaders(t *testing.T) {
	cfg := defaultConfig()
	headers := cfg.GetUpstreamHeaders()

	requiredHeaders := []string{
		"Accept",
		"Accept-Language",
		"User-Agent",
		"Origin",
		"Referer",
	}

	for _, header := range requiredHeaders {
		if _, ok := headers[header]; !ok {
			t.Errorf("GetUpstreamHeaders() missing required header: %s", header)
		}
	}

	// Check Origin format
	expectedOrigin := cfg.Upstream.Protocol + "//" + cfg.Upstream.Host
	if headers["Origin"] != expectedOrigin {
		t.Errorf("GetUpstreamHeaders() Origin = %s, want %s", headers["Origin"], expectedOrigin)
	}
}
