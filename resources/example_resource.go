package resources

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var EmbeddedResources = map[string]string{
	"info": "This is the hello example server.",
}

// ResourceHandler handles resource operations.
type ResourceHandler struct {
	logger *slog.Logger
}

// NewResourceHandler creates a new ResourceHandler with the provided logger.
func NewResourceHandler(logger *slog.Logger) *ResourceHandler {
	return &ResourceHandler{
		logger: logger,
	}
}

// HandleEmbeddedResource implements embedded resource handling.
func (rh *ResourceHandler) HandleEmbeddedResource(_ context.Context, _ *mcp.ServerSession,
	params *mcp.ReadResourceParams) (*mcp.ReadResourceResult, error) {
	rh.logger.Info("Processing resource request", "uri", params.URI)

	u, err := url.Parse(params.URI)
	if err != nil {
		rh.logger.Error("Failed to parse resource URI", "uri", params.URI, "error", err)
		return nil, err
	}
	if u.Scheme != "embedded" {
		rh.logger.Error("Invalid resource scheme", "scheme", u.Scheme, "expected", "embedded")
		return nil, fmt.Errorf("wrong scheme: %q", u.Scheme)
	}
	key := u.Opaque
	text, ok := EmbeddedResources[key]
	if !ok {
		rh.logger.Error("Resource not found", "key", key)
		return nil, fmt.Errorf("no embedded resource named %q", key)
	}

	// You can use context here for:
	// - Timeout handling
	// - Cancellation
	// - Request tracing
	// - API calls with context

	rh.logger.Info("Resource retrieved successfully", "key", key, "length", len(text))
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{URI: params.URI, MIMEType: "text/plain", Text: text},
		},
	}, nil
}
