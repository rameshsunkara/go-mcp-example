# VS Code MCP Configuration Guide

This guide helps you set up the Go MCP Example server with VS Code and the MCP extension.

## Prerequisites

- VS Code installed on your system
- VS Code MCP extension installed
- Go 1.24+ installed on your system
- US Data Analytics Program API key

## Getting Your API Key

1. Visit the [US Data Analytics Program API documentation](https://open.gsa.gov/api/dap/)
2. Sign up for an API key if you don't have one
3. Make note of your API key for configuration

## Installation Steps

### 1. Install VS Code MCP Extension

1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X / Cmd+Shift+X)
3. Search for "MCP" and install the official MCP extension
4. Restart VS Code if required

### 2. Set Up the Project

1. **Clone and build the project:**

   ```bash
   git clone https://github.com/rameshsunkara/go-mcp-example.git
   cd go-mcp-example
   make build
   ```

2. **Copy the VS Code configuration:**

   ```bash
   # Copy the MCP configuration to your .vscode folder
   cp docs/vscode/mcp.json .vscode/
   cp docs/vscode/settings.json .vscode/
   ```

3. **Update the configuration:**
   
   Edit `.vscode/mcp.json` and update the `command` path to point to your compiled binary.

### 3. Configure API Access

When you open the project in VS Code, the MCP extension will prompt you for:

- **US Data Gov API Key**: Your API key for accessing analytics data
- **US Data Gov Base URL**: API base URL (defaults to official endpoint)
- **Log Level**: Logging verbosity (info, debug, warn, error)
- **Log Format**: Output format (json, text)

## Configuration Files

### mcp.json

This file defines the MCP server configuration with input prompts:

```jsonc
{
  "inputs": [
    {
      "type": "promptString",
      "id": "us-data-gov-api-key",
      "description": "US Data Gov API Key",
      "password": true
    },
    // ... more inputs
  ],
  "servers": {
    "go-mcp-example-stdio": {
      "type": "stdio",
      "command": "/path/to/go-mcp-example/go-mcp-example",
      "env": {
        "API_KEY": "${input:us-data-gov-api-key}",
        "API_BASE_URL": "${input:us-data-gov-base-url}",
        "LOG_LEVEL": "${input:log_level}",
        "LOG_FORMAT": "${input:log_format}"
      }
    }
  }
}
```

### settings.json

VS Code workspace settings for optimal development experience:

```json
{
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v", "-race"],
  "files.exclude": {
    "**/.git": true,
    "**/coverage.out": true,
    "**/go-mcp-example": true
  }
}
```

## Usage

### Available Tools

Once configured, you'll have access to:

- **get_report**: Fetch analytics data from the Digital Analytics Program
- **Interactive Prompts**: Guided analytics workflows
- **Resources**: Documentation and examples

### Example Commands

In VS Code with the MCP extension:

1. **Fetch device statistics:**
   ```
   Use the get_report tool to get device usage data for the last 30 days
   ```

2. **Get traffic analysis:**
   ```
   Show me traffic trends using the analyze-traffic prompt
   ```

3. **Generate monthly report:**
   ```
   Create a monthly analytics report for January 2024
   ```

## Development Mode

For development, you can also run the server directly:

### Option 1: Using the Built Binary

Update your `.vscode/mcp.json` to point to the compiled binary:

```jsonc
{
  "servers": {
    "go-mcp-example-stdio": {
      "command": "/full/path/to/go-mcp-example/go-mcp-example",
      // ... rest of config
    }
  }
}
```

### Option 2: Using Go Run

For development iterations, you can use `go run`:

```jsonc
{
  "servers": {
    "go-mcp-example-dev": {
      "type": "stdio",
      "command": "go",
      "args": ["run", "main.go"],
      "cwd": "/path/to/go-mcp-example",
      "env": {
        // ... environment variables
      }
    }
  }
}
```

## Troubleshooting

### Common Issues

1. **MCP extension not loading:**
   - Ensure the MCP extension is properly installed
   - Restart VS Code
   - Check the VS Code output panel for MCP logs

2. **Server not starting:**
   - Verify the binary path in `mcp.json`
   - Check that the binary is executable: `chmod +x go-mcp-example`
   - Review environment variables

3. **API key errors:**
   - Verify your API key is valid
   - Check network connectivity
   - Test the API key with a direct HTTP request

4. **Build errors:**
   - Ensure Go 1.24+ is installed: `go version`
   - Run `go mod tidy` to resolve dependencies
   - Check for any compilation errors: `make build`

### Debug Mode

To enable debug logging:

1. Set the log level to "debug" when prompted by VS Code
2. Or manually edit your environment in `mcp.json`:

```jsonc
{
  "env": {
    "LOG_LEVEL": "debug",
    // ... other vars
  }
}
```

### Viewing Logs

- Check the VS Code Output panel
- Select "MCP" from the dropdown to see MCP-specific logs
- Look for server startup messages and any error details

## Customization

### Adding New Inputs

You can extend the `mcp.json` file to add more configuration options:

```jsonc
{
  "inputs": [
    // ... existing inputs
    {
      "type": "promptString",
      "id": "custom-setting",
      "description": "Your custom setting",
      "default": "default-value"
    }
  ]
}
```

### Multiple Server Configurations

You can run multiple instances with different configurations:

```jsonc
{
  "servers": {
    "go-mcp-example-prod": {
      "type": "stdio",
      "command": "/path/to/binary",
      // ... production config
    },
    "go-mcp-example-dev": {
      "type": "stdio", 
      "command": "go",
      "args": ["run", "main.go"],
      // ... development config
    }
  }
}
```

## Support

- Check the main [README.md](../../README.md) for detailed documentation
- Review the [OpenAPI specification](../../openapi.yaml) for API details
- Consult the [Claude Desktop setup guide](../claude-desktop/README.md) for comparison
- Create an issue on the project repository for bugs or feature requests
