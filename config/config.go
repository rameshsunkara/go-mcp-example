// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
)

// Config holds all configuration for the MCP server.
type Config struct {
	HTTPAddr   string
	LogLevel   string
	LogFormat  string
	APIKey     string // Secret - should only come from env vars for security, not flags
	APIBaseURL string
}

// GetEnv returns the value of an environment variable or a default value.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Load parses command-line flags (or custom args), validates the configuration, and returns it.
// If no args are provided, it uses os.Args[1:] (command-line arguments)
// If args are provided, it uses those instead (useful for testing).
func Load(args ...[]string) (*Config, error) {
	// Create a new FlagSet to avoid global state
	fs := flag.NewFlagSet("config", flag.ContinueOnError)

	httpAddr := fs.String("http", GetEnv("HTTP_ADDR", ""),
		"HTTP address to listen on (can also use HTTP_ADDR env var), if empty uses stdin/stdout")
	logLevel := fs.String("log-level", GetEnv("LOG_LEVEL", "info"),
		"Log level: debug, info, warn, error (can also use LOG_LEVEL env var)")
	logFormat := fs.String("log-format", GetEnv("LOG_FORMAT", "json"),
		"Log format: json, text (can also use LOG_FORMAT env var)")
	apiBaseURL := fs.String("api-base-url", GetEnv("API_BASE_URL", "https://api.gsa.gov/analytics/dap/v2"),
		"API base URL (can also use API_BASE_URL env var)")

	// Determine which arguments to parse
	var argsToUse []string
	if len(args) > 0 {
		argsToUse = args[0]
	} else {
		argsToUse = os.Args[1:]
	}

	if err := fs.Parse(argsToUse); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	cfg := &Config{
		HTTPAddr:   *httpAddr,
		LogLevel:   *logLevel,
		LogFormat:  *logFormat,
		APIKey:     os.Getenv("API_KEY"),
		APIBaseURL: *apiBaseURL,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if the configuration values are valid.
func (c *Config) Validate() error {
	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !slices.Contains(validLogLevels, strings.ToLower(c.LogLevel)) {
		return fmt.Errorf("invalid log level '%s', must be one of: %s", c.LogLevel, strings.Join(validLogLevels, ", "))
	}

	// Validate log format
	validLogFormats := []string{"json", "text"}
	if !slices.Contains(validLogFormats, strings.ToLower(c.LogFormat)) {
		return fmt.Errorf("invalid log format '%s', must be one of: %s", c.LogFormat, strings.Join(validLogFormats, ", "))
	}

	// Validate API base URL
	if c.APIBaseURL != "" {
		if _, err := url.Parse(c.APIBaseURL); err != nil {
			return fmt.Errorf("invalid API base URL '%s': %w", c.APIBaseURL, err)
		}
	}

	// Validate HTTP address format if provided
	if c.HTTPAddr != "" {
		// Basic validation - should contain a colon for host:port format
		if !strings.Contains(c.HTTPAddr, ":") {
			return fmt.Errorf("invalid HTTP address '%s', expected format 'host:port' or ':port'", c.HTTPAddr)
		}
	}

	// APIKey validation could be added here if needed
	// For example, checking minimum length, format, etc.

	return nil
}
