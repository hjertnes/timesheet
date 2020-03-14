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
	var home = os.Getenv("HOME")

	repo := models.New(fmt.Sprintf("%s/txt/timesheet.yaml", home))
	d, err := repo.Load()
	utils.ErrorHandler(err)

	var rr = read.New()

	var r = runner.New(d, rr)

	var rf = cmd.New(r)

	cmd.Run(rf, r)

	err = repo.Save(d)
	utils.ErrorHandler(err)
}
