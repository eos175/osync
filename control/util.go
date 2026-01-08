package control

import "time"

// NextDailyAt calculates the next occurrence of a specific time (hour, minute, second).
// It respects the time zone (Location) of the provided 'now' time.
func NextDailyAt(
	now time.Time,
	hour, minute, second int,
) time.Time {
	loc := now.Location()

	next := time.Date(
		now.Year(), now.Month(), now.Day(),
		hour, minute, second, 0,
		loc,
	)

	if !next.After(now) {
		// recalculate "tomorrow" correctly
		y, m, d := now.AddDate(0, 0, 1).Date()
		next = time.Date(y, m, d, hour, minute, second, 0, loc)
	}

	return next
}
