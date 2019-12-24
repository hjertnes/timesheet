// Package utils various util methods
package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hjertnes/timesheet/models"
)

func timeFromString(datestr string) (time.Time, error) {
	return time.Parse(time.RFC3339, datestr)
}

// TimeFromDateString turns a datestring into time.Time at 00:00:00
func TimeFromDateString(datestr string) (time.Time, error) {
	return timeFromString(fmt.Sprintf("%sT00:00:00Z", datestr))
}

// TimeFromDateStringAndTimeString turns a date string and a time string into time.Time
func TimeFromDateStringAndTimeString(datestr string, timestr string) (time.Time, error) {
	return timeFromString(fmt.Sprintf("%sT%s:00Z", datestr, timestr))
}

// TimeFromString turns a RFC3339 into time.Time
func TimeFromString(datestr string) (time.Time, error) {
	return time.Parse(time.RFC3339, fmt.Sprintf("%sZ", datestr))
}

// IntFromString turns a string into a int
func IntFromString(number string) (int, error) {
	return strconv.Atoi(number)
}

// ErrorHandler panics any error
func ErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

// IntOfMinutesToString turns a int into a string like 0h 30m
func IntOfMinutesToString(minutes int) string {
	var m = minutes

	var h int = 0

	for {
		if m < 60 {
			break
		}

		h++

		m -= 60
	}

	return fmt.Sprintf("%sh %sm", strconv.Itoa(h), strconv.Itoa(m))
}

// BuildListOf builds map of unique based on time.Time format
func BuildListOf(format string, items []models.Event) map[string]string {
	var years = make(map[string]string)

	for _, e := range items {
		var year = e.Start.Format(format)
		years[year] = year
	}

	return years
}

// FilterEventsFrom return all items that match format with string
func FilterEventsFrom(format string, items []models.Event, year string) []models.Event {
	events := make([]models.Event, 0)

	for _, i := range items {
		if i.Start.Format(format) == year {
			events = append(events, i)
		}
	}

	return events
}

// CountDaysNotExcluded count elements not excluded
func CountDaysNotExcluded(items []models.Event) int {
	var numberOfDays = 0

	var days = BuildListOf("2006-01-02", items)

	for _, day := range days {
		var dayEvents = FilterEventsFrom("2006-01-02", items, day)

		if !IsDayExcluded(dayEvents) {
			numberOfDays++
		}
	}

	return numberOfDays
}

// CalculateTotal count the difference between start and end in a list and add it all together
func CalculateTotal(items []models.Event) int {
	var total int = 0

	for _, i := range items {
		var diff = i.End.Sub(i.Start)
		total += int(diff.Minutes())
	}

	return total
}

// IsDayExcluded check if any in a list of events are excluded
func IsDayExcluded(events []models.Event) bool {
	var excluded = false

	for _, i := range events {
		if i.Excluded {
			excluded = true
		}
	}

	return excluded
}
