package utils

import (
	"errors"
	"os"
	"testing"
	"time"

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

func TestTimeFromDateStringAnTimeString2(t *testing.T) {

	var d, err = TimeFromDateStringAndTimeString2("2010-01-01", "08:00:00")

	assert.Nil(t, err)
	assert.Equal(t, d.Year(), 2010)
	assert.Equal(t, int(d.Month()), 1)
	assert.Equal(t, d.Day(), 1)
	assert.Equal(t, d.Hour(), 8)
	assert.Equal(t, d.Minute(), 0)
	assert.Equal(t, d.Second(), 0)

	_, err = TimeFromDateStringAndTimeString2("abc", "bc")
	assert.NotNil(t, err)
	_, err = TimeFromDateStringAndTimeString2("2010-01-01", "bc")
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

func TestExist(t *testing.T) {
	assert.False(t, exist("/tmp/file"))
	f, _ := OpenOrCreate("/tmp/file")
	_ = f.Close()
	assert.True(t, exist("/tmp/file"))
	os.Remove("/tmp/file")
}
