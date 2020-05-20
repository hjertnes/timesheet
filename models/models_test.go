package models

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_ = New("/tmp/filename")
	os.Remove("/tmp/filename")

}

func TestLoadSave(t *testing.T) {
	d := Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]DayItem),
	}
	d.Configuration["test"] = "1"
	d.Add(time.Now(), time.Now(), false, false)
	d.Add(time.Now(), time.Now(), false, true)
	assert.NotNil(t, d)
	r := New("/tmp/filename")
	r.Save(&d)
	r.Load()

	os.Remove("/tmp/filename")
}
