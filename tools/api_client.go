package tools

import (
	"net/http"
	"strings"
)

// APIClient configuration and utilities for making API requests.
type APIClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient HTTPClientInterface
}

// HTTPClientInterface defines the interface for HTTP clients (for testing).
type HTTPClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewAPIClient creates a new API client with the given base URL and API key.
func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		BaseURL:    strings.TrimSuffix(baseURL, "/"),
		APIKey:     apiKey,
		HTTPClient: &http.Client{}, // Default HTTP client
	}
}

// NewAPIClientWithHTTPClient creates a new API client with a custom HTTP client (for testing).
func NewAPIClientWithHTTPClient(baseURL, apiKey string, httpClient HTTPClientInterface) *APIClient {
	return &APIClient{
		BaseURL:    strings.TrimSuffix(baseURL, "/"),
		APIKey:     apiKey,
		HTTPClient: httpClient,
	}
}

// DoRequest makes an HTTP request with the configured headers.
func (c *APIClient) DoRequest(req *http.Request) (*http.Response, error) {
	// Add standard headers
	headers := c.HTTPHeaders()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.HTTPClient.Do(req)
}

// HTTPHeaders returns the HTTP headers needed for API requests.
func (c *APIClient) HTTPHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	if c.APIKey != "" {
		headers["X-API-KEY"] = c.APIKey
	}

	return headers
}
