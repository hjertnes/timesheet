package read

import (
	"bufio"
	"os"

	"github.com/hjertnes/timesheet/utils"
)

type IRead interface {
	Execute(f *os.File) string
}

type Read struct{}

//os.Stdin
func (r Read) Execute(f *os.File) string {
	var reader = bufio.NewReader(f)
	var line, err = reader.ReadString('\n')
	utils.ErrorHandler(err)
	return line
}
