package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/rameshsunkara/go-mcp-example/models"
)

func TestReportType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		reportType models.ReportType
		expected   string
	}{
		{
			name:       "devices report type",
			reportType: models.ReportTypeDevices,
			expected:   "devices",
		},
		{
			name:       "browsers report type",
			reportType: models.ReportTypeBrowsers,
			expected:   "browsers",
		},
		{
			name:       "operating systems report type",
			reportType: models.ReportTypeOperatingSystems,
			expected:   "operating-systems",
		},
		{
			name:       "languages report type",
			reportType: models.ReportTypeLanguages,
			expected:   "languages",
		},
		{
			name:       "countries report type",
			reportType: models.ReportTypeCountries,
			expected:   "countries",
		},
		{
			name:       "cities report type",
			reportType: models.ReportTypeCities,
			expected:   "cities",
		},
		{
			name:       "traffic report type",
			reportType: models.ReportTypeTraffic,
			expected:   "traffic",
		},
		{
			name:       "top pages report type",
			reportType: models.ReportTypeTopPages,
			expected:   "top-pages",
		},
		{
			name:       "downloads report type",
			reportType: models.ReportTypeDownloads,
			expected:   "downloads",
		},
		{
			name:       "active users report type",
			reportType: models.ReportTypeActiveUsers,
			expected:   "realtime",
		},
		{
			name:       "sources report type",
			reportType: models.ReportTypeSources,
			expected:   "traffic-sources",
		},
		{
			name:       "domains report type",
			reportType: models.ReportTypeDomains,
			expected:   "domains",
		},
		{
			name:       "agencies report type",
			reportType: models.ReportTypeAgencies,
			expected:   "agencies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.reportType.String()
			if result != tt.expected {
				t.Errorf("ReportType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReportType_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		reportType models.ReportType
		expected   bool
	}{
		{
			name:       "valid devices report type",
			reportType: models.ReportTypeDevices,
			expected:   true,
		},
		{
			name:       "valid browsers report type",
			reportType: models.ReportTypeBrowsers,
			expected:   true,
		},
		{
			name:       "valid operating systems report type",
			reportType: models.ReportTypeOperatingSystems,
			expected:   true,
		},
		{
			name:       "valid languages report type",
			reportType: models.ReportTypeLanguages,
			expected:   true,
		},
		{
			name:       "valid countries report type",
			reportType: models.ReportTypeCountries,
			expected:   true,
		},
		{
			name:       "valid cities report type",
			reportType: models.ReportTypeCities,
			expected:   true,
		},
		{
			name:       "valid traffic report type",
			reportType: models.ReportTypeTraffic,
			expected:   true,
		},
		{
			name:       "valid top pages report type",
			reportType: models.ReportTypeTopPages,
			expected:   true,
		},
		{
			name:       "valid downloads report type",
			reportType: models.ReportTypeDownloads,
			expected:   true,
		},
		{
			name:       "valid active users report type",
			reportType: models.ReportTypeActiveUsers,
			expected:   true,
		},
		{
			name:       "valid sources report type",
			reportType: models.ReportTypeSources,
			expected:   true,
		},
		{
			name:       "valid domains report type",
			reportType: models.ReportTypeDomains,
			expected:   true,
		},
		{
			name:       "valid agencies report type",
			reportType: models.ReportTypeAgencies,
			expected:   true,
		},
		{
			name:       "invalid empty report type",
			reportType: "",
			expected:   false,
		},
		{
			name:       "invalid unknown report type",
			reportType: "unknown-report",
			expected:   false,
		},
		{
			name:       "invalid random string report type",
			reportType: "invalid",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.reportType.IsValid()
			if result != tt.expected {
				t.Errorf("ReportType.IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetAllReportTypes(t *testing.T) {
	t.Parallel()

	reportTypes := models.GetAllReportTypes()

	// Test that we get the expected number of report types
	expectedCount := 13
	if len(reportTypes) != expectedCount {
		t.Errorf("GetAllReportTypes() returned %d types, want %d", len(reportTypes), expectedCount)
	}

	// Test that all returned types are valid
	for _, rt := range reportTypes {
		if !rt.IsValid() {
			t.Errorf("GetAllReportTypes() returned invalid report type: %s", rt)
		}
	}

	// Test that specific expected types are included
	expectedTypes := []models.ReportType{
		models.ReportTypeDevices,
		models.ReportTypeBrowsers,
		models.ReportTypeOperatingSystems,
		models.ReportTypeLanguages,
		models.ReportTypeCountries,
		models.ReportTypeCities,
		models.ReportTypeTraffic,
		models.ReportTypeTopPages,
		models.ReportTypeDownloads,
		models.ReportTypeActiveUsers,
		models.ReportTypeSources,
		models.ReportTypeDomains,
		models.ReportTypeAgencies,
	}

	for _, expected := range expectedTypes {
		found := false
		for _, actual := range reportTypes {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetAllReportTypes() missing expected type: %s", expected)
		}
	}
}

func TestReportArgs_JSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     models.ReportArgs
		expected string
	}{
		{
			name: "complete report args",
			args: models.ReportArgs{
				ReportName: "devices",
				Limit:      100,
				Page:       2,
				After:      "2024-01-01",
				Before:     "2024-01-31",
			},
			expected: `{"report_name":"devices","limit":100,"page":2,"after":"2024-01-01","before":"2024-01-31"}`,
		},
		{
			name: "minimal report args",
			args: models.ReportArgs{
				ReportName: "browsers",
			},
			expected: `{"report_name":"browsers"}`,
		},
		{
			name: "report args with only limit",
			args: models.ReportArgs{
				ReportName: "traffic",
				Limit:      50,
			},
			expected: `{"report_name":"traffic","limit":50}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			jsonData, err := json.Marshal(tt.args)
			if err != nil {
				t.Errorf("Failed to marshal ReportArgs: %v", err)
			}

			if string(jsonData) != tt.expected {
				t.Errorf("JSON marshaling = %s, want %s", string(jsonData), tt.expected)
			}

			// Test unmarshaling
			var unmarshaled models.ReportArgs
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal ReportArgs: %v", err)
			}

			if unmarshaled != tt.args {
				t.Errorf("Unmarshaled ReportArgs = %+v, want %+v", unmarshaled, tt.args)
			}
		})
	}
}

func TestReports_JSON(t *testing.T) {
	t.Parallel()

	report := models.Reports{
		ID:           1,
		ReportName:   "devices",
		ReportAgency: "GSA",
		Date:         "2024-01-01",
		Device:       "desktop",
		Visits:       1000,
		Users:        800,
		Pageviews:    1500,
		BounceRate:   0.35,
	}

	jsonData, err := json.Marshal(report)
	if err != nil {
		t.Errorf("Failed to marshal Reports: %v", err)
	}

	var unmarshaled models.Reports
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal Reports: %v", err)
	}

	// Test key fields
	if unmarshaled.ID != report.ID {
		t.Errorf("Unmarshaled ID = %d, want %d", unmarshaled.ID, report.ID)
	}
	if unmarshaled.ReportName != report.ReportName {
		t.Errorf("Unmarshaled ReportName = %s, want %s", unmarshaled.ReportName, report.ReportName)
	}
	if unmarshaled.Device != report.Device {
		t.Errorf("Unmarshaled Device = %s, want %s", unmarshaled.Device, report.Device)
	}
	if unmarshaled.Visits != report.Visits {
		t.Errorf("Unmarshaled Visits = %d, want %d", unmarshaled.Visits, report.Visits)
	}
}

func TestReports_ParseDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		date      string
		expectErr bool
		expected  time.Time
	}{
		{
			name:      "valid date",
			date:      "2024-01-01",
			expectErr: false,
			expected:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "valid date with different month",
			date:      "2023-12-31",
			expectErr: false,
			expected:  time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "invalid date format",
			date:      "2024/01/01",
			expectErr: true,
		},
		{
			name:      "invalid date format with time",
			date:      "2024-01-01 15:04:05",
			expectErr: true,
		},
		{
			name:      "empty date",
			date:      "",
			expectErr: true,
		},
		{
			name:      "invalid date",
			date:      "invalid-date",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			report := models.Reports{Date: tt.date}
			result, err := report.ParseDate()

			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseDate() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ParseDate() unexpected error: %v", err)
				}
				if !result.Equal(tt.expected) {
					t.Errorf("ParseDate() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestReportParams_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    models.ReportParams
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid params",
			params: models.ReportParams{
				Limit:  1000,
				Page:   1,
				After:  "2024-01-01",
				Before: "2024-01-31",
			},
			expectErr: false,
		},
		{
			name: "valid params with zero values",
			params: models.ReportParams{
				Limit: 0,
				Page:  0,
			},
			expectErr: false,
		},
		{
			name: "valid minimum limit",
			params: models.ReportParams{
				Limit: 1,
			},
			expectErr: false,
		},
		{
			name: "valid maximum limit",
			params: models.ReportParams{
				Limit: 10000,
			},
			expectErr: false,
		},
		{
			name: "invalid limit too high",
			params: models.ReportParams{
				Limit: 10001,
			},
			expectErr: true,
			errMsg:    "limit must be between 1 and 10000",
		},
		{
			name: "invalid limit negative",
			params: models.ReportParams{
				Limit: -1,
			},
			expectErr: true,
			errMsg:    "limit must be between 1 and 10000",
		},
		{
			name: "valid page number",
			params: models.ReportParams{
				Page: 1,
			},
			expectErr: false,
		},
		{
			name: "valid high page number",
			params: models.ReportParams{
				Page: 100,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.params.Validate()

			validateTestResult(t, err, tt.expectErr, tt.errMsg)
		})
	}
}

func TestError_JSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		error    models.Error
		expected string
	}{
		{
			name: "complete error",
			error: models.Error{
				Message:        "Invalid report name",
				RequiredFields: "report_name",
				Example:        `{"report_name": "devices"}`,
			},
			expected: `{"message":"Invalid report name","required_fields":"report_name",` +
				`"example":"{\"report_name\": \"devices\"}"}`,
		},
		{
			name: "minimal error",
			error: models.Error{
				Message: "Something went wrong",
			},
			expected: `{"message":"Something went wrong"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			jsonData, err := json.Marshal(tt.error)
			if err != nil {
				t.Errorf("Failed to marshal Error: %v", err)
			}

			if string(jsonData) != tt.expected {
				t.Errorf("JSON marshaling = %s, want %s", string(jsonData), tt.expected)
			}

			// Test unmarshaling
			var unmarshaled models.Error
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal Error: %v", err)
			}

			if unmarshaled != tt.error {
				t.Errorf("Unmarshaled Error = %+v, want %+v", unmarshaled, tt.error)
			}
		})
	}
}

func TestReportRequest_JSON(t *testing.T) {
	t.Parallel()

	request := models.ReportRequest{
		ReportName: "devices",
		AgencyName: "GSA",
		Domain:     "example.gov",
		Parameters: models.ReportParams{
			Limit:  500,
			Page:   2,
			After:  "2024-01-01",
			Before: "2024-01-31",
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Errorf("Failed to marshal ReportRequest: %v", err)
	}

	var unmarshaled models.ReportRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal ReportRequest: %v", err)
	}

	if unmarshaled.ReportName != request.ReportName {
		t.Errorf("Unmarshaled ReportName = %s, want %s", unmarshaled.ReportName, request.ReportName)
	}
	if unmarshaled.AgencyName != request.AgencyName {
		t.Errorf("Unmarshaled AgencyName = %s, want %s", unmarshaled.AgencyName, request.AgencyName)
	}
	if unmarshaled.Domain != request.Domain {
		t.Errorf("Unmarshaled Domain = %s, want %s", unmarshaled.Domain, request.Domain)
	}
	if unmarshaled.Parameters.Limit != request.Parameters.Limit {
		t.Errorf("Unmarshaled Parameters.Limit = %d, want %d", unmarshaled.Parameters.Limit, request.Parameters.Limit)
	}
}

func TestReportResponse_JSON(t *testing.T) {
	t.Parallel()

	response := models.ReportResponse{
		Data: []models.Reports{
			{
				ID:           1,
				ReportName:   "devices",
				ReportAgency: "GSA",
				Date:         "2024-01-01",
				Device:       "desktop",
				Visits:       1000,
			},
			{
				ID:           2,
				ReportName:   "devices",
				ReportAgency: "GSA",
				Date:         "2024-01-01",
				Device:       "mobile",
				Visits:       800,
			},
		},
		Error: nil,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal ReportResponse: %v", err)
	}

	var unmarshaled models.ReportResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal ReportResponse: %v", err)
	}

	if len(unmarshaled.Data) != len(response.Data) {
		t.Errorf("Unmarshaled Data length = %d, want %d", len(unmarshaled.Data), len(response.Data))
	}

	if unmarshaled.Data[0].Device != response.Data[0].Device {
		t.Errorf("Unmarshaled Data[0].Device = %s, want %s", unmarshaled.Data[0].Device, response.Data[0].Device)
	}
}

func TestReportResponse_WithError(t *testing.T) {
	t.Parallel()

	response := models.ReportResponse{
		Data: nil,
		Error: &models.Error{
			Message:        "Invalid report name",
			RequiredFields: "report_name",
		},
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal ReportResponse with error: %v", err)
	}

	var unmarshaled models.ReportResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal ReportResponse with error: %v", err)
	}

	if unmarshaled.Error == nil {
		t.Error("Expected error field to be present")
	} else if unmarshaled.Error.Message != response.Error.Message {
		t.Errorf("Unmarshaled Error.Message = %s, want %s", unmarshaled.Error.Message, response.Error.Message)
	}
}

// validateTestResult is a helper function to reduce complexity in validation tests.
func validateTestResult(t *testing.T, err error, expectErr bool, errMsg string) {
	t.Helper()

	if expectErr {
		if err == nil {
			t.Errorf("Validate() expected error but got none")
			return
		}

		if errMsg != "" && err.Error() != errMsg {
			// For specific error message checking, check if it contains the expected text
			if len(err.Error()) < len(errMsg) || err.Error()[:len(errMsg)] != errMsg {
				t.Errorf("Validate() error = %v, want error containing %v", err, errMsg)
			}
		}
		return
	}

	if err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}
}
