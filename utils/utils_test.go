package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/hjertnes/timesheet/models"
	"github.com/stretchr/testify/assert"
)

func TestIntOfMinutesToString(t *testing.T) {
	assert.Equal(t, IntOfMinutesToString(0), "0h 0m")
	assert.Equal(t, IntOfMinutesToString(60), "1h 0m")
	assert.Equal(t, IntOfMinutesToString(30), "0h 30m")
	assert.Equal(t, IntOfMinutesToString(125), "2h 5m")
}

func TestErrorHandler(t *testing.T) {
	assert.Panics(t, func() { ErrorHandler(errors.New("Test")) })
	assert.NotPanics(t, func() { ErrorHandler(nil) })
}

func TestTimeFromDateString(t *testing.T) {
	var d, err = TimeFromDateString("2010-01-01")

	assert.Nil(t, err)
	assert.Equal(t, d.Year(), 2010)
	assert.Equal(t, int(d.Month()), 1)
	assert.Equal(t, d.Day(), 1)

	_, err = TimeFromDateString("Hello World")

	assert.NotNil(t, err)
}

func TestTimeFromDateStringAnTimeString(t *testing.T) {
	var d, err = TimeFromDateStringAndTimeString("2010-01-01", "08:00")

	assert.Nil(t, err)
	assert.Equal(t, d.Year(), 2010)
	assert.Equal(t, int(d.Month()), 1)
	assert.Equal(t, d.Day(), 1)
	assert.Equal(t, d.Hour(), 8)
	assert.Equal(t, d.Minute(), 0)
	assert.Equal(t, d.Second(), 0)

	_, err = TimeFromDateStringAndTimeString("abc", "bc")
	assert.NotNil(t, err)
	_, err = TimeFromDateStringAndTimeString("2010-01-01", "bc")
	assert.NotNil(t, err)
}

func TestIntFromString(t *testing.T) {
	var d, err = IntFromString("1")

	assert.Nil(t, err)
	assert.Equal(t, d, 1)

	_, err = IntFromString("A")

	assert.NotNil(t, err)
}

func TestTimeFromString(t *testing.T) {
	var d, err = TimeFromString("2010-01-01T08:00:00")

	assert.Nil(t, err)
	assert.NotNil(t, d)
}

func unwrap(e time.Time, err error) time.Time {
	return e
}

func TestBuildListOf(t *testing.T) {
	var events = []models.Event{
		{Start: unwrap(TimeFromDateString("2010-01-01"))},
		{Start: unwrap(TimeFromDateString("2010-01-01"))},
		{Start: unwrap(TimeFromDateString("2010-01-02"))},
		{Start: unwrap(TimeFromDateString("2010-01-02"))},
		{Start: unwrap(TimeFromDateString("2010-01-03"))},
	}

	var days = BuildListOf("2006-01-02", events)

	assert.Len(t, days, 3)
	days = BuildListOf("2006", events)
	assert.Len(t, days, 1)
}

func TestFilterEventsFrom(t *testing.T) {
	var events = []models.Event{
		{Start: unwrap(TimeFromDateString("2010-01-01"))},
		{Start: unwrap(TimeFromDateString("2010-01-01"))},
		{Start: unwrap(TimeFromDateString("2010-01-02"))},
		{Start: unwrap(TimeFromDateString("2010-01-02"))},
		{Start: unwrap(TimeFromDateString("2010-01-03"))},
	}

	var days = FilterEventsFrom("2006", events, "2010")

	assert.Len(t, days, 5)
	days = FilterEventsFrom("2006-01-02", events, "2010-01-03")
	assert.Len(t, days, 1)
}

func TestCountDaysNotExcluded(t *testing.T) {
	var events = []models.Event{
		{Start: unwrap(TimeFromDateString("2010-01-01")), Excluded: true},
		{Start: unwrap(TimeFromDateString("2010-01-02")), Excluded: false},
		{Start: unwrap(TimeFromDateString("2010-01-03")), Excluded: false},
		{Start: unwrap(TimeFromDateString("2010-01-04")), Excluded: false},
		{Start: unwrap(TimeFromDateString("2010-01-05")), Excluded: true},
	}

	var days = CountDaysNotExcluded(events)

	assert.Equal(t, days, 3)
}

func TestCalculateTotal(t *testing.T) {
	var events = []models.Event{
		{
			Start: unwrap(TimeFromDateStringAndTimeString("2010-01-01", "08:00")),
			End:   unwrap(TimeFromDateStringAndTimeString("2010-01-01", "10:00")),
		},
		{
			Start: unwrap(TimeFromDateStringAndTimeString("2010-01-02", "08:00")),
			End:   unwrap(TimeFromDateStringAndTimeString("2010-01-02", "10:00")),
		},
		{
			Start: unwrap(TimeFromDateStringAndTimeString("2010-01-03", "08:00")),
			End:   unwrap(TimeFromDateStringAndTimeString("2010-01-03", "10:00")),
		},
	}

	var days = CalculateTotal(events)

	assert.Equal(t, days, 360)
}

/*


 */
