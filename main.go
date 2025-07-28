package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rameshsunkara/go-mcp-example/config"
	"github.com/rameshsunkara/go-mcp-example/log"
	"github.com/rameshsunkara/go-mcp-example/prompts"
	"github.com/rameshsunkara/go-mcp-example/resources"
	"github.com/rameshsunkara/go-mcp-example/tools"
)

const (
	readTimeoutSeconds  = 30
	writeTimeoutSeconds = 30
	idleTimeoutSeconds  = 60
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Create logger with specified level and format
	useTextFormat := cfg.LogFormat == "text"
	logger := log.New(cfg.LogLevel, useTextFormat)

	logger.Info("Starting MCP server",
		"name", "go-mcp-example",
		"log_level", cfg.LogLevel,
		"log_format", cfg.LogFormat)

	server := mcp.NewServer(&mcp.Implementation{Name: "go-mcp-example"}, nil)

	// Create shared API client for all analytics tools
	apiClient := tools.NewAPIClient(cfg.APIBaseURL, cfg.APIKey)

	// Create tools, prompts, and resources with logger, config, and shared API client
	reportsTool := tools.NewReportsTool(logger, cfg, apiClient)
	reportPrompts := prompts.NewReportPrompts(logger)
	resourceHandler := resources.NewResourceHandler(logger)

	// Register tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_report",
		Description: tools.GetReportToolDescription,
	}, reportsTool.GetReport)

	// Register prompts
	server.AddPrompt(&mcp.Prompt{
		Name:        "analyze-traffic",
		Description: "Analyze website traffic patterns and trends with guided data exploration",
	}, reportPrompts.AnalyzeTrafficPrompt)

	server.AddPrompt(&mcp.Prompt{
		Name:        "compare-reports",
		Description: "Compare and analyze two different report types for insights",
	}, reportPrompts.CompareReportsPrompt)

	server.AddPrompt(&mcp.Prompt{
		Name:        "monthly-report",
		Description: "Generate comprehensive monthly analytics report",
	}, reportPrompts.MonthlyReportPrompt)

	server.AddPrompt(&mcp.Prompt{
		Name:        "realtime-insights",
		Description: "Get real-time website analytics and insights",
	}, reportPrompts.RealTimeInsightsPrompt)

	// Register resources
	server.AddResource(&mcp.Resource{
		Name:     "info",
		MIMEType: "text/plain",
		URI:      "embedded:info",
	}, resourceHandler.HandleEmbeddedResource)

	if cfg.HTTPAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		logger.Info("MCP handler starting", "transport", "http", "address", cfg.HTTPAddr)

		// Create HTTP server with timeouts
		httpServer := &http.Server{
			Addr:         cfg.HTTPAddr,
			Handler:      handler,
			ReadTimeout:  readTimeoutSeconds * time.Second,
			WriteTimeout: writeTimeoutSeconds * time.Second,
			IdleTimeout:  idleTimeoutSeconds * time.Second,
		}

		if httpErr := httpServer.ListenAndServe(); httpErr != nil {
			logger.Error("HTTP server failed", "error", httpErr)
		}
	} else {
		logger.Info("MCP handler starting", "transport", "stdio")
		t := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
		if runErr := server.Run(context.Background(), t); runErr != nil {
			logger.Error("Server failed", "error", runErr)
		}
	}
}
