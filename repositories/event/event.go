package event

import (
	"github.com/hjertnes/timesheet/models"
	"github.com/jinzhu/gorm"
	"time"
)

type IEventRepository interface {
	GetAll() ([]models.Event, error)
	Add(start time.Time, end time.Time, excluded bool, off bool) error
	Delete(id int) error
	DeleteAll() error
}

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (e *EventRepository) getOne(id int) (*models.Event, error) {
	var event *models.Event
	var result = e.db.Where("id = ?", id).First(&event)
	return event, result.Error
}

func (e *EventRepository) GetAll() ([]models.Event, error) {
	var events []models.Event
	var result = e.db.Find(&events).Order("start asc")
	return events, result.Error
}

func (e *EventRepository) Add(start time.Time, end time.Time, excluded bool, off bool) error {
	var event = &models.Event{Start: start, End: end, Excluded: excluded, Off: off}
	var result = e.db.Create(&event)
	return result.Error
}

func (e *EventRepository) Delete(id int) error {
	var event, err = e.getOne(id)
	if err != nil {
		return err
	}
	var result = e.db.Delete(&event)
	return result.Error
}

func (e *EventRepository) DeleteAll() error {
	var result = e.db.Delete(&models.Event{})
	return result.Error
}
