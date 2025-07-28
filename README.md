# Go MCP Server Example

A robust Model Context Protocol (MCP) server implementation in Go, featuring analytics reporting tools with idiomatic Go architecture, comprehensive error handling, and enterprise-grade configurability.

## Overview

This project demonstrates a well-structured MCP server that provides analytics reporting capabilities by interfacing with external APIs. It showcases modern Go development practices including dependency injection, structured logging, configuration management, and comprehensive error handling.

## Features

### MCP Protocol Features

1. **Tools**: Analytics report fetching with configurable parameters
2. **Resources**: Extensible resource management system  
3. **Prompts**: Interactive prompt system for enhanced user experience
4. **Error Handling**: Consistent error responses with proper MCP error formatting

### Architecture Features

1. **Idiomatic Go Structure**: Clean separation of concerns with dedicated packages
2. **Dependency Injection**: Testable components with injectable dependencies
3. **Configuration Management**: Environment-based configuration with validation
4. **Structured Logging**: Context-aware logging with slog
5. **Type Safety**: Generated models from OpenAPI specifications
6. **Security**: Secure handling of API keys and sensitive configuration

### Development Features

1. **OpenAPI Integration**: Auto-generated models from API specifications
2. **Testable HTTP Client**: Injectable HTTP client interface for easy mocking
3. **Comprehensive Error Handling**: All errors returned as MCP-compatible responses
4. **Docker Support**: Containerized deployment with multi-stage builds
5. **VS Code Integration**: Configured for seamless MCP server development

## Project Structure

```text
go-mcp-example/
├── main.go                        # Entry point and MCP server setup
├── config/
│   └── config.go                  # Configuration management with validation
├── log/
│   └── logger.go                  # Structured logging with slog
├── models/
│   ├── types.go                   # Core data types and validation
│   ├── reports.go                 # Report-specific models
│   └── client.go                  # API client models
├── tools/
│   ├── api_client.go              # HTTP client with dependency injection
│   ├── reports_tools.go           # Analytics report fetching tool
│   └── descriptions.go            # Tool descriptions and documentation
├── prompts/
│   └── report_prompts.go          # Interactive prompts for analytics
├── resources/
│   └── example_resource.go        # MCP resources
├── docs/
│   ├── claude-desktop/            # Claude Desktop configuration and setup
│   │   ├── claude-desktop-config.json       # Production configuration
│   │   ├── claude-desktop-config-dev.json   # Development configuration
│   │   └── README.md                        # Claude Desktop setup guide
│   ├── vscode/                    # VS Code configuration and setup
│   │   ├── mcp.json              # MCP server configuration
│   │   ├── settings.json         # VS Code workspace settings
│   │   └── README.md             # VS Code setup guide
│   └── README.md                 # Documentation index
├── .vscode/
│   ├── mcp.json                   # VS Code MCP configuration
│   └── settings.json              # VS Code project settings
├── .env.example                   # Environment variables template
├── .gitignore                     # Git ignore patterns
├── openapi.yaml                   # API specification
├── Dockerfile                     # Multi-stage Docker build
├── Makefile                       # Development automation
└── .env                          # Environment configuration
```

## Architecture Flow

```mermaid
flowchart TD
    A[MCP Client] -->|JSON-RPC| B[MCP Server]
    B --> C{Request Type}
    C -->|Tool Call| D[Tools Package]
    C -->|Resource| E[Resources Package]
    C -->|Prompt| F[Prompts Package]
    
    D --> G[Reports Tool]
    G --> H[API Client]
    H --> I[External API]
    
    G --> J[Config]
    G --> K[Logger]
    
    L[Models] --> G
    L --> H
```

The MCP server handles three types of operations:

1. **Tools**: Execute analytics report fetching with configurable parameters
2. **Resources**: Provide access to static or dynamic resources
3. **Prompts**: Handle interactive prompts for enhanced user experience

## Quick Start

### Prerequisites

- Go 1.24+
- US Data Analytics Program API key ([Get your API key here](https://open.gsa.gov/api/dap/#authentication))
- Make (optional, for convenience commands)
- Docker (optional, for containerized deployment)

### Running the MCP Server

1. **Clone and setup**:

   ```bash
   git clone https://github.com/rameshsunkara/go-mcp-example.git
   cd go-mcp-example
   cp .env.example .env  # Configure your API settings
   ```

   **Important**: Edit the `.env` file and add your API key:

   ```bash
   API_KEY=your-actual-api-key-here
   ```

   Get your API key from: <https://open.gsa.gov/api/dap/#authentication>

2. **Install dependencies**:

   ```bash
   go mod download
   ```

3. **Run the server**:

   ```bash
   # Via stdio (default MCP transport)
   make run
   # OR: go run main.go
   
   # Via HTTP (for debugging)
   make run-http
   # OR: go run main.go --http localhost:8080
   ```

4. **Test with VS Code**:
   - Install the MCP extension in VS Code
   - The server will be auto-configured via `.vscode/mcp.json`
   - See the [VS Code MCP Configuration](#vs-code-mcp-configuration) section below for details

### VS Code MCP Configuration

This project includes a pre-configured VS Code setup in the `docs/vscode/` directory that automatically sets up the MCP server in VS Code.

For detailed setup instructions, see the [VS Code Setup Guide](docs/vscode/README.md).

#### Quick Setup

1. **Build the executable**: `make build`
2. **Copy VS Code configurations**:

   ```bash
   cp docs/vscode/mcp.json .vscode/
   cp docs/vscode/settings.json .vscode/
   ```

3. **Install VS Code MCP Extension**: Install the official MCP extension for VS Code
4. **Configure API Key**: When VS Code loads, it will prompt you to enter your configuration

#### Available Configurations

- **stdio mode**: Default mode using standard input/output
- **HTTP mode**: Debug mode using HTTP transport on localhost:8080

#### VS Code Configuration Files

- **[docs/vscode/mcp.json](docs/vscode/mcp.json)**: MCP server configuration with input prompts
- **[docs/vscode/settings.json](docs/vscode/settings.json)**: VS Code workspace settings optimized for Go development
- Add additional server instances

### Claude Desktop Configuration

This project includes pre-configured Claude Desktop configuration files in the `docs/claude-desktop/` directory for easy integration with Claude Desktop.

#### Configuration Files

Two configuration files are provided:

1. **`docs/claude-desktop/claude-desktop-config.json`** - Production configuration using the compiled binary
2. **`docs/claude-desktop/claude-desktop-config-dev.json`** - Development configuration using `go run`

#### Claude Desktop Quick Setup

For detailed setup instructions, see the [Claude Desktop Setup Guide](docs/claude-desktop/README.md).

**Production Setup:**

1. Build the binary: `make build`
2. Copy `docs/claude-desktop/claude-desktop-config.json` to your Claude Desktop config directory
3. Update the `command` path and `API_KEY` in the configuration
4. Restart Claude Desktop

**Development Setup:**

1. Copy `docs/claude-desktop/claude-desktop-config-dev.json` to your Claude Desktop config directory  
2. Update the `cwd` path and `API_KEY` in the configuration
3. Restart Claude Desktop

### Configuration

Configure via environment variables, `.env` file, or VS Code's `.vscode/mcp.json`:

```bash
# API Configuration
API_BASE_URL=https://api.gsa.gov/analytics/dap/v2
API_KEY=your-secret-api-key-here  # Get your API key: https://open.gsa.gov/api/dap/#authentication

# Logging Configuration  
LOG_LEVEL=info                    # debug, info, warn, error
LOG_FORMAT=json                   # json, text

# Server Configuration (optional)
HTTP_ADDR=localhost:8080          # Enable HTTP transport for debugging
```

### Available Tools

#### get_report - Analytics Report Fetching

The `get_report` tool provides comprehensive access to the [Digital Analytics Program (DAP) API](https://open.gsa.gov/api/dap/#reports), allowing you to fetch various analytics reports for U.S. federal government websites.

**Parameters:**

- `report_name` (required): The type of report to fetch
- `limit` (optional): Maximum number of records (1-10000, default 1000)
- `page` (optional): Page number for pagination (default 1)
- `after` (optional): Start date filter (YYYY-MM-DD format)
- `before` (optional): End date filter (YYYY-MM-DD format)

**Available Report Types:**

| Report Type | Description | Key Metrics |
|-------------|-------------|-------------|
| `devices` | Device usage statistics | Desktop, mobile, tablet breakdowns |
| `browsers` | Browser usage data | Chrome, Safari, Firefox, Edge usage |
| `operating-systems` | OS statistics | Windows, macOS, iOS, Android, Linux |
| `languages` | Language preferences | Visitor language settings |
| `countries` | Geographic data by country | Traffic by country |
| `cities` | Geographic data by city | Traffic by major cities |
| `traffic` | Traffic volume trends | Visits, users, pageviews over time |
| `top-pages` | Most visited pages | Page paths and their metrics |
| `downloads` | File download stats | Popular downloads and metrics |
| `realtime` | Real-time active users | Current active visitor count |
| `traffic-sources` | Traffic source analysis | Direct, referral, search, social |
| `domains` | Multi-domain analytics | Per-domain statistics |
| `agencies` | Agency-level analytics | Government agency breakdowns |

**Usage Examples:**

```bash
# Get device statistics
get_report("devices")

# Get top 50 pages
get_report("top-pages", limit=50)

# Get browser data for January 2024
get_report("browsers", after="2024-01-01", before="2024-01-31")

# Get traffic sources with pagination
get_report("traffic-sources", page=2, limit=100)

# Get real-time active users
get_report("realtime")
```

**Response Data:**

The tool returns structured JSON data with fields varying by report type:

- **Common fields**: `id`, `report_name`, `report_agency`, `date`
- **Metrics**: `visits`, `users`, `pageviews`, `bounce_rate`
- **Categorical**: `device`, `browser`, `os`, `country`, `page`
- **Behavioral**: `avg_session_duration`, `pageviews_per_session`

**Error Handling:**

- Invalid report names are validated and return descriptive errors
- Date format validation ensures YYYY-MM-DD format
- API authentication and rate limiting are handled gracefully
- Network timeouts and connectivity issues are reported clearly

## Troubleshooting

### Common Issues

#### Empty or No Results from API

The Digital Analytics Program (DAP) API may sometimes return empty results or no data for certain queries. This is a known limitation of the current API endpoint. Common scenarios include:

- **Recent dates**: Very recent data (last 24-48 hours) may not be available yet
- **Specific filters**: Certain combinations of filters may not have data
- **Low-traffic periods**: Some reports may be empty during low-traffic periods
- **API maintenance**: The API may be temporarily unavailable or returning limited data

**Workarounds:**

- Try querying data from a few days ago instead of today
- Use broader date ranges to increase the likelihood of finding data
- Check different report types to see if the issue is report-specific
- Try removing optional filters to get broader results

**Future Plans:**

I am aware of these API reliability issues and are evaluating more stable analytics APIs to provide better data consistency and availability. A migration to a more reliable data source is planned for a future release.

### Development Commands

The Makefile provides convenient development commands:

```bash
# Building and Running
make build                         # Build the MCP server binary
make run                           # Run the server via stdio (default MCP transport)
make run-http                      # Run the server via HTTP (for debugging)

# Testing and Quality
make test                          # Run tests with coverage
make coverage                      # Generate and display coverage report
make lint                          # Run the linter
make lint-fix                      # Run the linter and fix issues
make format                        # Format Go code
make tidy                          # Tidy Go modules

# Docker
make docker-build                  # Build Docker image
make docker-run                    # Run containerized server
make docker-clean                  # Clean Docker resources
```

### Docker Deployment

```bash
# Build Docker image
make docker-build

# Run containerized server
make docker-run
```

## Technology Stack

1. **MCP Protocol**: [Model Context Protocol Go SDK](https://github.com/modelcontextprotocol/go-sdk)
2. **Logging**: [slog](https://pkg.go.dev/log/slog) - Structured logging
3. **HTTP Client**: Standard `net/http` with dependency injection
4. **Configuration**: Environment variables with validation
5. **Container**: [Docker](https://www.docker.com/) with multi-stage builds

## Roadmap

- **Enhanced Analytics**: Add support for more analytics endpoints and data sources
- **Caching Layer**: Implement intelligent caching for improved performance
- **Authentication**: Add support for multiple authentication mechanisms
- **Testing**: Comprehensive unit and integration test coverage
- **Monitoring**: Add health checks and metrics collection
- **Documentation**: Auto-generated API documentation from OpenAPI specs

## Contributing

- Feel free to open Pull Requests with improvements
- Create issues for bugs or feature requests  
- Suggestions for architectural improvements are welcome

### What This Is

- A production-ready MCP server template with analytics capabilities
- A showcase of modern Go development practices and patterns
- A foundation for building custom MCP tools and integrations
- A reference implementation for MCP protocol handling in Go

### What This Is Not

- A one-size-fits-all solution (customize it for your specific needs)
