package cmd

import (
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/mock"
)

type RunnerMock struct {
	mock.Mock
}

func (m *RunnerMock) SettingsList() {
	m.Called()
}
func (m *RunnerMock) SettingsSet(key string, value string) {
	m.Called(key, value)
}
func (m *RunnerMock) List() {
	m.Called()
}
func (m *RunnerMock) Add(start time.Time, end time.Time, excluded bool) {
	m.Called(start, end, excluded)
}
func (m *RunnerMock) Off(date time.Time) {
	m.Called(date)
}
func (m *RunnerMock) Delete(id int) {
	m.Called(id)
}
func (m *RunnerMock) Setup() {
	m.Called()
}
func (m *RunnerMock) Backup(filename string) {
	m.Called(filename)
}
func (m *RunnerMock) Restore(filename string) {
	m.Called(filename)
}
func (m *RunnerMock) SummaryYear() {
	m.Called()
}
func (m *RunnerMock) SummaryDay() {
	m.Called()
}

func TestRun(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	Run(r, m)
}

func TestRunFuncList(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("List").Return()
	r.list(cmd, []string{})
}

func TestRunFuncOff(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("Off", mock.Anything).Return()
	r.off(cmd, []string{"2010-01-01"})
}

func TestRunFuncAdd(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("Add", mock.Anything, mock.Anything, mock.Anything).Return()
	r.add(cmd, []string{"2010-01-01", "08:00", "16:00"})
}

func TestRunFuncBackup(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("Backup", "filename").Return()
	r.backup(cmd, []string{"filename"})
}

func TestRunFuncRestore(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("Restore", "filename").Return()
	r.restore(cmd, []string{"filename"})
}

func TestRunFuncDelete(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("Delete", 1).Return()
	r.deleteOne(cmd, []string{"1"})
}

func TestRunFuncSettingsList(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("SettingsList").Return()
	r.settingsList(cmd, []string{})
}

func TestRunFuncSettingsSet(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("SettingsSet", "A", "B").Return()
	r.settingsSet(cmd, []string{"A", "B"})
}

func TestRunFuncSettingsSetup(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("Setup").Return()
	r.setup(cmd, []string{})
}

func TestRunFuncSummaryYear(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("SummaryYear").Return()
	r.summaryYear(cmd, []string{})
}

func TestRunFuncSummaryDay(t *testing.T) {
	var m = &RunnerMock{}

	var r = New(m)

	var cmd = &cobra.Command{}

	m.On("SummaryDay").Return()
	r.summaryDay(cmd, []string{})
}
