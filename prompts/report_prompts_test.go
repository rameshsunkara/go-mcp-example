package prompts_test

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rameshsunkara/go-mcp-example/prompts"
)

func TestNewReportPrompts(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	rp := prompts.NewReportPrompts(logger)

	if rp == nil {
		t.Error("NewReportPrompts() returned nil")
	}
}

func TestAnalyzeTrafficPrompt(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	rp := prompts.NewReportPrompts(logger)

	tests := []struct {
		name        string
		arguments   map[string]string
		expectInMsg string
	}{
		{
			name:        "with date range",
			arguments:   map[string]string{"date_range": "January 2024"},
			expectInMsg: "January 2024",
		},
		{
			name:        "without date range uses default",
			arguments:   map[string]string{},
			expectInMsg: "last 30 days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testAnalyzeTrafficPromptCase(t, rp, tt)
		})
	}
}

// testAnalyzeTrafficPromptCase handles individual test cases for AnalyzeTrafficPrompt.
func testAnalyzeTrafficPromptCase(t *testing.T, rp *prompts.ReportPrompts, tt struct {
	name        string
	arguments   map[string]string
	expectInMsg string
}) {
	params := &mcp.GetPromptParams{
		Arguments: tt.arguments,
	}

	result, err := rp.AnalyzeTrafficPrompt(context.Background(), nil, params)
	if err != nil {
		t.Errorf("AnalyzeTrafficPrompt() error = %v", err)
		return
	}

	validatePromptResult(t, result, "AnalyzeTrafficPrompt")
	validatePromptMessage(t, result, tt.expectInMsg, "AnalyzeTrafficPrompt")
	validateReportCalls(t, result, []string{"traffic", "devices", "browsers", "top-pages"}, "AnalyzeTrafficPrompt")
}

func TestCompareReportsPrompt(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	rp := prompts.NewReportPrompts(logger)

	tests := []struct {
		name        string
		arguments   map[string]string
		expectInMsg []string
	}{
		{
			name:        "with both reports specified",
			arguments:   map[string]string{"report1": "traffic", "report2": "countries"},
			expectInMsg: []string{"traffic", "countries"},
		},
		{
			name:        "with only report1 specified",
			arguments:   map[string]string{"report1": "languages"},
			expectInMsg: []string{"languages", "browsers"}, // browsers is default for report2
		},
		{
			name:        "with no reports specified uses defaults",
			arguments:   map[string]string{},
			expectInMsg: []string{"devices", "browsers"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			params := &mcp.GetPromptParams{
				Arguments: tt.arguments,
			}

			result, err := rp.CompareReportsPrompt(context.Background(), nil, params)
			if err != nil {
				t.Errorf("CompareReportsPrompt() error = %v", err)
				return
			}

			if result == nil {
				t.Error("CompareReportsPrompt() returned nil result")
				return
			}

			if len(result.Messages) == 0 {
				t.Error("CompareReportsPrompt() returned no messages")
				return
			}

			content, ok := result.Messages[0].Content.(*mcp.TextContent)
			if !ok {
				t.Error("CompareReportsPrompt() message content is not TextContent")
				return
			}

			for _, expected := range tt.expectInMsg {
				if !strings.Contains(content.Text, expected) {
					t.Errorf("CompareReportsPrompt() message should contain %q, got %q", expected, content.Text)
				}
			}
		})
	}
}

func TestMonthlyReportPrompt(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	rp := prompts.NewReportPrompts(logger)

	tests := []struct {
		name        string
		arguments   map[string]string
		expectInMsg []string
	}{
		{
			name:        "with month and year specified",
			arguments:   map[string]string{"month": "03", "year": "2024"},
			expectInMsg: []string{"03/2024", "2024-03-01", "2024-03-31"},
		},
		{
			name:        "with defaults when not specified",
			arguments:   map[string]string{},
			expectInMsg: []string{"01/2024", "2024-01-01", "2024-01-31"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testMonthlyReportPromptCase(t, rp, tt)
		})
	}
}

// testMonthlyReportPromptCase handles individual test cases for MonthlyReportPrompt.
func testMonthlyReportPromptCase(t *testing.T, rp *prompts.ReportPrompts, tt struct {
	name        string
	arguments   map[string]string
	expectInMsg []string
}) {
	params := &mcp.GetPromptParams{
		Arguments: tt.arguments,
	}

	result, err := rp.MonthlyReportPrompt(context.Background(), nil, params)
	if err != nil {
		t.Errorf("MonthlyReportPrompt() error = %v", err)
		return
	}

	validatePromptResult(t, result, "MonthlyReportPrompt")
	validateMultipleMessages(t, result, tt.expectInMsg, "MonthlyReportPrompt")
	validateReportCalls(t, result, []string{"traffic", "devices", "top-pages", "countries"}, "MonthlyReportPrompt")
}

// validateMultipleMessages validates that the prompt contains multiple expected messages.
func validateMultipleMessages(t *testing.T, result *mcp.GetPromptResult, expectedTexts []string, promptName string) {
	if len(result.Messages) == 0 {
		return // Already validated in validatePromptResult
	}

	content, ok := result.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Errorf("%s message content is not TextContent", promptName)
		return
	}

	for _, expected := range expectedTexts {
		if !strings.Contains(content.Text, expected) {
			t.Errorf("%s message should contain %q, got %q", promptName, expected, content.Text)
		}
	}
}

func TestRealTimeInsightsPrompt(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	rp := prompts.NewReportPrompts(logger)

	params := &mcp.GetPromptParams{
		Arguments: map[string]string{},
	}

	result, err := rp.RealTimeInsightsPrompt(context.Background(), nil, params)
	if err != nil {
		t.Errorf("RealTimeInsightsPrompt() error = %v", err)
		return
	}

	if result == nil {
		t.Error("RealTimeInsightsPrompt() returned nil result")
		return
	}

	if result.Description == "" {
		t.Error("RealTimeInsightsPrompt() returned empty description")
	}

	if len(result.Messages) == 0 {
		t.Error("RealTimeInsightsPrompt() returned no messages")
		return
	}

	content, ok := result.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Error("RealTimeInsightsPrompt() message content is not TextContent")
		return
	}

	// Check that it mentions the expected real-time reports
	expectedReports := []string{"realtime", "traffic", "top-pages"}
	for _, report := range expectedReports {
		if !strings.Contains(content.Text, `get_report("`+report+`")`) {
			t.Errorf("RealTimeInsightsPrompt() should mention get_report(%q)", report)
		}
	}

	// Check that it mentions real-time concepts
	realTimeKeywords := []string{"real-time", "current", "active users"}
	foundKeyword := false
	for _, keyword := range realTimeKeywords {
		if strings.Contains(strings.ToLower(content.Text), strings.ToLower(keyword)) {
			foundKeyword = true
			break
		}
	}
	if !foundKeyword {
		t.Error("RealTimeInsightsPrompt() should mention real-time related keywords")
	}
}

func TestPromptResultStructure(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	rp := prompts.NewReportPrompts(logger)

	// Test all prompt methods to ensure they return properly structured results
	prompts := []struct {
		name   string
		method func(context.Context, *mcp.ServerSession, *mcp.GetPromptParams) (*mcp.GetPromptResult, error)
	}{
		{"AnalyzeTrafficPrompt", rp.AnalyzeTrafficPrompt},
		{"CompareReportsPrompt", rp.CompareReportsPrompt},
		{"MonthlyReportPrompt", rp.MonthlyReportPrompt},
		{"RealTimeInsightsPrompt", rp.RealTimeInsightsPrompt},
	}

	for _, prompt := range prompts {
		t.Run(prompt.name, func(t *testing.T) {
			t.Parallel()

			params := &mcp.GetPromptParams{
				Arguments: map[string]string{},
			}

			result, err := prompt.method(context.Background(), nil, params)
			if err != nil {
				t.Errorf("%s error = %v", prompt.name, err)
				return
			}

			validatePromptResult(t, result, prompt.name)
		})
	}
}

// Helper functions to reduce cognitive complexity

// validatePromptResult validates the basic structure of a prompt result.
func validatePromptResult(t *testing.T, result *mcp.GetPromptResult, promptName string) {
	if result == nil {
		t.Errorf("%s returned nil result", promptName)
		return
	}

	if result.Description == "" {
		t.Errorf("%s returned empty description", promptName)
	}

	if len(result.Messages) == 0 {
		t.Errorf("%s returned no messages", promptName)
	}
}

// validatePromptMessage validates the message content and expected text.
func validatePromptMessage(t *testing.T, result *mcp.GetPromptResult, expectedText, promptName string) {
	if len(result.Messages) == 0 {
		return // Already validated in validatePromptResult
	}

	message := result.Messages[0]
	if message.Role != "user" {
		t.Errorf("%s message role = %v, want user", promptName, message.Role)
	}

	content, ok := message.Content.(*mcp.TextContent)
	if !ok {
		t.Errorf("%s message content is not TextContent", promptName)
		return
	}

	if !strings.Contains(content.Text, expectedText) {
		t.Errorf("%s message text should contain %q, got %q", promptName, expectedText, content.Text)
	}
}

// validateReportCalls validates that the prompt mentions the expected report calls.
func validateReportCalls(t *testing.T, result *mcp.GetPromptResult, expectedReports []string, promptName string) {
	if len(result.Messages) == 0 {
		return // Already validated in validatePromptResult
	}

	content, ok := result.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		return // Already validated in validatePromptMessage
	}

	for _, report := range expectedReports {
		// For monthly reports, check for get_report with the report name (may have additional parameters)
		pattern := `get_report("` + report + `"`
		if !strings.Contains(content.Text, pattern) {
			t.Errorf("%s should mention %s", promptName, pattern)
		}
	}
}
