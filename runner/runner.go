// Package runner runs and prints output
package runner

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hjertnes/timesheet/models"
	"github.com/hjertnes/timesheet/read"
	"github.com/hjertnes/timesheet/utils"
	"github.com/olekukonko/tablewriter"
)

// Runner methods
type Runner interface {
	SettingsList()
	SettingsSet(key string, value string)
	List()
	Add(start time.Time, end time.Time, excluded bool)
	Off(date time.Time)
	Setup()
	SummaryYear()
	SummaryDay()
}

type runner struct {
	document *models.Document
	reader   read.Read
}

// New constructor
func New(d *models.Document, r read.Read) Runner {
	return &runner{
		reader:   r,
		document: d,
	}
}

func (r *runner) settingToInt(name string) int {
	var err error

	var result int

	setting, ok := r.document.Configuration[name]
	if !ok {
		err = errors.New("Key not found")
	}
	utils.ErrorHandler(err)
	result, err = strconv.Atoi(setting)

	utils.ErrorHandler(err)

	return result
}

func (r *runner) getSettings() (int, int) {
	return r.settingToInt("workday"), r.settingToInt("break")
}

// SettingsList prints a table of settings
func (r *runner) SettingsList() {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Key", "Value"})

	for key, v := range r.document.Configuration {
		table.Append([]string{key, v})
	}

	table.Render()
}

// SettingsSet adds or updates a setting
func (r *runner) SettingsSet(key string, value string) {
	r.document.Configuration[key] = value
}

// List lists events
func (r *runner) List() {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Start", "End", "Off", "Excluded"})

	for _, yearValue := range r.document.Items {
		for day, dayItem := range yearValue {
			for _, item := range dayItem.Events {
				table.Append([]string{
					fmt.Sprint(day, " ", item.Start),
					fmt.Sprint(day, " ", item.End),
					strconv.FormatBool(false),
					strconv.FormatBool(dayItem.Excluded),
				})
			}
			if len(dayItem.Events) == 0 {
				table.Append([]string{
					fmt.Sprint(day),
					fmt.Sprint(day),
					strconv.FormatBool(true),
					strconv.FormatBool(dayItem.Excluded),
				})
			}
		}
	}

	table.Render()
}

//Add add event
func (r *runner) Add(start time.Time, end time.Time, excluded bool) {
	r.document.Add(start, end, excluded, false)
}

// Off add a day as "off"
func (r *runner) Off(date time.Time) {
	r.document.Add(date, date, false, true)
}

// Setup settings
func (r *runner) Setup() {
	fmt.Println("Setup")
	fmt.Println("This will replace your current settings but not your data")
	fmt.Print("Work day in minutes: ")

	var workDayMinutes = r.reader.Execute(os.Stdin)

	fmt.Print("Break in minutes: ")

	var breakInMinutes = r.reader.Execute(os.Stdin)

	r.document.Configuration["workday"] = strings.Trim(workDayMinutes, "\n")
	r.document.Configuration["break"] = strings.Trim(breakInMinutes, "\n")
}

// SummaryYear show summary per year with difference between expected hours and actual hours
func (r *runner) SummaryYear() {
	var workday, breaktime = r.getSettings()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Year", "Expected", "Total", "Difference"})

	var data = make([][]string, 0)

	for year, yearValue := range r.document.Items {
		var numberOfDays int = 0
		var expected int = 0
		var total int = 0

		for day, dayItem := range yearValue {
			if !dayItem.Excluded {
				numberOfDays++
			}
			for _, item := range dayItem.Events {
				s, err := utils.TimeFromDateStringAndTimeString2(day, item.Start)
				utils.ErrorHandler(err)
				e, err := utils.TimeFromDateStringAndTimeString2(day, item.End)
				utils.ErrorHandler(err)

				var diff = e.Sub(s)
				total += int(diff.Minutes())
			}

		}

		expected = numberOfDays * workday

		total -= (numberOfDays * breaktime)

		var diff int = total - expected

		data = append(data, []string{
			year,
			utils.IntOfMinutesToString(expected),
			utils.IntOfMinutesToString(total),
			utils.IntOfMinutesToString(diff),
		})
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	for _, e := range data {
		table.Append(e)
	}

	table.Render()
}

// SummaryDay shows list of dates and sum of hours on that day
func (r *runner) SummaryDay() {
	var _, breaktime = r.getSettings()

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Date", "Hours"})

	var data = make([][]string, 0)

	for _, yearValue := range r.document.Items {
		for day, dayItem := range yearValue {
			var total int = 0
			for _, item := range dayItem.Events {
				s, err := utils.TimeFromDateStringAndTimeString2(day, item.Start)
				utils.ErrorHandler(err)
				e, err := utils.TimeFromDateStringAndTimeString2(day, item.End)
				utils.ErrorHandler(err)

				var diff = e.Sub(s)
				total += int(diff.Minutes())
			}
			if !dayItem.Excluded {
				total -= breaktime
			}
			if total > 0 {
				data = append(data, []string{day, utils.IntOfMinutesToString(total)})
			}

		}
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	for _, e := range data {
		table.Append(e)
	}

	table.Render()
}
