package service

// ProcessingStats is a simple struct for storing counts
type ProcessingStats struct {
	TelemetryEventCount int `json:"telemetry_count"`
	AlertEventCount     int `json:"alert_count"`
	PositionEventCount  int `json:"position_count"`
}

var (
	// Stats contains a per-process count of processing counts
	Stats ProcessingStats
)
