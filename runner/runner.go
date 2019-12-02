package runner

import (
	"bufio"
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
	EventRepository "github.com/hjertnes/timesheet/repositories/event"
	SettingsRepository "github.com/hjertnes/timesheet/repositories/settings"
	"github.com/hjertnes/timesheet/utils"
	"github.com/olekukonko/tablewriter"
)

type IRunner interface {
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

type Runner struct {
	eventRepository    EventRepository.IEventRepository
	settingsRepository SettingsRepository.ISettingsRepository
}

func NewRunner(e EventRepository.IEventRepository, s SettingsRepository.ISettingsRepository) IRunner {
	return &Runner{
		eventRepository:    e,
		settingsRepository: s,
	}
}

func (r *Runner) settingToInt(name string) int {
	var err error
	var setting *models.Setting
	var result int
	setting, err = r.settingsRepository.GetOne(name)
	utils.ErrorHandler(err)
	result, err = strconv.Atoi(setting.Value)
	utils.ErrorHandler(err)
	return result
}

func (r *Runner) getSettings() (int, int) {
	return r.settingToInt("workday"), r.settingToInt("break")
}

func (r *Runner) SettingsList() {
	var items, err = r.settingsRepository.GetAll()
	if err != nil {
		utils.ErrorHandler(err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})
	for _, v := range items {
		table.Append([]string{v.Key, v.Value})
	}
	table.Render()
}

func (r *Runner) SettingsSet(key string, value string) {
	var err = r.settingsRepository.AddOrUpdate(key, value)
	utils.ErrorHandler(err)
}

func (r *Runner) List() {
	var items, err = r.eventRepository.GetAll()
	utils.ErrorHandler(err)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Start", "End", "Off", "Excluded"})
	for _, v := range items {
		table.Append([]string{strconv.FormatUint(v.ID, 2), v.Start.String(), v.End.String(), strconv.FormatBool(v.Off), strconv.FormatBool(v.Excluded)})
	}
	table.Render()

}

func (r *Runner) Add(start time.Time, end time.Time, excluded bool) {
	var err = r.eventRepository.Add(start, end, excluded, false)
	utils.ErrorHandler(err)
}

func (r *Runner) Off(date time.Time) {
	var err = r.eventRepository.Add(date, date, false, true)
	utils.ErrorHandler(err)
}

func (r *Runner) Delete(id int) {
	var err = r.eventRepository.Delete(id)
	utils.ErrorHandler(err)
}

func read() string {
	var reader = bufio.NewReader(os.Stdin)
	var line, err = reader.ReadString('\n')
	utils.ErrorHandler(err)
	return line
}

func (r *Runner) Setup() {
	fmt.Println("Setup")
	fmt.Println("This will replace your current settings but not your data")
	fmt.Print("Work day in minutes: ")
	var workDayMinutes = read()
	fmt.Print("Break in minutes: ")
	var breakInMinutes = read()
	r.settingsRepository.AddOrUpdate("workday", strings.Trim(workDayMinutes, "\n"))
	r.settingsRepository.AddOrUpdate("break", strings.Trim(breakInMinutes, "\n"))
}

func (r *Runner) Backup(filename string) {
	var allSettings, _ = r.settingsRepository.GetAll()
	var allEvents, _ = r.eventRepository.GetAll()
	var settings = make([]backupmodels.Setting, 0)
	var events = make([]backupmodels.Event, 0)
	for _, setting := range allSettings {
		settings = append(settings, backupmodels.Setting{Key: setting.Key, Value: setting.Value})
	}
	for _, event := range allEvents {
		events = append(events, backupmodels.Event{Start: event.Start, End: event.End, Excluded: event.Excluded, Off: event.Off})
	}
	var document = backupmodels.Document{
		Settings: settings,
		Events:   events,
	}
	file, _ := json.MarshalIndent(document, "", " ")

	var err = ioutil.WriteFile(filename, file, 0644)
	utils.ErrorHandler(err)
}

func (r *Runner) Restore(filename string) {
	var obj backupmodels.Document
	data, err := ioutil.ReadFile(filename)
	utils.ErrorHandler(err)
	err = json.Unmarshal(data, &obj)
	utils.ErrorHandler(err)
	r.eventRepository.DeleteAll()
	r.settingsRepository.DeleteAll()
	for _, e := range obj.Events {
		err = r.eventRepository.Add(e.Start, e.End, e.Excluded, e.Off)
		utils.ErrorHandler(err)
	}
	for _, e := range obj.Settings {
		err = r.settingsRepository.AddOrUpdate(e.Key, e.Value)
		utils.ErrorHandler(err)
	}
}

func (r *Runner) SummaryYear() {
	var workday, breaktime = r.getSettings()
	var items, err = r.eventRepository.GetAll()
	utils.ErrorHandler(err)
	var years = utils.BuildListOf("2006", items)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Year", "Expected", "Total", "Difference"})
	var data = make([][]string, 0)
	for _, e := range years {
		events := utils.FilterEventsFrom("2006", items, e)
		var numberOfDays = utils.CountDaysNotExcluded(events)
		var expected int = int(numberOfDays) * workday
		var total int = utils.CalculateTotal(events)
		total = total - (numberOfDays * breaktime)
		var diff int = int(total) - expected
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

func (r *Runner) SummaryDay() {
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
		if excluded == false {
			total = total - breaktime
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
