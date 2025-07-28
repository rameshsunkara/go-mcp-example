package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rameshsunkara/go-mcp-example/config"
	"github.com/rameshsunkara/go-mcp-example/models"
)

// ReportsTool handles analytics report fetching operations.
type ReportsTool struct {
	logger    *slog.Logger
	config    *config.Config
	apiClient *APIClient
}

// NewReportsTool creates a new ReportsTool with the provided logger, config, and API client.
func NewReportsTool(logger *slog.Logger, cfg *config.Config, apiClient *APIClient) *ReportsTool {
	return &ReportsTool{
		logger:    logger,
		config:    cfg,
		apiClient: apiClient,
	}
}

// GetReport implements the get_report tool.
func (rt *ReportsTool) GetReport(ctx context.Context, _ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[models.ReportArgs]) (*mcp.CallToolResultFor[struct{}], error) {
	args := params.Arguments

	rt.logger.InfoContext(ctx, "Processing get_report tool call",
		"report_name", args.ReportName,
		"limit", args.Limit)

	// Validate report type
	reportType := models.ReportType(args.ReportName)
	if !reportType.IsValid() {
		return nil, fmt.Errorf("invalid report type '%s'. Valid types: %v", args.ReportName, models.GetAllReportTypes())
	}

	// Build request parameters
	reportParams := models.ReportParams{
		Limit:  args.Limit,
		Page:   args.Page,
		After:  args.After,
		Before: args.Before,
	}

	// Set defaults if not provided
	if reportParams.Limit == 0 {
		reportParams.Limit = 1000
	}
	if reportParams.Page == 0 {
		reportParams.Page = 1
	}

	// Validate parameters
	if err := reportParams.Validate(); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// Build the API URL
	apiURL, err := rt.buildReportsURL(args.ReportName, reportParams)
	if err != nil {
		return nil, fmt.Errorf("failed to build API URL: %w", err)
	}

	rt.logger.InfoContext(ctx, "Making API request", "url", apiURL)

	// Make HTTP request
	reports, fetchErr := rt.fetchReports(ctx, apiURL)
	if fetchErr != nil {
		// All errors are returned as MCP errors for consistent user experience
		result := &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Request failed: " + fetchErr.Error()},
			},
			IsError: true,
		}
		return result, nil //nolint:nilerr // MCP tools return nil error when IsError is true
	}

	// Check if no data was returned
	if len(reports) == 0 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No data found for report: " + args.ReportName},
			},
			IsError: true,
		}, nil
	}

	// Format response
	response := models.ReportResponse{
		Data: reports,
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Analytics Report: %s\n\nFound %d records:\n\n%s",
					args.ReportName, len(reports), string(responseJSON)),
			},
		},
	}, nil
}

// buildReportsURL builds a complete URL for fetching report data.
func (rt *ReportsTool) buildReportsURL(reportName string, params models.ReportParams) (string, error) {
	// Build base URL
	baseURL := rt.apiClient.BaseURL + "/reports/" + reportName + "/data"

	// Parse URL to add query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Add query parameters
	q := u.Query()

	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}

	if params.After != "" {
		q.Set("after", params.After)
	}

	if params.Before != "" {
		q.Set("before", params.Before)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// fetchReports makes the HTTP request to fetch analytics data.
func (rt *ReportsTool) fetchReports(ctx context.Context, apiURL string) ([]models.Reports, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make request using APIClient
	resp, err := rt.apiClient.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		rt.logger.ErrorContext(ctx, "API request failed",
			"status_code", resp.StatusCode,
			"response", string(body))

		// Return formatted error message
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var reports []models.Reports
	if parseErr := json.Unmarshal(body, &reports); parseErr != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", parseErr)
	}

	rt.logger.InfoContext(ctx, "Successfully fetched reports", "count", len(reports))
	return reports, nil
}
