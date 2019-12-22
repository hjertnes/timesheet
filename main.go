package main

import (
	"github.com/hjertnes/timesheet/cmd"
	"github.com/hjertnes/timesheet/database"
	"github.com/hjertnes/timesheet/read"
	EventRepository "github.com/hjertnes/timesheet/repositories/event"
	SettingsRepository "github.com/hjertnes/timesheet/repositories/settings"
	"github.com/hjertnes/timesheet/runner"
)

func main() {
	var database = database.Open()
	var settingsRepository SettingsRepository.ISettingsRepository = SettingsRepository.NewSettingsRepository(database.Db)
	var eventRepository EventRepository.IEventRepository = EventRepository.NewEventRepository(database.Db)
	var read = read.Read{}
	var r = runner.NewRunner(eventRepository, settingsRepository, read)
	var rf = cmd.NewRunFunc(r)
	cmd.Run(rf, r)
	defer database.Db.Close()
}
