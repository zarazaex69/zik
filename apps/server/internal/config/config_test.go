package config

import (
	"os"
	"strings"
	"sync"
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

func Test_getEnvInt(t *testing.T) {
	// Save original env
	origEnv := os.Getenv("TEST_INT_VAR")
	t.Cleanup(func() {
		os.Setenv("TEST_INT_VAR", origEnv)
	})

	tests := []struct {
		name         string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "env var set to valid int",
			envValue:     "123",
			defaultValue: 10,
			expected:     123,
		},
		{
			name:         "env var not set",
			envValue:     "", // Simulates unset env var
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "env var set to invalid int",
			envValue:     "abc",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "env var set to empty string",
			envValue:     " ", // space treated as invalid int
			defaultValue: 10,
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue == "" {
				os.Unsetenv("TEST_INT_VAR")
			} else {
				os.Setenv("TEST_INT_VAR", tt.envValue)
			}

			result := getEnvInt("TEST_INT_VAR", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnvInt() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func Test_getEnvBool(t *testing.T) {
	// Save original env
	origEnv := os.Getenv("TEST_BOOL_VAR")
	t.Cleanup(func() {
		os.Setenv("TEST_BOOL_VAR", origEnv)
	})

	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "env var set to true",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "env var set to 1",
			envValue:     "1",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "env var set to TRUE (case-insensitive)",
			envValue:     "TRUE",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "env var set to false",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "env var set to 0",
			envValue:     "0",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "env var not set",
			envValue:     "", // Simulates unset env var
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "env var set to invalid string",
			envValue:     "something",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "env var set to empty string", // space treated as invalid bool
			envValue:     " ",
			defaultValue: true,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue == "" {
				os.Unsetenv("TEST_BOOL_VAR")
			} else {
				os.Setenv("TEST_BOOL_VAR", tt.envValue)
			}

			result := getEnvBool("TEST_BOOL_VAR", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnvBool() = %t, want %t", result, tt.expected)
			}
		})
	}
}

func TestLoadConfig_Default(t *testing.T) {
	// Clear all relevant environment variables to ensure defaults are used initially
	os.Unsetenv("PORT")
	os.Unsetenv("HOST")
	os.Unsetenv("DEBUG")
	os.Unsetenv("ZAI_TOKEN")
	os.Unsetenv("MODEL")
	os.Unsetenv("THINK_MODE")

	cfg, err := loadConfig("")
	if err != nil {
		t.Fatalf("loadConfig() with empty path returned an error: %v", err)
	}

	// Verify default values
	expectedDefault := defaultConfig()
	if cfg.Server.Port != expectedDefault.Server.Port {
		t.Errorf("loadConfig() default port = %d, want %d", cfg.Server.Port, expectedDefault.Server.Port)
	}
	if cfg.Model.Default != expectedDefault.Model.Default {
		t.Errorf("loadConfig() default model = %s, want %s", cfg.Model.Default, expectedDefault.Model.Default)
	}

	// Test with environment override
	os.Setenv("PORT", "9000")
	t.Cleanup(func() { os.Unsetenv("PORT") })

	cfgWithEnv, err := loadConfig("") // Reload config to pick up env var
	if err != nil {
		t.Fatalf("loadConfig() with empty path and env var returned an error: %v", err)
	}
	if cfgWithEnv.Server.Port != 9000 {
		t.Errorf("loadConfig() with env override port = %d, want %d", cfgWithEnv.Server.Port, 9000)
	}
}

// createTempConfigFile creates a temporary file with the given content and returns its path.
// It also ensures the file is cleaned up after the test.
func createTempConfigFile(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tmpfile.Close() // Close immediately to allow other processes to open it

	_, err = tmpfile.WriteString(content)
	if err != nil {
		os.Remove(tmpfile.Name()) // Clean up if write fails
		t.Fatalf("failed to write to temp file: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(tmpfile.Name())
	})

	return tmpfile.Name()
}

func TestLoadConfig_FromFile(t *testing.T) {
	// Sample valid config content
	validConfigContent := `
server:
  port: 8888
  host: "127.0.0.1"
  debug: true
upstream:
  host: "test.z.ai"
  token: "test-token-from-file"
model:
  default: "test-model"
  think_mode: "strip"
`
	configPath := createTempConfigFile(t, validConfigContent)

	// Clear relevant env vars to ensure file config is used
	os.Unsetenv("PORT")
	os.Unsetenv("HOST")
	os.Unsetenv("DEBUG")
	os.Unsetenv("ZAI_TOKEN")
	os.Unsetenv("MODEL")
	os.Unsetenv("THINK_MODE")
	
	// Ensure that global config is reset for this test to avoid sync.Once issues
	// This is a hack for testing singletons and should ideally be avoided by refactoring
	// for testability, but given the current structure, it's necessary for isolation.
	cfg = nil // reset global config
	once = sync.Once{} // reset sync.Once

	loadedConfig, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("loadConfig() from valid file returned an unexpected error: %v", err)
	}

	// Assert values loaded from file
	if loadedConfig.Server.Port != 8888 {
		t.Errorf("loadedConfig.Server.Port = %d, want %d", loadedConfig.Server.Port, 8888)
	}
	if loadedConfig.Upstream.Host != "test.z.ai" {
		t.Errorf("loadedConfig.Upstream.Host = %s, want %s", loadedConfig.Upstream.Host, "test.z.ai")
	}
	if loadedConfig.Upstream.Token != "test-token-from-file" {
		t.Errorf("loadedConfig.Upstream.Token = %s, want %s", loadedConfig.Upstream.Token, "test-token-from-file")
	}
	if loadedConfig.Model.ThinkMode != "strip" {
		t.Errorf("loadedConfig.Model.ThinkMode = %s, want %s", loadedConfig.Model.ThinkMode, "strip")
	}
	// Check that anonymous is correctly set to false because token is provided
	if loadedConfig.Upstream.Anonymous != false {
		t.Errorf("loadedConfig.Upstream.Anonymous = %t, want %t", loadedConfig.Upstream.Anonymous, false)
	}

	// Test with environment override on top of file config
	os.Setenv("PORT", "9999")
	os.Setenv("ZAI_TOKEN", "env-token")
	t.Cleanup(func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ZAI_TOKEN")
	})

	// Ensure that global config is reset for this sub-test
	cfg = nil // reset global config
	once = sync.Once{} // reset sync.Once
	
	loadedConfigWithEnv, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("loadConfig() from file with env override returned an error: %v", err)
	}

	if loadedConfigWithEnv.Server.Port != 9999 {
		t.Errorf("loadedConfigWithEnv.Server.Port = %d, want %d", loadedConfigWithEnv.Server.Port, 9999)
	}
	if loadedConfigWithEnv.Upstream.Token != "env-token" {
		t.Errorf("loadedConfigWithEnv.Upstream.Token = %s, want %s", loadedConfigWithEnv.Upstream.Token, "env-token")
	}
	if loadedConfigWithEnv.Upstream.Anonymous != false {
		t.Errorf("loadedConfigWithEnv.Upstream.Anonymous = %t, want %t", loadedConfigWithEnv.Upstream.Anonymous, false)
	}
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	// Ensure that global config is reset for this test
	cfg = nil
	once = sync.Once{}

	_, err := loadConfig("/path/to/non/existent/config.yaml")
	if err == nil {
		t.Error("loadConfig() from non-existent file expected an error, but got none")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("loadConfig() from non-existent file got unexpected error: %v", err)
	}
}

func TestLoadConfig_InvalidYaml(t *testing.T) {
	invalidConfigContent := `
server:
  port: 8080
  host: "0.0.0.0"
  debug: true
  invalid_field: [
    this is not valid yaml
`
	configPath := createTempConfigFile(t, invalidConfigContent)

	// Ensure that global config is reset for this test
	cfg = nil
	once = sync.Once{}

	_, err := loadConfig(configPath)
	if err == nil {
		t.Error("loadConfig() from invalid YAML expected an error, but got none")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to parse config file") {
		t.Errorf("loadConfig() from invalid YAML got unexpected error: %v", err)
	}
}

// TestLoad_Idempotency ensures that Load always returns the same config instance
// after the first successful call, due to sync.Once.
func TestLoad_Idempotency(t *testing.T) {
	// Use a temporary config file for the initial load
	configContent := `server: {port: 1111}`
	configPath := createTempConfigFile(t, configContent)

	// Clear global state to ensure a fresh load
	cfg = nil
	once = sync.Once{}

	// First load
	firstConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("first Load() failed: %v", err)
	}

	// Attempt a second load with different parameters - it should return the same config
	// and ignore the new path due to sync.Once
	secondConfig, err := Load("/path/to/another/config.yaml")
	if err != nil {
		t.Fatalf("second Load() failed: %v", err)
	}

	if firstConfig != secondConfig {
		t.Error("Load() is not idempotent: returned different config instances")
	}

	if firstConfig.Server.Port != 1111 {
		t.Errorf("Loaded config port was %d, expected 1111", firstConfig.Server.Port)
	}

	// Verify that the second load did not change the config
	if secondConfig.Server.Port != 1111 {
		t.Errorf("Second load changed config port to %d, expected 1111", secondConfig.Server.Port)
	}
}

// TestGet_Fallback ensures Get returns a default config if Load hasn't been called.
func TestGet_Fallback(t *testing.T) {
	// Clear global state to ensure Get falls back
	cfg = nil
	once = sync.Once{}

	defaultCfg := Get()

	if defaultCfg == nil {
		t.Fatal("Get() returned nil config during fallback")
	}
	// Check if it's indeed the default config (e.g., check a default port)
	if defaultCfg.Server.Port != defaultConfig().Server.Port {
		t.Errorf("Get() fallback config port was %d, expected default %d", defaultCfg.Server.Port, defaultConfig().Server.Port)
	}

	// Ensure that subsequent Get calls return the same instance
	anotherCfg := Get()
	if defaultCfg != anotherCfg {
		t.Error("Get() subsequent calls returned different instances during fallback")
	}
}

// TestLoad_Error ensures that Load correctly propagates errors from loadConfig.
func TestLoad_Error(t *testing.T) {
	// Clear global state to ensure a fresh load attempt
	cfg = nil
	once = sync.Once{}

	// Attempt to load from a non-existent path
	_, err := Load("/path/to/definitely/non/existent/file.yaml")
	if err == nil {
		t.Error("Load() from non-existent file expected an error, but got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("Load() from non-existent file got unexpected error: %v", err)
	}
}



