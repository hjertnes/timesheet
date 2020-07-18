package read

import (
	"git.sr.ht/~hjertnes/timesheet/utils"

	"os"
	"testing"
)

func Test(t *testing.T) {
	New()
	var r = &read{}

	var f, _ = os.Create("./test")

	var _, err = f.WriteString("Test\n")

	utils.ErrorHandler(err)

	_ = f.Close()

	f, _ = os.Open("./test")

	_ = r.Execute(f)

	_ = os.Remove("./test")
}
