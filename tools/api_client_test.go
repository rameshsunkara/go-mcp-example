package tools_test

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rameshsunkara/go-mcp-example/tools"
)

// MockHTTPClient is a mock implementation of HTTPClientInterface for testing.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewAPIClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		baseURL     string
		apiKey      string
		expectedURL string
		expectedKey string
	}{
		{
			name:        "basic URL without trailing slash",
			baseURL:     "https://api.example.com",
			apiKey:      "test-key",
			expectedURL: "https://api.example.com",
			expectedKey: "test-key",
		},
		{
			name:        "URL with trailing slash",
			baseURL:     "https://api.example.com/",
			apiKey:      "test-key",
			expectedURL: "https://api.example.com",
			expectedKey: "test-key",
		},
		{
			name:        "URL with multiple trailing slashes",
			baseURL:     "https://api.example.com///",
			apiKey:      "test-key",
			expectedURL: "https://api.example.com//",
			expectedKey: "test-key",
		},
		{
			name:        "empty API key",
			baseURL:     "https://api.example.com",
			apiKey:      "",
			expectedURL: "https://api.example.com",
			expectedKey: "",
		},
		{
			name:        "complex URL with path",
			baseURL:     "https://api.example.com/v1/analytics/",
			apiKey:      "secret-key-123",
			expectedURL: "https://api.example.com/v1/analytics",
			expectedKey: "secret-key-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := tools.NewAPIClient(tt.baseURL, tt.apiKey)

			if client == nil {
				t.Fatal("NewAPIClient returned nil")
			}

			if client.BaseURL != tt.expectedURL {
				t.Errorf("BaseURL = %v, want %v", client.BaseURL, tt.expectedURL)
			}

			if client.APIKey != tt.expectedKey {
				t.Errorf("APIKey = %v, want %v", client.APIKey, tt.expectedKey)
			}

			if client.HTTPClient == nil {
				t.Error("HTTPClient should not be nil")
			}
		})
	}
}

func TestNewAPIClientWithHTTPClient(t *testing.T) {
	t.Parallel()

	mockClient := &MockHTTPClient{}
	baseURL := "https://api.example.com/"
	apiKey := "test-key"

	client := tools.NewAPIClientWithHTTPClient(baseURL, apiKey, mockClient)

	if client == nil {
		t.Fatal("NewAPIClientWithHTTPClient returned nil")
	}

	if client.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL = %v, want %v", client.BaseURL, "https://api.example.com")
	}

	if client.APIKey != apiKey {
		t.Errorf("APIKey = %v, want %v", client.APIKey, apiKey)
	}

	if client.HTTPClient != mockClient {
		t.Error("HTTPClient should be the provided mock client")
	}
}

func TestAPIClient_HTTPHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		apiKey          string
		expectedHeaders map[string]string
	}{
		{
			name:   "with API key",
			apiKey: "test-api-key",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
				"X-API-KEY":    "test-api-key",
			},
		},
		{
			name:   "without API key",
			apiKey: "",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		},
		{
			name:   "with complex API key",
			apiKey: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
				"X-API-KEY":    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := tools.NewAPIClient("https://api.example.com", tt.apiKey)
			headers := client.HTTPHeaders()

			// Check that all expected headers are present
			for key, expectedValue := range tt.expectedHeaders {
				if actualValue, exists := headers[key]; !exists {
					t.Errorf("Expected header %s not found", key)
				} else if actualValue != expectedValue {
					t.Errorf("Header %s = %v, want %v", key, actualValue, expectedValue)
				}
			}

			// Check that no unexpected headers are present
			for key := range headers {
				if _, expected := tt.expectedHeaders[key]; !expected {
					t.Errorf("Unexpected header found: %s = %v", key, headers[key])
				}
			}

			// Verify the total number of headers
			if len(headers) != len(tt.expectedHeaders) {
				t.Errorf("Got %d headers, want %d", len(headers), len(tt.expectedHeaders))
			}
		})
	}
}

func TestAPIClient_DoRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		apiKey          string
		requestURL      string
		mockResponse    *http.Response
		mockError       error
		expectedHeaders map[string]string
		expectError     bool
	}{
		{
			name:       "successful request with API key",
			apiKey:     "test-key",
			requestURL: "https://api.example.com/test",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
			},
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
				"X-API-KEY":    "test-key",
			},
			expectError: false,
		},
		{
			name:       "successful request without API key",
			apiKey:     "",
			requestURL: "https://api.example.com/test",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
			},
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
			expectError: false,
		},
		{
			name:        "request with HTTP client error",
			apiKey:      "test-key",
			requestURL:  "https://api.example.com/test",
			mockError:   &mockError{message: "network error"},
			expectError: true,
		},
		{
			name:       "request with 404 response",
			apiKey:     "test-key",
			requestURL: "https://api.example.com/notfound",
			mockResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader(`{"error": "not found"}`)),
			},
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
				"X-API-KEY":    "test-key",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testDoRequestCase(t, tt)
		})
	}
}

// testDoRequestCase handles individual test cases for DoRequest to reduce cognitive complexity.
func testDoRequestCase(t *testing.T, tt struct {
	name            string
	apiKey          string
	requestURL      string
	mockResponse    *http.Response
	mockError       error
	expectedHeaders map[string]string
	expectError     bool
}) {
	var capturedRequest *http.Request
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedRequest = req
			if tt.mockError != nil {
				return nil, tt.mockError
			}
			return tt.mockResponse, nil
		},
	}

	client := tools.NewAPIClientWithHTTPClient("https://api.example.com", tt.apiKey, mockClient)

	// Create a test request
	req, err := http.NewRequest(http.MethodGet, tt.requestURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Execute the request
	resp, err := client.DoRequest(req)

	// Check error expectation
	if tt.expectError {
		if err == nil {
			t.Error("Expected error but got none")
		}
		return
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if resp == nil {
		t.Error("Response should not be nil")
		return
	}

	validateRequestHeaders(t, capturedRequest, tt.expectedHeaders, tt.apiKey)
	validateResponseStatus(t, resp, tt.mockResponse)
}

// validateRequestHeaders validates that the request has the expected headers.
func validateRequestHeaders(t *testing.T, capturedRequest *http.Request,
	expectedHeaders map[string]string, apiKey string) {
	if capturedRequest == nil {
		t.Error("Request was not captured")
		return
	}

	// Check that all expected headers were set
	for key, expectedValue := range expectedHeaders {
		if actualValue := capturedRequest.Header.Get(key); actualValue != expectedValue {
			t.Errorf("Header %s = %v, want %v", key, actualValue, expectedValue)
		}
	}

	// Verify X-API-KEY header is only present when API key is provided
	apiKeyHeader := capturedRequest.Header.Get("X-API-KEY")
	if apiKey == "" {
		if apiKeyHeader != "" {
			t.Errorf("X-API-KEY header should not be present when API key is empty, got %v", apiKeyHeader)
		}
	} else {
		if apiKeyHeader != apiKey {
			t.Errorf("X-API-KEY header = %v, want %v", apiKeyHeader, apiKey)
		}
	}
}

// validateResponseStatus validates the response status code.
func validateResponseStatus(t *testing.T, resp *http.Response, mockResponse *http.Response) {
	if mockResponse != nil && resp.StatusCode != mockResponse.StatusCode {
		t.Errorf("Status code = %v, want %v", resp.StatusCode, mockResponse.StatusCode)
	}
}

func TestAPIClient_DoRequest_HeaderOverride(t *testing.T) {
	t.Parallel()

	mockClient := &MockHTTPClient{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"success": true}`))),
			}, nil
		},
	}

	client := tools.NewAPIClientWithHTTPClient("https://api.example.com", "test-key", mockClient)

	// Create a request with pre-existing headers
	req, err := http.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set some headers that should be overridden
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Accept", "text/html")
	req.Header.Set("X-API-KEY", "old-key")
	req.Header.Set("Custom-Header", "should-remain")

	// Execute the request
	_, err = client.DoRequest(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify that API client headers override the request headers
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type should be overridden to application/json, got %v", req.Header.Get("Content-Type"))
	}

	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("Accept should be overridden to application/json, got %v", req.Header.Get("Accept"))
	}

	if req.Header.Get("X-API-KEY") != "test-key" {
		t.Errorf("X-API-KEY should be overridden to test-key, got %v", req.Header.Get("X-API-KEY"))
	}

	// Verify that custom headers are preserved
	if req.Header.Get("Custom-Header") != "should-remain" {
		t.Errorf("Custom-Header should be preserved, got %v", req.Header.Get("Custom-Header"))
	}
}

// mockError is a helper for testing error scenarios.
type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

func TestAPIClient_Integration(t *testing.T) {
	t.Parallel()

	// Test the integration between different methods
	baseURL := "https://api.example.com/"
	apiKey := "integration-test-key"

	client := tools.NewAPIClient(baseURL, apiKey)

	// Test that HTTPHeaders returns consistent results
	headers1 := client.HTTPHeaders()
	headers2 := client.HTTPHeaders()

	if len(headers1) != len(headers2) {
		t.Error("HTTPHeaders should return consistent results")
	}

	for key, value := range headers1 {
		if headers2[key] != value {
			t.Errorf("HTTPHeaders inconsistent: %s = %v vs %v", key, value, headers2[key])
		}
	}

	// Verify the BaseURL trimming worked correctly
	if strings.HasSuffix(client.BaseURL, "/") {
		t.Error("BaseURL should have trailing slash trimmed")
	}

	// Verify that the HTTPClient interface is properly set
	if client.HTTPClient == nil {
		t.Error("HTTPClient should not be nil after NewAPIClient")
	}
}
