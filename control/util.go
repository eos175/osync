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

// NextWeeklyAt calculates the next occurrence of a specific weekday and time.
func NextWeeklyAt(
	now time.Time,
	weekday time.Weekday,
	hour, minute, second int,
) time.Time {
	loc := now.Location()

	// Calculate days until the target weekday
	daysUntil := int(weekday - now.Weekday())
	if daysUntil < 0 {
		daysUntil += 7
	}

	next := time.Date(
		now.Year(), now.Month(), now.Day()+daysUntil,
		hour, minute, second, 0,
		loc,
	)

	// If the calculated time is in the past (e.g., today is Monday 10am, and we want Monday 9am),
	// add a full week.
	if !next.After(now) {
		next = next.AddDate(0, 0, 7)
	}

	return next
}

// NextMonthlyAt calculates the next occurrence of a specific day of the month and time.
// Note: If the day exceeds the days in the target month (e.g. 31st of February),
// standard Go time normalization applies (e.g., March 3rd).
func NextMonthlyAt(
	now time.Time,
	day, hour, minute, second int,
) time.Time {
	loc := now.Location()

	next := time.Date(
		now.Year(), now.Month(), day,
		hour, minute, second, 0,
		loc,
	)

	if !next.After(now) {
		// Move to next month
		next = next.AddDate(0, 1, 0)
	}

	return next
}
