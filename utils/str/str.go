package str

import "time"

func GetTodayTimeRange() (time.Time, time.Time) {
	now := time.Now()

	// Get start of day (00:00:00)
	startOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, // hour
		0, // minute
		0, // second
		0, // nanosecond
		now.Location(),
	)

	// Get end of day (23:59:59)
	endOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		23,        // hour
		59,        // minute
		59,        // second
		999999999, // nanosecond (optional, for complete precision)
		now.Location(),
	)

	return startOfDay, endOfDay
}
