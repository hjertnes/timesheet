package main

import (
	"github.com/hjertnes/timesheet/cmd"
	"github.com/hjertnes/timesheet/runner"
)

func main() {
	var r = runner.NewRunner()
	cmd.Run(r)
}
