package runner

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/hjertnes/timesheet/models"
	"github.com/hjertnes/timesheet/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SettingsRepositoryMock struct {
	mock.Mock
}

func (m *SettingsRepositoryMock) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *SettingsRepositoryMock) GetAll() ([]models.Setting, error) {
	args := m.Called()
	return args.Get(0).([]models.Setting), args.Error(1)
}

func (m *SettingsRepositoryMock) GetOne(key string) (*models.Setting, error) {
	args := m.Called(key)
	return args.Get(0).(*models.Setting), args.Error(1)
}

func (m *SettingsRepositoryMock) AddOrUpdate(key string, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

type EventRepositoryMock struct {
	mock.Mock
}

func (m *EventRepositoryMock) GetAll() ([]models.Event, error) {
	args := m.Called()
	return args.Get(0).([]models.Event), args.Error(1)
}

func (m *EventRepositoryMock) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *EventRepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *EventRepositoryMock) Add(start time.Time, end time.Time, excluded bool, off bool) error {
	args := m.Called(start, end, excluded, off)
	return args.Error(0)
}

type ReadMock struct {
	mock.Mock
}

func (m *ReadMock) Execute(f *os.File) string {
	args := m.Called(f)
	return args.String(0)
}

func NewRunnerWithMocks() (*SettingsRepositoryMock, *EventRepositoryMock, *ReadMock, IRunner) {
	var s = &SettingsRepositoryMock{}
	var e = &EventRepositoryMock{}
	var rm = &ReadMock{}
	var r = NewRunner(e, s, rm)
	return s, e, rm, r
}
func TestSettingsSet(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		var s, _, _, runner = NewRunnerWithMocks()
		s.On("AddOrUpdate", "A", "B").Return(nil)
		runner.SettingsSet("A", "B")
		s.AssertExpectations(t)
		s.AssertNumberOfCalls(t, "AddOrUpdate", 1)
	})
	t.Run("Err", func(t *testing.T) {
		var s, _, _, runner = NewRunnerWithMocks()
		s.On("AddOrUpdate", "A", "B").Return(errors.New("Error"))
		assert.Panics(t, func() { runner.SettingsSet("A", "B") })
		s.AssertExpectations(t)
		s.AssertNumberOfCalls(t, "AddOrUpdate", 1)
	})
}
func TestSettingsList(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var s, _, _, runner = NewRunnerWithMocks()
		s.On("GetAll").Return([]models.Setting{}, nil)
		runner.SettingsList()
		s.AssertExpectations(t)
		s.AssertNumberOfCalls(t, "GetAll", 1)
	})
	t.Run("One", func(t *testing.T) {
		var s, _, _, runner = NewRunnerWithMocks()
		s.On("GetAll").Return([]models.Setting{models.Setting{Key: "A", Value: "B"}}, nil)
		runner.SettingsList()
		s.AssertExpectations(t)
		s.AssertNumberOfCalls(t, "GetAll", 1)
	})
	t.Run("Two", func(t *testing.T) {
		var s, _, _, runner = NewRunnerWithMocks()
		s.On("GetAll").Return([]models.Setting{
			models.Setting{Key: "A", Value: "B"},
			models.Setting{Key: "B", Value: "C"},
		}, nil)
		runner.SettingsList()
		s.AssertExpectations(t)
		s.AssertNumberOfCalls(t, "GetAll", 1)
	})
	t.Run("error", func(t *testing.T) {
		var s, _, _, runner = NewRunnerWithMocks()
		s.On("GetAll").Return([]models.Setting{}, errors.New("Error"))
		assert.Panics(t, func() { runner.SettingsList() })
		s.AssertExpectations(t)
		s.AssertNumberOfCalls(t, "GetAll", 1)
	})
}

func TestAdd(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var start, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "08:00")
		var end, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "16:00")
		var _, e, _, runner = NewRunnerWithMocks()
		e.On("Add", start, end, false, false).Return(nil)
		runner.Add(start, end, false)
		e.AssertExpectations(t)
		e.AssertNumberOfCalls(t, "Add", 1)
	})
}

func TestList(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var start, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "08:00")
		var end, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "16:00")
		var _, e, _, runner = NewRunnerWithMocks()
		e.On("GetAll").Return([]models.Event{
			models.Event{
				Start:    start,
				End:      end,
				Excluded: false,
				Off:      false,
			},
		}, nil)
		runner.List()
		e.AssertExpectations(t)
		e.AssertNumberOfCalls(t, "GetAll", 1)
	})
}

func TestOff(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var start, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "00:00")
		var _, e, _, runner = NewRunnerWithMocks()
		e.On("Add", start, start, false, true).Return(nil)
		runner.Off(start)
		e.AssertExpectations(t)
		e.AssertNumberOfCalls(t, "Add", 1)
	})
}

func TestBackup(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var start, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "08:00")
		var end, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "16:00")
		var s, e, _, runner = NewRunnerWithMocks()
		s.On("GetAll").Return([]models.Setting{
			models.Setting{
				Key:   "A",
				Value: "B",
			},
		}, nil)
		e.On("GetAll").Return([]models.Event{
			models.Event{
				Start:    start,
				End:      end,
				Excluded: false,
				Off:      false,
			},
		}, nil)
		runner.Backup("test.json")
		e.AssertExpectations(t)
		s.AssertExpectations(t)
		e.AssertNumberOfCalls(t, "GetAll", 1)
		s.AssertNumberOfCalls(t, "GetAll", 1)
	})
}

func TestRestore(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var s, e, _, runner = NewRunnerWithMocks()
		s.On("DeleteAll").Return(nil)
		e.On("DeleteAll").Return(nil)
		s.On("AddOrUpdate", mock.Anything, mock.Anything).Return(nil)
		e.On("Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		runner.Restore("test.json")
		e.AssertExpectations(t)
		s.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var _, e, _, runner = NewRunnerWithMocks()
		e.On("Delete", mock.Anything).Return(nil)
		runner.Delete(1)
		e.AssertExpectations(t)
	})
}

func TestSetup(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var s, _, r, runner = NewRunnerWithMocks()
		r.On("Execute", os.Stdin).Return("0")
		s.On("AddOrUpdate", mock.Anything, mock.Anything).Return(nil)
		runner.Setup()
		r.AssertExpectations(t)
		r.AssertNumberOfCalls(t, "Execute", 2)
		s.AssertNumberOfCalls(t, "AddOrUpdate", 2)
	})
}
func getSomeEvents() []models.Event {
	var start1, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "08:00")
	var end1, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "16:00")
	var start2, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "17:00")
	var end2, _ = utils.TimeFromDateStringAnTimeString("2010-01-01", "18:00")
	var start3, _ = utils.TimeFromDateStringAnTimeString("2010-01-02", "08:00")
	var end3, _ = utils.TimeFromDateStringAnTimeString("2010-01-02", "16:00")
	var start4, _ = utils.TimeFromDateStringAnTimeString("2011-01-02", "17:00")
	var end4, _ = utils.TimeFromDateStringAnTimeString("2011-01-02", "18:00")
	return []models.Event{
		models.Event{
			Start:    start1,
			End:      end1,
			Excluded: false,
			Off:      false,
		},
		models.Event{
			Start:    start2,
			End:      end2,
			Excluded: false,
			Off:      false,
		},
		models.Event{
			Start:    start3,
			End:      end3,
			Excluded: false,
			Off:      false,
		},
		models.Event{
			Start:    start4,
			End:      end4,
			Excluded: false,
			Off:      false,
		},
	}
}

func TestSummaryYear(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var s, e, _, runner = NewRunnerWithMocks()
		s.On("GetOne", "break").Return(&models.Setting{
			Value: "0",
		}, nil)
		s.On("GetOne", "workday").Return(&models.Setting{
			Value: "0",
		}, nil)
		e.On("GetAll", mock.Anything).Return(getSomeEvents(), nil)
		runner.SummaryYear()
		e.AssertExpectations(t)
	})
}

func TestSummaryDay(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var s, e, _, runner = NewRunnerWithMocks()
		s.On("GetOne", "break").Return(&models.Setting{
			Value: "0",
		}, nil)
		s.On("GetOne", "workday").Return(&models.Setting{
			Value: "0",
		}, nil)
		e.On("GetAll", mock.Anything).Return(getSomeEvents(), nil)
		runner.SummaryDay()
		e.AssertExpectations(t)
	})
}
