package handler

import "testing"

func TestParseOptionalHourDefaults(t *testing.T) {
	start, err := parseOptionalHour("", 0)
	if err != nil || start != 0 {
		t.Fatalf("start hour = %d, err = %v; want 0, nil", start, err)
	}

	end, err := parseOptionalHour("", 24)
	if err != nil || end != 24 {
		t.Fatalf("end hour = %d, err = %v; want 24, nil", end, err)
	}
}

func TestParseOptionalHourRange(t *testing.T) {
	for _, value := range []string{"0", "14", "24"} {
		if _, err := parseOptionalHour(value, 0); err != nil {
			t.Errorf("parseOptionalHour(%q) returned error: %v", value, err)
		}
	}

	for _, value := range []string{"-1", "25", "not-an-hour"} {
		if _, err := parseOptionalHour(value, 0); err == nil {
			t.Errorf("parseOptionalHour(%q) returned nil error", value)
		}
	}
}
