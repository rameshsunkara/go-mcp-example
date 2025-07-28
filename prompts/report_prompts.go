package prompts

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ReportPrompts handles report-related prompt operations.
type ReportPrompts struct {
	logger *slog.Logger
}

// NewReportPrompts creates a new ReportPrompts with the provided logger.
func NewReportPrompts(logger *slog.Logger) *ReportPrompts {
	return &ReportPrompts{
		logger: logger,
	}
}

// AnalyzeTrafficPrompt provides guidance for analyzing traffic data.
func (rp *ReportPrompts) AnalyzeTrafficPrompt(_ context.Context, _ *mcp.ServerSession,
	params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	dateRange := params.Arguments["date_range"]
	if dateRange == "" {
		dateRange = "last 30 days"
	}

	rp.logger.Info("Processing analyze traffic prompt", "date_range", dateRange)

	return &mcp.GetPromptResult{
		Description: "Analyze website traffic patterns and trends",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: `Analyze the website traffic data for ` + dateRange + `. Please:

1. Use get_report("traffic") to fetch overall traffic metrics
2. Use get_report("devices") to understand device usage patterns
3. Use get_report("browsers") to see browser preferences
4. Use get_report("top-pages") to identify most popular content

Provide insights on:
- Traffic trends and patterns
- User behavior and preferences
- Device and browser usage
- Content performance
- Recommendations for improvement`,
				},
			},
		},
	}, nil
}

// CompareReportsPrompt helps compare different report types.
func (rp *ReportPrompts) CompareReportsPrompt(_ context.Context, _ *mcp.ServerSession,
	params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	report1 := params.Arguments["report1"]
	report2 := params.Arguments["report2"]

	if report1 == "" {
		report1 = "devices"
	}
	if report2 == "" {
		report2 = "browsers"
	}

	rp.logger.Info("Processing compare reports prompt", "report1", report1, "report2", report2)

	return &mcp.GetPromptResult{
		Description: "Compare and analyze two different report types",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: `Compare ` + report1 + ` and ` + report2 + ` reports. Please:

1. Use get_report("` + report1 + `") to fetch the first report
2. Use get_report("` + report2 + `") to fetch the second report
3. Analyze the data from both reports

Provide a comparative analysis including:
- Key metrics from each report
- Trends and patterns observed
- Correlations between the two data sets
- Actionable insights and recommendations
- Data visualization suggestions`,
				},
			},
		},
	}, nil
}

// MonthlyReportPrompt generates comprehensive monthly analytics.
func (rp *ReportPrompts) MonthlyReportPrompt(_ context.Context, _ *mcp.ServerSession,
	params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	month := params.Arguments["month"]
	year := params.Arguments["year"]

	if month == "" || year == "" {
		month = "01"
		year = "2024"
	}

	rp.logger.Info("Processing monthly report prompt", "month", month, "year", year)

	return &mcp.GetPromptResult{
		Description: "Generate comprehensive monthly analytics report",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: buildMonthlyReportText(month, year),
				},
			},
		},
	}, nil
}

// buildMonthlyReportText builds the monthly report prompt text.
func buildMonthlyReportText(month, year string) string {
	dateRange := `after="` + year + `-` + month + `-01", before="` + year + `-` + month + `-31"`

	return `Generate a comprehensive monthly analytics report for ` + month + `/` + year + `. Please:

1. Use get_report("traffic", ` + dateRange + `) for traffic data
2. Use get_report("devices", ` + dateRange + `) for device breakdown
3. Use get_report("top-pages", ` + dateRange + `) for content performance
4. Use get_report("countries", ` + dateRange + `) for geographic data

Create a comprehensive report with:
- Executive summary of key metrics
- Traffic trends and growth patterns
- User demographics and behavior
- Content performance analysis
- Geographic distribution insights
- Month-over-month comparisons (if available)
- Recommendations for the next month`
}

// RealTimeInsightsPrompt provides real-time analytics guidance.
func (rp *ReportPrompts) RealTimeInsightsPrompt(_ context.Context, _ *mcp.ServerSession,
	_ *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	rp.logger.Info("Processing real-time insights prompt")

	return &mcp.GetPromptResult{
		Description: "Get real-time website analytics and insights",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: `Provide real-time analytics insights. Please:

1. Use get_report("realtime") to get current active users
2. Use get_report("traffic") to get recent traffic trends
3. Use get_report("top-pages") to see what content is currently popular

Analyze and provide:
- Current website activity levels
- Real-time user engagement
- Popular content right now
- Traffic patterns compared to historical data
- Immediate optimization opportunities
- Alert-worthy trends or anomalies`,
				},
			},
		},
	}, nil
}
