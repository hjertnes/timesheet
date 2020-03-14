// Package utils various util methods
package utils

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
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

// TimeFromDateStringAndTimeString2 turns a date string and a time string into time.Time
func TimeFromDateStringAndTimeString2(datestr string, timestr string) (time.Time, error) {
	return timeFromString(fmt.Sprintf("%sT%sZ", datestr, timestr))
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

func exist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// OpenOrCreate opens the file, and creates it first if it doesn't exist
func OpenOrCreate(filename string) (io.ReadWriteCloser, error) {
	e := exist(filename)
	if !e {
		f, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		err = f.Close()
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(filename, os.O_RDWR, 0600)
}
