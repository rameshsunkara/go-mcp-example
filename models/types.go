// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package models

// ReportType represents the different types of reports available.
type ReportType string

const (
	// Common report types based on the API documentation.
	ReportTypeDevices          ReportType = "devices"
	ReportTypeBrowsers         ReportType = "browsers"
	ReportTypeOperatingSystems ReportType = "operating-systems"
	ReportTypeLanguages        ReportType = "languages"
	ReportTypeCountries        ReportType = "countries"
	ReportTypeCities           ReportType = "cities"
	ReportTypeTraffic          ReportType = "traffic"
	ReportTypeTopPages         ReportType = "top-pages"
	ReportTypeDownloads        ReportType = "downloads"
	ReportTypeActiveUsers      ReportType = "realtime"
	ReportTypeSources          ReportType = "traffic-sources"
	ReportTypeDomains          ReportType = "domains"
	ReportTypeAgencies         ReportType = "agencies"
)

// String returns the string representation of the report type.
func (rt ReportType) String() string {
	return string(rt)
}

// IsValid checks if the report type is valid.
func (rt ReportType) IsValid() bool {
	switch rt {
	case ReportTypeDevices, ReportTypeBrowsers, ReportTypeOperatingSystems,
		ReportTypeLanguages, ReportTypeCountries, ReportTypeCities,
		ReportTypeTraffic, ReportTypeTopPages, ReportTypeDownloads,
		ReportTypeActiveUsers, ReportTypeSources, ReportTypeDomains,
		ReportTypeAgencies:
		return true
	default:
		return false
	}
}

// GetAllReportTypes returns all available report types.
func GetAllReportTypes() []ReportType {
	return []ReportType{
		ReportTypeDevices,
		ReportTypeBrowsers,
		ReportTypeOperatingSystems,
		ReportTypeLanguages,
		ReportTypeCountries,
		ReportTypeCities,
		ReportTypeTraffic,
		ReportTypeTopPages,
		ReportTypeDownloads,
		ReportTypeActiveUsers,
		ReportTypeSources,
		ReportTypeDomains,
		ReportTypeAgencies,
	}
}

// ReportArgs represents the arguments for fetching a report.
type ReportArgs struct {
	ReportName string `json:"report_name" jsonschema:"required" jsonschema_description:"Name of the report"`
	Limit      int    `json:"limit,omitempty" jsonschema_description:"Limit results (1-10000, default 1000)"`
	Page       int    `json:"page,omitempty" jsonschema_description:"Page number (default 1)"`
	After      string `json:"after,omitempty" jsonschema_description:"Start date (YYYY-MM-DD format)"`
	Before     string `json:"before,omitempty" jsonschema_description:"End date (YYYY-MM-DD format)"`
}
