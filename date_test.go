package main

import (
	"testing"
	"time"
)

func TestWorkWeek(t *testing.T) {
	cases := []struct {
		input        time.Time
		expectedDays []string
	}{
		{
			// Thrs April 1st
			input: time.Date(2021, 04, 01, 0, 0, 0, 0, time.UTC),
			expectedDays: []string{
				"2021-03-29",
				"2021-03-30",
				"2021-03-31",
				"2021-04-01",
				"2021-04-02",
			},
		},
		{
			// Sat April 3rd
			input: time.Date(2021, 04, 03, 0, 0, 0, 0, time.UTC),
			expectedDays: []string{
				"2021-03-29",
				"2021-03-30",
				"2021-03-31",
				"2021-04-01",
				"2021-04-02",
			},
		},
		{
			// Sun February 28th
			input: time.Date(2021, 02, 28, 0, 0, 0, 0, time.UTC),
			expectedDays: []string{
				"2021-03-01",
				"2021-03-02",
				"2021-03-03",
				"2021-03-04",
				"2021-03-05",
			},
		},
	}

	fmt := func(d time.Time) string { return d.Format("2006-01-02") }

	for _, tc := range cases {
		t.Run(fmt(tc.input), func(t *testing.T) {
			r := workWeek(tc.input)

			t.Log("dateRange.days:")
			for _, day := range r.days {
				t.Logf("  %s", fmt(day))
			}

			if got, want := len(tc.expectedDays), len(r.days); got != want {
				t.Errorf("got %d days; want %d", got, want)
			}
			for i, want := range tc.expectedDays {
				if got := fmt(r.days[i]); got != want {
					t.Errorf("days[%d] = %q; expected %q", i, got, want)
				}
			}
		})
	}

}
