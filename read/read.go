// Package read from *os.File
package read

import (
	"bufio"
	"os"

	"git.sr.ht/~hjertnes/timesheet/utils"
)

// Read Interface for reading files usually stdin
type Read interface {
	Execute(f *os.File) string
}

type read struct{}

// New constructor
func New() Read {
	return &read{}
}

//Execute read os.File
func (r read) Execute(f *os.File) string {
	var reader = bufio.NewReader(f)

	var line, err = reader.ReadString('\n')

	utils.ErrorHandler(err)

	return line
}
