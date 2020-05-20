package main

import (
	"fmt"
	"os"

	"github.com/hjertnes/timesheet/cmd"
	"github.com/hjertnes/timesheet/models"
	"github.com/hjertnes/timesheet/read"
	"github.com/hjertnes/timesheet/runner"
	"github.com/hjertnes/timesheet/utils"
)

func main() {
	home := os.Getenv("HOME")
	filename := fmt.Sprintf("%s/txt/timesheet.yaml", home)

	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		repo := models.New(filename)
		d := &models.Document{
			Configuration: map[string]string{},
		}
		d.Configuration["workday"] = "450"
		d.Configuration["break"] = "30"
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
