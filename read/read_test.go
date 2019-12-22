package read

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	var r = &Read{}
	var f, _ = os.Create("./test")
	f.WriteString("Test\n")
	f.Close()
	f, _ = os.Open("./test")
	r.Execute(f)
	os.Remove("./test")
}
