package parser

import (
	"testing"
	"time"
)

func TestParseDateStandardFormats(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
		checkFn  func(time.Time) bool
		desc     string
	}{
		{"2026-02-15", false, func(t time.Time) bool { return t.Year() == 2026 && t.Month() == 2 && t.Day() == 15 }, "ISO date"},
		{"2026-02-15 14:30", false, func(t time.Time) bool { return t.Hour() == 14 && t.Minute() == 30 }, "ISO datetime"},
		{"2026-02-15 2:30PM", false, func(t time.Time) bool { return t.Hour() == 14 && t.Minute() == 30 }, "ISO date with 12h time"},
		{"Jan 15, 2026", false, func(t time.Time) bool { return t.Month() == 1 && t.Day() == 15 }, "Month day year"},
		{"", true, nil, "empty string"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result, err := ParseDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}
			if tt.checkFn != nil && !tt.checkFn(result) {
				t.Errorf("check failed for input %q, got %v", tt.input, result)
			}
		})
	}
}

func TestParseDateNaturalLanguage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		input   string
		wantErr bool
		checkFn func(time.Time) bool
		desc    string
	}{
		{"today", false, func(t time.Time) bool { return t.Day() == now.Day() }, "today"},
		{"tomorrow", false, func(t time.Time) bool { return t.Day() == now.AddDate(0, 0, 1).Day() }, "tomorrow"},
		{"next week", false, func(t time.Time) bool { return t.After(now) }, "next week"},
		{"next month", false, func(t time.Time) bool { return t.After(now) }, "next month"},
		{"in 2 days", false, func(t time.Time) bool { return t.Day() == now.AddDate(0, 0, 2).Day() }, "in 2 days"},
		{"in 3 hours", false, func(t time.Time) bool { return t.After(now) }, "in 3 hours"},
		{"in 30 minutes", false, func(t time.Time) bool { return t.After(now) }, "in 30 minutes"},
		{"in 1 week", false, func(t time.Time) bool { return t.After(now) }, "in 1 week"},
		{"eod", false, func(t time.Time) bool { return t.Hour() == 17 }, "end of day"},
		{"end of day", false, func(t time.Time) bool { return t.Hour() == 17 }, "end of day long"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result, err := ParseDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}
			if tt.checkFn != nil && !tt.checkFn(result) {
				t.Errorf("check failed for input %q, got %v", tt.input, result)
			}
		})
	}
}

func TestParseDateNextWeekday(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		desc    string
	}{
		{"next monday", false, "next monday"},
		{"next friday", false, "next friday"},
		{"next sunday", false, "next sunday"},
		{"next monday at 2pm", false, "next monday at 2pm"},
		{"next friday at 3:30pm", false, "next friday at 3:30pm"},
		{"next invalid", true, "invalid day"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result, err := ParseDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}
			if result.Before(time.Now()) {
				t.Errorf("next weekday should be in the future, got %v", result)
			}
		})
	}
}

func TestParseDateWithTime(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		hour    int
		min     int
		desc    string
	}{
		{"today at 5pm", false, 17, 0, "today at 5pm"},
		{"tomorrow at 3:30pm", false, 15, 30, "tomorrow at 3:30pm"},
		{"today at 9am", false, 9, 0, "today at 9am"},
		{"today at 17:00", false, 17, 0, "today at 17:00"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result, err := ParseDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}
			if result.Hour() != tt.hour || result.Minute() != tt.min {
				t.Errorf("for input %q: expected %02d:%02d, got %02d:%02d", tt.input, tt.hour, tt.min, result.Hour(), result.Minute())
			}
		})
	}
}

func TestParseDateNaturalFormats(t *testing.T) {
	now := time.Now()

	tests := []struct {
		input   string
		wantErr bool
		checkFn func(time.Time) bool
		desc    string
	}{
		// Day-first without year
		{"21 Mar", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Year() == now.Year() && t.Hour() == 9
		}, "21 Mar (no time)"},
		{"21 Mar 2pm", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Hour() == 14 && t.Minute() == 0
		}, "21 Mar 2pm"},
		{"21 Mar 2:30pm", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Hour() == 14 && t.Minute() == 30
		}, "21 Mar 2:30pm"},
		{"21 march 2pm", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Hour() == 14
		}, "21 march 2pm (full month)"},

		// Month-first without year
		{"Mar 21", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Year() == now.Year() && t.Hour() == 9
		}, "Mar 21 (no time)"},
		{"mar 21 2pm", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Hour() == 14
		}, "mar 21 2pm"},
		{"march 21 2pm", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Hour() == 14
		}, "march 21 2pm"},

		// With "at"
		{"march 21 at 2pm", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Hour() == 14
		}, "march 21 at 2pm"},
		{"21 mar at 3:30pm", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Hour() == 15 && t.Minute() == 30
		}, "21 mar at 3:30pm"},

		// With year
		{"21 March 2026", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Year() == 2026
		}, "21 March 2026"},
		{"mar 21 2026", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Year() == 2026
		}, "mar 21 2026"},
		{"mar 21 2026 2pm", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Year() == 2026 && t.Hour() == 14
		}, "mar 21 2026 2pm"},

		// ISO + bare 12h time
		{"2026-03-21 2pm", false, func(t time.Time) bool {
			return t.Year() == 2026 && t.Month() == 3 && t.Day() == 21 && t.Hour() == 14
		}, "ISO date + 2pm"},
		{"2026-03-21 2:30pm", false, func(t time.Time) bool {
			return t.Year() == 2026 && t.Month() == 3 && t.Day() == 21 && t.Hour() == 14 && t.Minute() == 30
		}, "ISO date + 2:30pm"},

		// Single-digit day
		{"5 mar", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 5 && t.Hour() == 9
		}, "5 mar (single digit day)"},
		{"mar 5", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 5 && t.Hour() == 9
		}, "mar 5 (single digit month-first)"},

		// Day boundaries
		{"1 jan", false, func(t time.Time) bool {
			return t.Month() == 1 && t.Day() == 1
		}, "1 jan (min day)"},
		{"31 jan", false, func(t time.Time) bool {
			return t.Month() == 1 && t.Day() == 31
		}, "31 jan (max day)"},

		// 24h time with natural date
		{"21 mar 14:30", false, func(t time.Time) bool {
			return t.Month() == 3 && t.Day() == 21 && t.Hour() == 14 && t.Minute() == 30
		}, "21 mar 14:30 (24h time)"},

		// Invalid month-day combos (Go time.Date overflow)
		{"feb 30", true, nil, "feb 30 (invalid)"},
		{"feb 31", true, nil, "feb 31 (invalid)"},
		{"apr 31", true, nil, "apr 31 (invalid)"},
		{"feb 29 2026", true, nil, "feb 29 non-leap year"},

		// Edge cases
		{"not a date", true, nil, "invalid input"},
		{"32 Mar", true, nil, "invalid day"},
		{"0 mar", true, nil, "day zero"},
		{"march", true, nil, "month name alone"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result, err := ParseDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got %v", tt.input, result)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}
			if tt.checkFn != nil && !tt.checkFn(result) {
				t.Errorf("check failed for input %q, got %v", tt.input, result)
			}
		})
	}
}

func TestParseDateTimeOnly(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		hour    int
		min     int
		desc    string
	}{
		{"5pm", false, 17, 0, "5pm"},
		{"9am", false, 9, 0, "9am"},
		{"3:30pm", false, 15, 30, "3:30pm"},
		{"17:00", false, 17, 0, "17:00"},
		{"12am", false, 0, 0, "midnight"},
		{"12pm", false, 12, 0, "noon"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result, err := ParseDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}
			if result.Hour() != tt.hour || result.Minute() != tt.min {
				t.Errorf("for input %q: expected %02d:%02d, got %02d:%02d", tt.input, tt.hour, tt.min, result.Hour(), result.Minute())
			}
		})
	}
}
