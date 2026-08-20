package service

import (
	"testing"
	"time"

	"power-price-monitor/model"
)

func TestSelectConsecutiveValuesPrefersLongestQualifyingWindow(t *testing.T) {
	values := []model.Value{
		{Value: 30, DateTime: "2026-08-20T14:00:00.000+00:00"},
		{Value: 40, DateTime: "2026-08-20T15:00:00.000+00:00"},
		{Value: 90, DateTime: "2026-08-20T16:00:00.000+00:00"},
		{Value: 100, DateTime: "2026-08-20T17:00:00.000+00:00"},
		{Value: 20, DateTime: "2026-08-20T18:00:00.000+00:00"},
	}

	selected := selectConsecutiveValues(values, 3, 60, time.UTC)

	if len(selected) != 3 {
		t.Fatalf("selected %d values, want 3", len(selected))
	}
	if selected[0].DateTime != "2026-08-20T14:00:00.000+00:00" {
		t.Errorf("selected window starts at %s, want 14:00", selected[0].DateTime)
	}
}

func TestSelectConsecutiveValuesUsesWindowAverageForThreshold(t *testing.T) {
	values := []model.Value{
		{Value: 100, DateTime: "2026-08-20T14:00:00.000+00:00"},
		{Value: 20, DateTime: "2026-08-20T15:00:00.000+00:00"},
	}

	selected := selectConsecutiveValues(values, 2, 60, time.UTC)

	if len(selected) != 2 {
		t.Fatalf("selected %d values, want 2", len(selected))
	}
}

func TestSelectConsecutiveValuesRespectsGapsAndMaxHours(t *testing.T) {
	values := []model.Value{
		{Value: 10, DateTime: "2026-08-20T14:00:00.000+00:00"},
		{Value: 10, DateTime: "2026-08-20T16:00:00.000+00:00"},
		{Value: 10, DateTime: "2026-08-20T17:00:00.000+00:00"},
	}

	selected := selectConsecutiveValues(values, 2, 20, time.UTC)

	if len(selected) != 2 || selected[0].DateTime != "2026-08-20T16:00:00.000+00:00" {
		t.Errorf("selected window = %#v, want 16:00-17:00", selected)
	}
}

func TestSelectConsecutiveValuesUsesHalfOpenHourRangeInput(t *testing.T) {
	values := []model.Value{
		{Value: 10, DateTime: "2026-08-20T13:00:00.000+00:00"},
		{Value: 10, DateTime: "2026-08-20T14:00:00.000+00:00"},
		{Value: 10, DateTime: "2026-08-20T20:00:00.000+00:00"},
		{Value: 10, DateTime: "2026-08-20T21:00:00.000+00:00"},
	}

	filtered := filterValuesByHourRange(values, 14, 21, time.UTC)

	if len(filtered) != 2 || filtered[0].DateTime != "2026-08-20T14:00:00.000+00:00" || filtered[1].DateTime != "2026-08-20T20:00:00.000+00:00" {
		t.Fatalf("filtered values = %#v, want 14:00 and 20:00", filtered)
	}
}
