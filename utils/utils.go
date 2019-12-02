package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hjertnes/timesheet/models"
)

func timeFromString(datestr string) (time.Time, error) {
	return time.Parse(time.RFC3339, datestr)
}

func TimeFromDateString(datestr string) (time.Time, error) {
	return timeFromString(fmt.Sprintf("%sT00:00:00Z", datestr))
}

func TimeFromDateStringAnTimeString(datestr string, timestr string) (time.Time, error) {
	return timeFromString(fmt.Sprintf("%sT%s:00Z", datestr, timestr))
}

func TimeFromString(datestr string) (time.Time, error) {
	return time.Parse(time.RFC3339, fmt.Sprintf("%sZ", datestr))
}

func IntFromString(number string) (int, error) {
	return strconv.Atoi(number)
}

func ErrorHandler(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func IntOfMinutesToString(minutes int) string {

	var m = minutes
	var h int = 0
	for {
		if m < 60 {
			break
		}
		h++
		m = m - 60
	}
	return fmt.Sprintf("%sh %sm", strconv.Itoa(h), strconv.Itoa(m))
}

func BuildListOf(format string, items []models.Event) map[string]string {
	var years = make(map[string]string)

	for _, e := range items {
		var year = e.Start.Format(format)
		years[year] = year
	}
	return years
}

func FilterEventsFrom(format string, items []models.Event, year string) []models.Event {
	events := make([]models.Event, 0)
	for _, i := range items {
		if i.Start.Format(format) == year {
			events = append(events, i)
		}
	}
	return events
}

func CountDaysNotExcluded(items []models.Event) int {
	var numberOfDays = 0
	var days = BuildListOf("2006-01-02", items)
	for _, day := range days {
		var dayEvents = FilterEventsFrom("2006-01-02", items, day)
		if IsDayExcluded(dayEvents) == false {
			numberOfDays++
		}
	}
	return numberOfDays
}

func CalculateTotal(items []models.Event) int {
	var total int = 0

	for _, i := range items {
		var diff = i.End.Sub(i.Start)
		total += int(diff.Minutes())
	}
	return total
}

func IsDayExcluded(events []models.Event) bool {
	var excluded = false
	for _, i := range events {
		if i.Excluded == true {
			excluded = true
		}
	}
	return excluded
}
