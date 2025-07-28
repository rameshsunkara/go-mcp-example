package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/rameshsunkara/go-mcp-example/log"
)

func TestParseLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{
			name:     "debug level lowercase",
			input:    "debug",
			expected: slog.LevelDebug,
		},
		{
			name:     "debug level uppercase",
			input:    "DEBUG",
			expected: slog.LevelDebug,
		},
		{
			name:     "info level lowercase",
			input:    "info",
			expected: slog.LevelInfo,
		},
		{
			name:     "info level uppercase",
			input:    "INFO",
			expected: slog.LevelInfo,
		},
		{
			name:     "warn level lowercase",
			input:    "warn",
			expected: slog.LevelWarn,
		},
		{
			name:     "warn level uppercase",
			input:    "WARN",
			expected: slog.LevelWarn,
		},
		{
			name:     "error level lowercase",
			input:    "error",
			expected: slog.LevelError,
		},
		{
			name:     "error level uppercase",
			input:    "ERROR",
			expected: slog.LevelError,
		},
		{
			name:     "mixed case",
			input:    "WaRn",
			expected: slog.LevelWarn,
		},
		{
			name:     "unknown level defaults to info",
			input:    "unknown",
			expected: slog.LevelInfo,
		},
		{
			name:     "empty string defaults to info",
			input:    "",
			expected: slog.LevelInfo,
		},
		{
			name:     "invalid level defaults to info",
			input:    "invalid_level",
			expected: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := log.ParseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNew_LoggerCreation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		level         string
		useTextFormat bool
	}{
		{
			name:          "debug level with JSON format",
			level:         "debug",
			useTextFormat: false,
		},
		{
			name:          "info level with text format",
			level:         "info",
			useTextFormat: true,
		},
		{
			name:          "warn level with JSON format",
			level:         "warn",
			useTextFormat: false,
		},
		{
			name:          "error level with text format",
			level:         "error",
			useTextFormat: true,
		},
		{
			name:          "invalid level defaults to info with JSON format",
			level:         "invalid",
			useTextFormat: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := log.New(tt.level, tt.useTextFormat)
			if logger == nil {
				t.Error("New() returned nil logger")
			}

			// Test that the logger is functional by calling a method
			// We can't easily test the output without redirecting stderr,
			// but we can at least verify the logger doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Logger panicked: %v", r)
				}
			}()

			// Test that the logger methods don't panic
			logger.Info("test message")
		})
	}
}

func TestNew_LoggerLevels(t *testing.T) {
	tests := getLoggerLevelTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture output
			var buf bytes.Buffer

			// Create a custom handler that writes to our buffer
			opts := &slog.HandlerOptions{
				Level: log.ParseLevel(tt.level),
			}
			handler := slog.NewJSONHandler(&buf, opts)
			logger := slog.New(handler)

			// Log a message at the test level
			logger.Log(context.TODO(), tt.testLogLevel, "test message", "key", "value")

			// Check if the message was logged
			output := buf.String()
			hasOutput := len(strings.TrimSpace(output)) > 0

			validateLogOutput(t, tt.shouldLog, hasOutput, output)
		})
	}
}

// getLoggerLevelTestCases returns test cases for logger level testing.
func getLoggerLevelTestCases() []struct {
	name         string
	level        string
	logLevel     slog.Level
	shouldLog    bool
	testLogLevel slog.Level
} {
	return []struct {
		name         string
		level        string
		logLevel     slog.Level
		shouldLog    bool
		testLogLevel slog.Level
	}{
		{
			name:         "debug logger logs debug messages",
			level:        "debug",
			logLevel:     slog.LevelDebug,
			shouldLog:    true,
			testLogLevel: slog.LevelDebug,
		},
		{
			name:         "debug logger logs info messages",
			level:        "debug",
			logLevel:     slog.LevelDebug,
			shouldLog:    true,
			testLogLevel: slog.LevelInfo,
		},
		{
			name:         "info logger ignores debug messages",
			level:        "info",
			logLevel:     slog.LevelInfo,
			shouldLog:    false,
			testLogLevel: slog.LevelDebug,
		},
		{
			name:         "info logger logs info messages",
			level:        "info",
			logLevel:     slog.LevelInfo,
			shouldLog:    true,
			testLogLevel: slog.LevelInfo,
		},
		{
			name:         "warn logger ignores info messages",
			level:        "warn",
			logLevel:     slog.LevelWarn,
			shouldLog:    false,
			testLogLevel: slog.LevelInfo,
		},
		{
			name:         "warn logger logs warn messages",
			level:        "warn",
			logLevel:     slog.LevelWarn,
			shouldLog:    true,
			testLogLevel: slog.LevelWarn,
		},
		{
			name:         "error logger ignores warn messages",
			level:        "error",
			logLevel:     slog.LevelError,
			shouldLog:    false,
			testLogLevel: slog.LevelWarn,
		},
		{
			name:         "error logger logs error messages",
			level:        "error",
			logLevel:     slog.LevelError,
			shouldLog:    true,
			testLogLevel: slog.LevelError,
		},
	}
}

// validateLogOutput validates whether log output matches expectations.
func validateLogOutput(t *testing.T, shouldLog, hasOutput bool, output string) {
	t.Helper()

	if shouldLog && !hasOutput {
		t.Errorf("Expected message to be logged but got no output")
		return
	}

	if !shouldLog && hasOutput {
		t.Errorf("Expected no output but got: %s", output)
		return
	}

	// If we expect output, verify it's valid JSON with expected fields
	if shouldLog && hasOutput {
		validateJSONLogEntry(t, output)
	}
}

// validateJSONLogEntry validates that a log entry is valid JSON with expected fields.
func validateJSONLogEntry(t *testing.T, output string) {
	t.Helper()

	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Errorf("Failed to parse JSON log output: %v", err)
		return
	}

	validateLogEntryFields(t, logEntry)
}

// validateLogEntryFields validates that a log entry has expected fields and values.
func validateLogEntryFields(t *testing.T, logEntry map[string]interface{}) {
	t.Helper()

	// Check for expected fields
	if _, exists := logEntry["time"]; !exists {
		t.Error("Log entry missing 'time' field")
	}
	if _, exists := logEntry["level"]; !exists {
		t.Error("Log entry missing 'level' field")
	}
	if _, exists := logEntry["msg"]; !exists {
		t.Error("Log entry missing 'msg' field")
	}
	if logEntry["msg"] != "test message" {
		t.Errorf("Expected message 'test message', got %v", logEntry["msg"])
	}
	if logEntry["key"] != "value" {
		t.Errorf("Expected key 'value', got %v", logEntry["key"])
	}
}

func TestNew_TextFormat(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a custom text handler that writes to our buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(&buf, opts)
	logger := slog.New(handler)

	// Log a message
	logger.Info("test message", "key", "value")

	// Check that the output is in text format (not JSON)
	output := buf.String()
	if len(strings.TrimSpace(output)) == 0 {
		t.Error("Expected text output but got no output")
	}

	// Text format should not be valid JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err == nil {
		t.Error("Expected text format but got valid JSON")
	}

	// Text format should contain the message and key-value pairs
	if !strings.Contains(output, "test message") {
		t.Error("Text output should contain the log message")
	}
	if !strings.Contains(output, "key=value") {
		t.Error("Text output should contain key-value pairs")
	}
}

func TestNew_JSONFormat(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a custom JSON handler that writes to our buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(&buf, opts)
	logger := slog.New(handler)

	// Log a message
	logger.Info("test message", "key", "value")

	// Check that the output is in JSON format
	output := buf.String()
	if len(strings.TrimSpace(output)) == 0 {
		t.Error("Expected JSON output but got no output")
	}

	// Should be valid JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Errorf("Expected valid JSON but got parse error: %v", err)
	}

	// Check JSON structure
	if logEntry["msg"] != "test message" {
		t.Errorf("Expected message 'test message', got %v", logEntry["msg"])
	}
	if logEntry["key"] != "value" {
		t.Errorf("Expected key 'value', got %v", logEntry["key"])
	}
}
