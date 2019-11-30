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
	Import(filename string)
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
	var exist, err = r.settingsRepository.Exist(key)
	utils.ErrorHandler(err)
	if exist {
		var err = r.settingsRepository.Update(key, value)
		utils.ErrorHandler(err)
	} else {
		var err = r.settingsRepository.Add(key, value)
		utils.ErrorHandler(err)
	}
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
	r.settingsRepository.Delete("workday")
	r.settingsRepository.Delete("break")
	r.settingsRepository.Add("workday", strings.Trim(workDayMinutes, "\n"))
	r.settingsRepository.Add("break", strings.Trim(breakInMinutes, "\n"))
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

	_ = ioutil.WriteFile(filename, file, 0644)

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
		r.eventRepository.Add(e.Start, e.End, e.Excluded, e.Off)
	}
	for _, e := range obj.Settings {
		r.settingsRepository.Add(e.Key, e.Value)
	}
}
func (r *Runner) Import(filename string) {
	type Event struct {
		Start    string
		End      string
		Excluded bool
		Off      bool
	}
	type Document struct {
		Events []Event
	}

	var obj Document
	data, err := ioutil.ReadFile(filename)
	utils.ErrorHandler(err)
	err = json.Unmarshal(data, &obj)
	utils.ErrorHandler(err)
	for _, e := range obj.Events {
		var start time.Time
		var end time.Time
		var err error
		start, err = utils.TimeFromString(e.Start)
		utils.ErrorHandler(err)
		end, err = utils.TimeFromString(e.End)
		utils.ErrorHandler(err)
		r.eventRepository.Add(start, end, e.Excluded, e.Off)
	}
}
func (r *Runner) SummaryYear() {
	var items, err = r.eventRepository.GetAll()
	utils.ErrorHandler(err)

	var years map[int]int = make(map[int]int)

	for _, e := range items {
		var year = e.Start.Year()
		years[year] = year
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Year", "Expected", "Total", "Difference"})
	var data = make([][]string, 0)
	for _, e := range years {
		events := make([]models.Event, 0)
		days := make(map[string]bool)
		for _, i := range items {
			if i.Start.Year() == e {
				events = append(events, i)
			}
		}
		for _, i := range events {
			var day = i.Start.Format("2006-01-02")
			if days[day] {
				if days[day] == false {
					days[day] = i.Excluded
				}
			} else {
				days[day] = i.Excluded
			}
		}

		var numberOfDays = 0
		for _, i := range days {
			if i == false {
				numberOfDays++
			}
		}

		var expected int = int(numberOfDays) * 425
		var total int = 0

		for _, i := range events {
			var diff = i.End.Sub(i.Start)
			total += int(diff.Minutes())
		}
		total = total - (numberOfDays * 30)
		var diff int = int(total) - expected
		data = append(data, []string{strconv.Itoa(e), utils.IntOfMinutesToString(expected), utils.IntOfMinutesToString(total), utils.IntOfMinutesToString(diff)})
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

	var items, err = r.eventRepository.GetAll()
	utils.ErrorHandler(err)

	var days = make(map[string]string)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Hours"})
	var data = make([][]string, 0)
	for _, i := range items {
		days[i.Start.Format("2006-01-02")] = i.Start.Format("2006-01-02")
	}
	for _, day := range days {
		var excluded = false
		var events []models.Event = make([]models.Event, 0)
		for _, i := range items {
			if i.Start.Format("2006-01-02") == day {
				events = append(events, i)
				if i.Excluded == true {
					excluded = true
				}
			}
		}
		var total int = 0
		for _, event := range events {
			var diff = event.End.Sub(event.Start)
			total = total + int(diff.Minutes())
		}
		if excluded == false {
			total = total - 30
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
