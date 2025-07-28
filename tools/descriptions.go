// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tools

// GetReportToolDescription contains the detailed description for the get_report tool.
const GetReportToolDescription = `Fetch analytics reports from the Digital Analytics Program (DAP) API ` +
	`with optional filtering and pagination.

The DAP provides analytics data for U.S. federal government websites. This tool allows you to ` +
	`retrieve various analytics reports with flexible filtering options.

PARAMETERS:
- report_name (required): The type of report to fetch
- limit (optional): Maximum number of records to return (1-10000, default 1000)
- page (optional): Page number for pagination (default 1, 1-based indexing)
- after (optional): Start date filter in YYYY-MM-DD format
- before (optional): End date filter in YYYY-MM-DD format

AVAILABLE REPORT TYPES:
- "devices": Device types used by visitors (desktop, mobile, tablet)
- "browsers": Browser usage statistics (Chrome, Safari, Firefox, etc.)
- "operating-systems": Operating system statistics (Windows, macOS, iOS, etc.)
- "languages": Language preferences of visitors
- "countries": Geographic breakdown by country
- "cities": Geographic breakdown by city
- "traffic": Traffic volume and trends over time
- "top-pages": Most visited pages and their metrics
- "downloads": File download statistics and popular downloads
- "realtime": Real-time active user statistics
- "traffic-sources": Traffic source analysis (direct, referral, search, etc.)
- "domains": Analytics by domain for multi-domain agencies
- "agencies": Analytics aggregated by government agency

EXAMPLES:
- get_report("devices") - Get device statistics with default settings
- get_report("browsers", limit=50) - Get browser stats limited to 50 results
- get_report("traffic", after="2024-01-01", before="2024-01-31") - Get traffic for January 2024
- get_report("top-pages", page=2, limit=100) - Get second page of top pages (100 per page)
- get_report("realtime") - Get current active users

RESPONSE FORMAT:
Returns JSON data containing analytics metrics. The response structure varies by report type but ` +
	`typically includes numerical metrics (visits, users, pageviews), categorical data (device types, ` +
	`browser names), time-series data, geographic information, and behavioral metrics.

NOTE: This tool requires a valid API key to be configured via the API_KEY environment variable. ` +
	`The API provides analytics data for U.S. federal government websites participating in the ` +
	`Digital Analytics Program.`
