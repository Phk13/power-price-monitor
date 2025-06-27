package model

import "time"

// IndicatorResponse represents the top-level structure of the API response
type IndicatorResponse struct {
	Indicator Indicator `json:"indicator"`
}

// Indicator contains the details of the power price indicator
type Indicator struct {
	Values []Value `json:"values"`
}

// Value represents a single data point for power price
type Value struct {
	Value    float64 `json:"value"`
	DateTime string  `json:"datetime"`
	GeoID    int     `json:"geo_id"`
}

// CachedPowerData holds cached power price data with metadata
type CachedPowerData struct {
	Values     []Value   `json:"values"`
	CachedAt   time.Time `json:"cached_at"`
	ValidUntil time.Time `json:"valid_until"`
}

// OptimalHoursResponse represents the API response for optimal hours
type OptimalHoursResponse struct {
	OptimalHours       []string `json:"optimal_hours"`
	NextStart          string   `json:"next_start,omitempty"`
	TotalHoursSelected int      `json:"total_hours_selected"`
	CurrentPrice       float64  `json:"current_price"`
	ThresholdUsed      float64  `json:"threshold_used"`
	MaxHoursUsed       int      `json:"max_hours_used"`
}
