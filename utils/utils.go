package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
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
