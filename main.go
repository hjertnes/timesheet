package main

import (
	"github.com/hjertnes/timesheet/cmd"
	"github.com/hjertnes/timesheet/database"
	EventRepository "github.com/hjertnes/timesheet/repositories/event"
	SettingsRepository "github.com/hjertnes/timesheet/repositories/settings"
	"github.com/hjertnes/timesheet/runner"
)

func main() {
	var database = database.OpenTest()
	var settingsRepository SettingsRepository.ISettingsRepository = SettingsRepository.NewSettingsRepository(database.Db)
	var eventRepository EventRepository.IEventRepository = EventRepository.NewEventRepository(database.Db)
	var r = runner.NewRunner(eventRepository, settingsRepository)
	cmd.Run(r)
	defer database.Db.Close()
}
