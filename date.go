package main

import (
	"time"
)

const (
	day  = 24 * time.Hour
	week = 7 * day
)

// relWorkWeek returns a Mon-Fri range relative
// to weekOffset number of weeks from now. The offset
// can be negative, 0, and positive.
func relWorkWeek(weekOffset int) dateRange {
	w := time.Now().Local()
	if weekOffset != 0 {
		w = w.Add(time.Duration(weekOffset) * week)
	}
	return workWeek(w)
}

// workWeek returns a Mon-Fri range for the given date. If date is
// Sunday, range will be the following days. If date is Saturday, range
// will be the previous days.
func workWeek(date time.Time) dateRange {
	y, m, d := date.Date()
	wd := date.Weekday()

	monday, friday := d+int(time.Monday-wd), d+int(time.Friday-wd)

	days := make([]time.Time, 0, friday-monday)
	for d := monday; d <= friday; d++ {
		// time.Date allows negative 'day' values, and also day values
		// that extend past the current month.
		days = append(days, time.Date(y, m, d, 0, 0, 0, 0, date.Location()))
	}

	startYear, startWeek := days[0].ISOWeek()

	return dateRange{
		startYear: startYear,
		startWeek: startWeek,
		date:      date,
		days:      days,
	}
}

// dateRange is a list of days. The startYear
// and startWeek values are the ISO 8601 values
// for the start of the range.
type dateRange struct {
	startYear, startWeek int
	date                 time.Time
	days                 []time.Time
}
