package api

import (
	"fmt"
	"strings"
	"time"
)

// Date forms accepted by the metering bindings.
//
// The wire format is YYYY-MM-DD, but ingext addresses days as a dayIndex --
// YYYYMMDD, UTC, the same form utils.ParseDayIndex and the lake's @dayIndex
// field use. Operators think in dayIndex because that is what they see in
// index names and dump paths, so both are accepted everywhere a date is taken
// and normalised to the wire form here rather than at each call site.
const (
	dateLayout     = "2006-01-02"
	dayIndexLayout = "20060102"
	monthLayout    = "2006-01"
	monthTightLay  = "200601"
)

// NormalizeDate accepts YYYY-MM-DD or a YYYYMMDD dayIndex and returns the wire
// form.
//
// Parsing rather than pattern-matching, so 2026-02-30 is rejected here instead
// of becoming a range the server answers with an empty ledger.
func NormalizeDate(s string) (string, error) {
	s = strings.TrimSpace(s)
	switch len(s) {
	case len(dateLayout):
		t, err := time.Parse(dateLayout, s)
		if err != nil {
			return "", fmt.Errorf("invalid date %q: %w", s, err)
		}
		return t.Format(dateLayout), nil
	case len(dayIndexLayout):
		t, err := time.Parse(dayIndexLayout, s)
		if err != nil {
			return "", fmt.Errorf("invalid dayIndex %q: %w", s, err)
		}
		return t.Format(dateLayout), nil
	default:
		return "", fmt.Errorf("date %q must be YYYY-MM-DD or a YYYYMMDD dayIndex", s)
	}
}

// DayIndex renders a wire date as a YYYYMMDD dayIndex.
func DayIndex(date string) string {
	t, err := time.Parse(dateLayout, strings.TrimSpace(date))
	if err != nil {
		return date
	}
	return t.Format(dayIndexLayout)
}

// MonthRange expands YYYY-MM (or YYYYMM) into the first and last day of that
// month, in wire form.
//
// A month still in progress is clamped to YESTERDAY, not to today: today has
// not finished, is never billable, and asking for it would make every
// current-month report end on a day that can only be `open` or `in_progress`.
// The clamp is reported so a caller can say what it actually covered.
func MonthRange(month string) (from, to string, partial bool, err error) {
	month = strings.TrimSpace(month)
	var t time.Time
	switch len(month) {
	case len(monthLayout):
		t, err = time.Parse(monthLayout, month)
	case len(monthTightLay):
		t, err = time.Parse(monthTightLay, month)
	default:
		return "", "", false, fmt.Errorf("month %q must be YYYY-MM or YYYYMM", month)
	}
	if err != nil {
		return "", "", false, fmt.Errorf("invalid month %q: %w", month, err)
	}

	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1)

	yesterday := time.Now().UTC().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	if last.After(yesterday) {
		if yesterday.Before(first) {
			// The whole month is in the future, or today is its first day, so
			// there is not one finished day in it yet.
			return "", "", false, fmt.Errorf("month %s has no completed day yet", month)
		}
		last = yesterday
		partial = true
	}
	return first.Format(dateLayout), last.Format(dateLayout), partial, nil
}
