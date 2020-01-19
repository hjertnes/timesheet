// Package runner runs and prints output
package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hjertnes/timesheet/backupmodels"
	"github.com/hjertnes/timesheet/models"
	"github.com/hjertnes/timesheet/read"
	EventRepository "github.com/hjertnes/timesheet/repositories/event"
	SettingsRepository "github.com/hjertnes/timesheet/repositories/settings"
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
	Delete(id int)
	Setup()
	Backup(filename string)
	Restore(filename string)
	SummaryYear()
	SummaryDay()
}

type runner struct {
	eventRepository    EventRepository.Repository
	settingsRepository SettingsRepository.Repository
	reader             read.Read
}

// New constructor
func New(e EventRepository.Repository, s SettingsRepository.Repository, r read.Read) Runner {
	return &runner{
		eventRepository:    e,
		settingsRepository: s,
		reader:             r,
	}
}

func (r *runner) settingToInt(name string) int {
	var err error

	var setting *models.Setting

	var result int

	setting, err = r.settingsRepository.GetOne(name)
	utils.ErrorHandler(err)
	result, err = strconv.Atoi(setting.Value)

	utils.ErrorHandler(err)

	return result
}

func (r *runner) getSettings() (int, int) {
	return r.settingToInt("workday"), r.settingToInt("break")
}

// SettingsList prints a table of settings
func (r *runner) SettingsList() {
	var items, err = r.settingsRepository.GetAll()

	utils.ErrorHandler(err)

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Key", "Value"})

	for _, v := range items {
		table.Append([]string{v.Key, v.Value})
	}

	table.Render()
}

// SettingsSet adds or updates a setting
func (r *runner) SettingsSet(key string, value string) {
	var err = r.settingsRepository.AddOrUpdate(key, value)

	utils.ErrorHandler(err)
}

// List lists events
func (r *runner) List() {
	var items, err = r.eventRepository.GetAll()

	utils.ErrorHandler(err)

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Id", "Start", "End", "Off", "Excluded"})

	for _, v := range items {
		table.Append([]string{
			strconv.FormatUint(v.ID, 10),
			v.Start.String(),
			v.End.String(),
			strconv.FormatBool(v.Off),
			strconv.FormatBool(v.Excluded),
		})
	}

	table.Render()
}

//Add add event
func (r *runner) Add(start time.Time, end time.Time, excluded bool) {
	var err = r.eventRepository.Add(start, end, excluded, false)

	utils.ErrorHandler(err)
}

// Off add a day as "off"
func (r *runner) Off(date time.Time) {
	var err = r.eventRepository.Add(date, date, false, true)

	utils.ErrorHandler(err)
}

// Delete event
func (r *runner) Delete(id int) {
	var err = r.eventRepository.Delete(id)

	utils.ErrorHandler(err)
}

// Setup settings
func (r *runner) Setup() {
	fmt.Println("Setup")
	fmt.Println("This will replace your current settings but not your data")
	fmt.Print("Work day in minutes: ")

	var workDayMinutes = r.reader.Execute(os.Stdin)

	fmt.Print("Break in minutes: ")

	var breakInMinutes = r.reader.Execute(os.Stdin)

	err := r.settingsRepository.AddOrUpdate("workday", strings.Trim(workDayMinutes, "\n"))
	utils.ErrorHandler(err)
	err = r.settingsRepository.AddOrUpdate("break", strings.Trim(breakInMinutes, "\n"))
	utils.ErrorHandler(err)
}

// Backup to file
func (r *runner) Backup(filename string) {
	var allSettings, _ = r.settingsRepository.GetAll()

	var allEvents, _ = r.eventRepository.GetAll()

	var settings = make([]backupmodels.Setting, 0)

	var events = make([]backupmodels.Event, 0)

	for _, setting := range allSettings {
		settings = append(
			settings,
			backupmodels.Setting{
				Key:   setting.Key,
				Value: setting.Value,
			},
		)
	}

	for _, event := range allEvents {
		events = append(
			events,
			backupmodels.Event{
				Start:    event.Start,
				End:      event.End,
				Excluded: event.Excluded,
				Off:      event.Off,
			},
		)
	}

	var document = backupmodels.Document{
		Settings: settings,
		Events:   events,
	}

	file, _ := json.MarshalIndent(document, "", " ")

	var err = ioutil.WriteFile(filename, file, 0644)

	utils.ErrorHandler(err)
}

// Restore from backup
func (r *runner) Restore(filename string) {
	var obj backupmodels.Document

	data, err := ioutil.ReadFile(filename) //nolint

	utils.ErrorHandler(err)
	err = json.Unmarshal(data, &obj)
	utils.ErrorHandler(err)
	err = r.eventRepository.DeleteAll()
	utils.ErrorHandler(err)
	err = r.settingsRepository.DeleteAll()
	utils.ErrorHandler(err)

	for _, e := range obj.Events {
		err = r.eventRepository.Add(e.Start, e.End, e.Excluded, e.Off)
		utils.ErrorHandler(err)
	}

	for _, e := range obj.Settings {
		err = r.settingsRepository.AddOrUpdate(e.Key, e.Value)
		utils.ErrorHandler(err)
	}
}

// SummaryYear show summary per year with difference between expected hours and actual hours
func (r *runner) SummaryYear() {
	var workday, breaktime = r.getSettings()

	var items, err = r.eventRepository.GetAll()

	utils.ErrorHandler(err)

	var years = utils.BuildListOf("2006", items)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Year", "Expected", "Total", "Difference"})

	var data = make([][]string, 0)

	for _, e := range years {
		events := utils.FilterEventsFrom("2006", items, e)

		var numberOfDays int = utils.CountDaysNotExcluded(events)

		var expected int = numberOfDays * workday

		var total int = utils.CalculateTotal(events)

		total -= (numberOfDays * breaktime)

		var diff int = total - expected

		data = append(data, []string{
			e,
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

	var items, err = r.eventRepository.GetAll()

	utils.ErrorHandler(err)

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Date", "Hours"})

	var data = make([][]string, 0)

	var days = utils.BuildListOf("2006-01-02", items)

	for _, day := range days {
		var events = utils.FilterEventsFrom("2006-01-02", items, day)

		var excluded = utils.IsDayExcluded(events)

		var total = utils.CalculateTotal(events)

		if !excluded {
			total -= breaktime
		}

		if total > 0 {
			data = append(data, []string{day, utils.IntOfMinutesToString(total)})
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
