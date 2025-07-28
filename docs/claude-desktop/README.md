# Claude Desktop Setup Guide

This guide helps you set up the Go MCP Example server with Claude Desktop.

## Prerequisites

- Go 1.24+ installed on your system
- Claude Desktop application
- US Data Analytics Program API key

## Getting Your API Key

1. Visit the [US Data Analytics Program API documentation](https://open.gsa.gov/api/dap/)
2. Sign up for an API key if you don't have one
3. Make note of your API key for configuration

## Installation Options

### Option 1: Using Compiled Binary (Recommended)

1. **Clone and build the project:**
   ```bash
   git clone https://github.com/rameshsunkara/go-mcp-example.git
   cd go-mcp-example
   make build
   ```

2. **Copy the configuration:**
   ```bash
   cp docs/claude-desktop/claude-desktop-config.json /path/to/claude/config/
   ```

3. **Update the configuration:**
   Edit your Claude Desktop configuration file and update:
   - `command`: Path to your compiled `go-mcp-example` binary
   - `API_KEY`: Your actual API key

### Option 2: Using Go Run (Development)

1. **Clone the project:**
   ```bash
   git clone https://github.com/rameshsunkara/go-mcp-example.git
   cd go-mcp-example
   ```

2. **Copy the development configuration:**
   ```bash
   cp docs/claude-desktop/claude-desktop-config-dev.json /path/to/claude/config/
   ```

3. **Update the configuration:**
   Edit your Claude Desktop configuration file and update:
   - `cwd`: Path to your project directory
   - `API_KEY`: Your actual API key

## Configuration File Locations

### macOS

```text
~/Library/Application Support/Claude/claude_desktop_config.json
```

### Windows

```text
%APPDATA%\Claude\claude_desktop_config.json
```

### Linux

```text
~/.config/claude/claude_desktop_config.json
```

## Testing the Setup

1. **Restart Claude Desktop** after updating the configuration

2. **Test the connection** by asking Claude:

   ```text
   Can you use the get_report tool to show me device usage statistics?
   ```

3. **Verify the tools are available:**
   - `get_report` - Fetch analytics data
   - Interactive prompts for traffic analysis and reporting

## Troubleshooting

### Common Issues

1. **"Command not found" error:**
   - Verify the `command` path in your configuration
   - Ensure the binary is executable: `chmod +x go-mcp-example`

2. **"API key error":**
   - Verify your API key is correctly set in the configuration
   - Test your API key with a direct HTTP request

3. **"Go not found" error (development mode):**
   - Ensure Go is installed and in your PATH
   - Verify with: `go version`

4. **Server not starting:**
   - Check the Claude Desktop logs
   - Verify all environment variables are set correctly
   - Ensure port 8080 is available (if using HTTP mode)

### Debug Mode

To enable debug logging, set `LOG_LEVEL` to `debug` in your configuration:

```json
{
  "mcpServers": {
    "go-mcp-example": {
      "command": "/path/to/go-mcp-example",
      "env": {
        "LOG_LEVEL": "debug",
        // ... other config
      }
    }
  }
}
```

## Support

- Check the [README.md](../../README.md) for detailed documentation
- Review the [OpenAPI specification](../../openapi.yaml) for API details
- Create an issue on the project repository for bugs or feature requests
