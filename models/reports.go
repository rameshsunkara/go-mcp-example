// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"
	"time"
)

// Reports represents the analytics report data structure.
type Reports struct {
	// Required fields
	ID           int    `json:"id" jsonschema:"description=Generated unique identifier"`
	ReportName   string `json:"report_name" jsonschema:"description=The name of the data point's report"`
	ReportAgency string `json:"report_agency" jsonschema:"description=The name of the data point's agency"`
	Date         string `json:"date" jsonschema:"description=The date the data in the data point corresponds to"`

	// Optional fields - depend on the report requested
	ActiveVisitors             int     `json:"active_visitors,omitempty"`
	AvgSessionDuration         float64 `json:"avg_session_duration,omitempty"`
	BounceRate                 float64 `json:"bounce_rate,omitempty"`
	Browser                    string  `json:"browser,omitempty"`
	City                       string  `json:"city,omitempty"`
	Country                    string  `json:"country,omitempty"`
	Device                     string  `json:"device,omitempty" jsonschema:"description=the type of device of the visitor"`
	Domain                     string  `json:"domain,omitempty"`
	EventLabel                 string  `json:"event_label,omitempty"`
	FileName                   string  `json:"file_name,omitempty"`
	Hour                       string  `json:"hour,omitempty"`
	LandingPage                string  `json:"landing_page,omitempty"`
	Language                   string  `json:"language,omitempty"`
	LanguageCode               string  `json:"language_code,omitempty"`
	MobileDevice               string  `json:"mobile_device,omitempty"`
	OS                         string  `json:"os,omitempty" jsonschema:"description=the operating system of the visitor"`
	OSVersion                  string  `json:"os_version,omitempty" jsonschema:"description=the operating system version"`
	Page                       string  `json:"page,omitempty" jsonschema:"description=the path of the page visited"`
	PageTitle                  string  `json:"page_title,omitempty"`
	Pageviews                  int     `json:"pageviews,omitempty"`
	PageviewsPerSession        int     `json:"pageviews_per_session,omitempty"`
	ScreenResolution           string  `json:"screen_resolution,omitempty"`
	SessionDefaultChannelGroup string  `json:"session_default_channel_group,omitempty"`
	Source                     string  `json:"source,omitempty"`
	TotalEvents                int     `json:"total_events,omitempty"`
	Users                      int     `json:"users,omitempty"`
	Visits                     int     `json:"visits,omitempty"`
}

// Error represents an API error response.
type Error struct {
	Message        string `json:"message" jsonschema:"description=Error message"`
	RequiredFields string `json:"required_fields,omitempty" jsonschema:"description=Required fields that are missing"`
	Example        string `json:"example,omitempty" jsonschema:"description=Example of correct usage"`
}

// ReportParams represents query parameters for report requests.
type ReportParams struct {
	Limit  int    `json:"limit,omitempty" jsonschema:"description=Limit the number of data points (max 10,000)"`
	Page   int    `json:"page,omitempty" jsonschema:"description=Pages through results (e.g., page=2 for next 1000)"`
	After  string `json:"after,omitempty" jsonschema:"description=Limit results to dates on or after (YYYY-MM-DD)"`
	Before string `json:"before,omitempty" jsonschema:"description=Limit results to dates on or before (YYYY-MM-DD)"`
}

// ValidateReportParams validates the report parameters.
func (p *ReportParams) Validate() error {
	if p.Limit != 0 && (p.Limit < 1 || p.Limit > 10000) {
		return fmt.Errorf("limit must be between 1 and 10000, got %d", p.Limit)
	}

	if p.Page != 0 && p.Page < 1 {
		return fmt.Errorf("page must be >= 1, got %d", p.Page)
	}

	// TODO: Add date format validation for After and Before
	// Expected format: YYYY-MM-DD

	return nil
}

// ReportRequest represents a request for analytics data.
type ReportRequest struct {
	ReportName string       `json:"report_name" jsonschema:"description=Name of the report"`
	AgencyName string       `json:"agency_name,omitempty" jsonschema:"description=Name of the agency"`
	Domain     string       `json:"domain,omitempty" jsonschema:"description=Name of the domain"`
	Parameters ReportParams `json:"parameters" jsonschema:"description=Query parameters"`
}

// ReportResponse represents the response containing analytics data.
type ReportResponse struct {
	Data  []Reports `json:"data" jsonschema:"description=Array of report data"`
	Error *Error    `json:"error,omitempty" jsonschema:"description=Error information if request failed"`
}

// ParseDate helper function to parse the date string into time.Time.
func (r *Reports) ParseDate() (time.Time, error) {
	return time.Parse("2006-01-02", r.Date)
}
