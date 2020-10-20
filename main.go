package main

import (
	"fmt"
	"os"

	"git.sr.ht/~hjertnes/timesheet/cmd"
	"git.sr.ht/~hjertnes/timesheet/models"
	"git.sr.ht/~hjertnes/timesheet/read"
	"git.sr.ht/~hjertnes/timesheet/runner"
	"git.sr.ht/~hjertnes/timesheet/utils"
)

func main() {
	home := os.Getenv("HOME")
	filename := fmt.Sprintf("%s/txt/timesheet.yaml", home)

	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		repo := models.New(filename)
		d := &models.Document{
			Configuration: map[string]string{
				"workday": "450",
				"break": "30",
			},
			Items: map[string]map[string]models.DayItem{},
		}

		err = repo.Save(d)
		utils.ErrorHandler(err)
	}

	repo := models.New(filename)
	d, err := repo.Load()
	utils.ErrorHandler(err)

	var rr = read.New()

	var r = runner.New(d, rr)

	var rf = cmd.New(r)

	cmd.Run(rf, r)

	err = repo.Save(d)
	utils.ErrorHandler(err)
}
