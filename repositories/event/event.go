// Package event repository
package event

import (
	"time"

	"github.com/hjertnes/timesheet/models"
	"github.com/jinzhu/gorm"
)

// Repository for events
type Repository interface {
	GetAll() ([]models.Event, error)
	Add(start time.Time, end time.Time, excluded bool, off bool) error
	Delete(id int) error
	DeleteAll() error
}

type repository struct {
	db *gorm.DB
}

// New constructor
func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (e *repository) getOne(id int) (*models.Event, error) {
	var event *models.Event = &models.Event{}

	var result = e.db.Where("id = ?", id).First(&event)

	return event, result.Error
}

// GetAll returns all events
func (e *repository) GetAll() ([]models.Event, error) {
	var events []models.Event

	var result = e.db.Find(&events).Order("start asc")

	return events, result.Error
}

// Add creates new event
func (e *repository) Add(start time.Time, end time.Time, excluded bool, off bool) error {
	var event = &models.Event{Start: start, End: end, Excluded: excluded, Off: off}

	var result = e.db.Create(&event)

	return result.Error
}

// Delete removes evenet
func (e *repository) Delete(id int) error {
	var event, err = e.getOne(id)

	if err != nil {
		return err
	}

	var result = e.db.Delete(event)

	return result.Error
}

// DeleteAll deletes all evenst
func (e *repository) DeleteAll() error {
	var result = e.db.Delete(&models.Event{})
	return result.Error
}
