package control

import (
	"testing"
	"time"
)

func TestNextDailyAt(t *testing.T) {
	// 2023-01-01 12:00:00
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		hour     int
		minute   int
		second   int
		expected time.Time
	}{
		{
			name:     "Same day later time",
			hour:     14,
			minute:   30,
			second:   0,
			expected: time.Date(2023, 1, 1, 14, 30, 0, 0, time.UTC),
		},
		{
			name:     "Next day earlier time",
			hour:     10,
			minute:   0,
			second:   0,
			expected: time.Date(2023, 1, 2, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "Next day same time (exact logic check)",
			hour:     12,
			minute:   0,
			second:   0,
			// Since !next.After(now) includes equality, same time counts as "past/present" so moves to tomorrow
			expected: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NextDailyAt(now, tc.hour, tc.minute, tc.second)
			if !got.Equal(tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestNextWeeklyAt(t *testing.T) {
	// 2023-01-01 is a Sunday
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		weekday  time.Weekday
		hour     int
		minute   int
		second   int
		expected time.Time
	}{
		{
			name:     "Next Monday (same week)",
			weekday:  time.Monday,
			hour:     9,
			minute:   0,
			second:   0,
			expected: time.Date(2023, 1, 2, 9, 0, 0, 0, time.UTC),
		},
		{
			name:     "Next Saturday (same week)",
			weekday:  time.Saturday,
			hour:     23,
			minute:   59,
			second:   59,
			expected: time.Date(2023, 1, 7, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "Next Sunday (next week, since today is Sunday 12:00)",
			weekday:  time.Sunday,
			hour:     10,
			minute:   0,
			second:   0,
			expected: time.Date(2023, 1, 8, 10, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NextWeeklyAt(now, tc.weekday, tc.hour, tc.minute, tc.second)
			if !got.Equal(tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestNextMonthlyAt(t *testing.T) {
	// 2023-01-15 12:00:00
	now := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		day      int
		hour     int
		minute   int
		second   int
		expected time.Time
	}{
		{
			name:     "Same month later day",
			day:      20,
			hour:     10,
			minute:   0,
			second:   0,
			expected: time.Date(2023, 1, 20, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "Next month earlier day",
			day:      5,
			hour:     9,
			minute:   0,
			second:   0,
			expected: time.Date(2023, 2, 5, 9, 0, 0, 0, time.UTC),
		},
		{
			name:     "Next month (today is past 15th 12:00)",
			day:      15,
			hour:     11,
			minute:   0,
			second:   0,
			expected: time.Date(2023, 2, 15, 11, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NextMonthlyAt(now, tc.day, tc.hour, tc.minute, tc.second)
			if !got.Equal(tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
