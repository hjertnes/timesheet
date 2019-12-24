package main

import (
	"github.com/hjertnes/timesheet/cmd"
	"github.com/hjertnes/timesheet/database"
	"github.com/hjertnes/timesheet/read"
	EventRepository "github.com/hjertnes/timesheet/repositories/event"
	SettingsRepository "github.com/hjertnes/timesheet/repositories/settings"
	"github.com/hjertnes/timesheet/runner"
	"github.com/hjertnes/timesheet/utils"
)

func main() {
	var database = database.Open()

	var settingsRepository SettingsRepository.Repository = SettingsRepository.New(database.Db)

	var eventRepository EventRepository.Repository = EventRepository.New(database.Db)

	var rr = read.New()

	var r = runner.New(eventRepository, settingsRepository, rr)

	var rf = cmd.New(r)

	cmd.Run(rf, r)

	var err error

	defer func() {
		err = database.Db.Close()
	}()

	utils.ErrorHandler(err)
}
