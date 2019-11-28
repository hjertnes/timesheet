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

func IntFromString(number string) (int, error) {
	return strconv.Atoi(number)
}

func ErrorHandler(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
