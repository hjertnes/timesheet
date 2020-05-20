package runner

import (
	"os"
	"testing"
	"time"

	"github.com/hjertnes/timesheet/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ReadMock struct {
	mock.Mock
}

func (r ReadMock) Execute(f *os.File) string {
	args := r.Called(f)
	return args.String(0)
}

func TestNew(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}
	r := ReadMock{}
	_ = New(&d, r)

}

func TestSettingToInt(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}
	d.Configuration["a"] = "1"
	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}
	v := r.settingToInt("a")
	assert.Equal(t, 1, v)

}
func TestSettingToIntFails(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}
	d.Configuration["a"] = "1"
	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}
	assert.Panics(t, func() { r.settingToInt("b") })
}

func TestGetSettings(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}
	d.Configuration["workday"] = "1"
	d.Configuration["break"] = "2"

	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}

	workday, breaktime := r.getSettings()
	assert.Equal(t, workday, 1)
	assert.Equal(t, breaktime, 2)
}

func TestSettingsList(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}
	d.Configuration["workday"] = "1"
	d.Configuration["break"] = "2"

	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}

	r.SettingsList()
}

func TestSettingsSet(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}

	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}

	r.SettingsSet("a", "b")
	assert.Equal(t, d.Configuration["a"], "b")
}

func TestAdd(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}

	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}

	r.Add(time.Now(), time.Now(), false)
	v := d.Items["2020"]
	assert.NotNil(t, v)
}

func TestOff(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}

	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}

	r.Off(time.Now())
	v := d.Items["2020"]
	assert.NotNil(t, v)
}

func TestList(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}

	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}

	r.Add(time.Now(), time.Now(), false)
	r.Add(time.Now(), time.Now(), true)
	r.List()
	r.Off(time.Now())
	r.List()
}

func TestSummary(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}

	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}

	d.Configuration["workday"] = "1"
	d.Configuration["break"] = "2"

	r.Add(time.Now(), time.Now(), false)
	r.SummaryYear()
	r.Add(time.Now(), time.Now(), true)
	r.SummaryYear()
	r.Off(time.Now())
	r.SummaryYear()

}

func TestSummaryDay(t *testing.T) {
	d := models.Document{
		Configuration: make(map[string]string),
		Items:         make(map[string]map[string]models.DayItem),
	}

	rm := ReadMock{}
	r := &runner{
		reader:   rm,
		document: &d,
	}

	d.Configuration["workday"] = "1"
	d.Configuration["break"] = "2"

	r.Add(time.Now(), time.Now(), false)
	r.SummaryDay()
	r.Add(time.Now(), time.Now(), true)
	r.SummaryDay()
	r.Off(time.Now())
	r.SummaryDay()

}
