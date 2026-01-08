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
