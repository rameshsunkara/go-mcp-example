package config_test

import (
	"strings"
	"testing"

	"github.com/rameshsunkara/go-mcp-example/config"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns environment variable when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "from_env",
			expected:     "from_env",
		},
		{
			name:         "returns default when environment variable not set",
			key:          "NONEXISTENT_KEY",
			defaultValue: "default_value",
			envValue:     "",
			expected:     "default_value",
		},
		{
			name:         "returns environment variable when empty string set",
			key:          "EMPTY_TEST_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if needed
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}

			result := config.GetEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  config.Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with all fields",
			config: config.Config{
				HTTPAddr:   "localhost:8080",
				LogLevel:   "info",
				LogFormat:  "json",
				APIKey:     "test-key",
				APIBaseURL: "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "valid config with minimal fields",
			config: config.Config{
				LogLevel:   "debug",
				LogFormat:  "text",
				APIBaseURL: "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "valid config with port only HTTP address",
			config: config.Config{
				HTTPAddr:   ":8080",
				LogLevel:   "warn",
				LogFormat:  "json",
				APIBaseURL: "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: config.Config{
				LogLevel:   "invalid",
				LogFormat:  "json",
				APIBaseURL: "https://api.example.com",
			},
			wantErr: true,
			errMsg:  "invalid log level 'invalid'",
		},
		{
			name: "invalid log format",
			config: config.Config{
				LogLevel:   "info",
				LogFormat:  "invalid",
				APIBaseURL: "https://api.example.com",
			},
			wantErr: true,
			errMsg:  "invalid log format 'invalid'",
		},
		{
			name: "invalid API base URL",
			config: config.Config{
				LogLevel:   "info",
				LogFormat:  "json",
				APIBaseURL: "://invalid-url",
			},
			wantErr: true,
			errMsg:  "invalid API base URL",
		},
		{
			name: "invalid HTTP address format",
			config: config.Config{
				HTTPAddr:   "localhost8080",
				LogLevel:   "info",
				LogFormat:  "json",
				APIBaseURL: "https://api.example.com",
			},
			wantErr: true,
			errMsg:  "invalid HTTP address 'localhost8080'",
		},
		{
			name: "case insensitive log level - uppercase",
			config: config.Config{
				LogLevel:   "INFO",
				LogFormat:  "json",
				APIBaseURL: "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "case insensitive log format - uppercase",
			config: config.Config{
				LogLevel:   "info",
				LogFormat:  "JSON",
				APIBaseURL: "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "empty API base URL is valid",
			config: config.Config{
				LogLevel:   "info",
				LogFormat:  "json",
				APIBaseURL: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Config.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Config.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else if err != nil {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := getLoadTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables for this test
			setEnvironmentVariables(t, tt.envVars)

			got, err := config.Load(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Load() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			validateLoadResult(t, got, tt.want)
		})
	}
}

func TestLoadWithEmptyArgs(t *testing.T) {
	// This test verifies that Load() works when called with empty arguments
	// instead of using the default command line arguments behavior

	// Clean environment - using t.Setenv with empty strings ensures clean state
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("LOG_FORMAT", "")
	t.Setenv("API_KEY", "")
	t.Setenv("API_BASE_URL", "")

	// Pass empty slice to avoid parsing test runner flags
	cfg, err := config.Load([]string{})
	if err != nil {
		t.Errorf("Load() with empty args should not error, got %v", err)
	}

	// Should have default values
	if cfg.LogLevel != "info" {
		t.Errorf("Load() LogLevel = %v, want info", cfg.LogLevel)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("Load() LogFormat = %v, want json", cfg.LogFormat)
	}
	if cfg.APIBaseURL != "https://api.gsa.gov/analytics/dap/v2" {
		t.Errorf("Load() APIBaseURL = %v, want https://api.gsa.gov/analytics/dap/v2", cfg.APIBaseURL)
	}
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// setEnvironmentVariables sets environment variables for a test case.
func setEnvironmentVariables(t *testing.T, envVars map[string]string) {
	t.Helper()

	// Always clear all config-related environment variables first
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("LOG_FORMAT", "")
	t.Setenv("API_KEY", "")
	t.Setenv("API_BASE_URL", "")

	// Then set the specific ones for this test
	for key, value := range envVars {
		t.Setenv(key, value)
	}
}

// validateLoadResult validates the config returned by Load matches expectations.
func validateLoadResult(t *testing.T, got *config.Config, want config.Config) {
	t.Helper()
	if got.HTTPAddr != want.HTTPAddr {
		t.Errorf("Load() HTTPAddr = %v, want %v", got.HTTPAddr, want.HTTPAddr)
	}
	if got.LogLevel != want.LogLevel {
		t.Errorf("Load() LogLevel = %v, want %v", got.LogLevel, want.LogLevel)
	}
	if got.LogFormat != want.LogFormat {
		t.Errorf("Load() LogFormat = %v, want %v", got.LogFormat, want.LogFormat)
	}
	if got.APIKey != want.APIKey {
		t.Errorf("Load() APIKey = %v, want %v", got.APIKey, want.APIKey)
	}
	if got.APIBaseURL != want.APIBaseURL {
		t.Errorf("Load() APIBaseURL = %v, want %v", got.APIBaseURL, want.APIBaseURL)
	}
}

// getLoadTestCases returns the test cases for TestLoad.
func getLoadTestCases() []struct {
	name    string
	args    []string
	envVars map[string]string
	want    config.Config
	wantErr bool
	errMsg  string
} {
	return []struct {
		name    string
		args    []string
		envVars map[string]string
		want    config.Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "default values with no args or env vars",
			args: []string{},
			want: config.Config{
				HTTPAddr:   "",
				LogLevel:   "info",
				LogFormat:  "json",
				APIKey:     "",
				APIBaseURL: "https://api.gsa.gov/analytics/dap/v2",
			},
			wantErr: false,
		},
		{
			name: "command line args override defaults",
			args: []string{"--http", "localhost:9090", "--log-level", "debug", "--log-format", "text"},
			want: config.Config{
				HTTPAddr:   "localhost:9090",
				LogLevel:   "debug",
				LogFormat:  "text",
				APIKey:     "",
				APIBaseURL: "https://api.gsa.gov/analytics/dap/v2",
			},
			wantErr: false,
		},
		{
			name: "environment variables override defaults",
			args: []string{},
			envVars: map[string]string{
				"HTTP_ADDR":    "0.0.0.0:8080",
				"LOG_LEVEL":    "warn",
				"LOG_FORMAT":   "text",
				"API_KEY":      "secret-key",
				"API_BASE_URL": "https://custom.api.com",
			},
			want: config.Config{
				HTTPAddr:   "0.0.0.0:8080",
				LogLevel:   "warn",
				LogFormat:  "text",
				APIKey:     "secret-key",
				APIBaseURL: "https://custom.api.com",
			},
			wantErr: false,
		},
		{
			name: "command line args override environment variables",
			args: []string{"--http", "127.0.0.1:7070", "--log-level", "error"},
			envVars: map[string]string{
				"HTTP_ADDR": "0.0.0.0:8080",
				"LOG_LEVEL": "debug",
				"API_KEY":   "env-key",
			},
			want: config.Config{
				HTTPAddr:   "127.0.0.1:7070",
				LogLevel:   "error",
				LogFormat:  "json",
				APIKey:     "env-key",
				APIBaseURL: "https://api.gsa.gov/analytics/dap/v2",
			},
			wantErr: false,
		},
		{
			name:    "invalid log level returns error",
			args:    []string{"--log-level", "invalid"},
			want:    config.Config{},
			wantErr: true,
			errMsg:  "invalid log level",
		},
		{
			name:    "invalid flag returns error",
			args:    []string{"--invalid-flag", "value"},
			want:    config.Config{},
			wantErr: true,
			errMsg:  "failed to parse flags",
		},
		{
			name: "custom API base URL",
			args: []string{"--api-base-url", "https://my-api.com/v1"},
			want: config.Config{
				HTTPAddr:   "",
				LogLevel:   "info",
				LogFormat:  "json",
				APIKey:     "",
				APIBaseURL: "https://my-api.com/v1",
			},
			wantErr: false,
		},
	}
}
