package api

// ValidatePayloadRequestObject reflects /v1/validate request schema bounds natively mapped from OpenAPI
type ValidatePayloadRequestObject struct {
	SchemaName string `json:"schema_name"`
	Version    string `json:"version,omitempty"`
	Payload    string `json:"payload"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidatePayloadResponseObject struct {
	Status      string            `json:"status"` // PASS, FAIL
	CoercedJson *string           `json:"coerced_json,omitempty"`
	RawJson     string            `json:"raw_json"`
	Errors      []ValidationError `json:"errors,omitempty"`
}

// TelemetryStats reflects /v1/stats metrics data payload dynamically feeding React frontend
type TelemetryStats struct {
	KPI        map[string]interface{}   `json:"kpi"`
	Heatmap    []map[string]interface{} `json:"heatmap"`
	Timeseries []map[string]interface{} `json:"timeseries"`
}
