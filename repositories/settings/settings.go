// Package settings repository for settings
package settings

import (
	"github.com/hjertnes/timesheet/models"
	"github.com/hjertnes/timesheet/utils"

	"github.com/jinzhu/gorm"
	//Sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Repository Interface for the repo
type Repository interface {
	GetOne(key string) (*models.Setting, error)
	GetAll() ([]models.Setting, error)
	AddOrUpdate(key string, value string) error
	DeleteAll() error
	Exist(key string) (bool, error)
}

// SettingsRepository struct for the repo
type repository struct {
	db *gorm.DB
}

// New constructor
func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Exist check if a key exist
func (s *repository) Exist(key string) (bool, error) {
	var count int

	var result = s.db.Model(&models.Setting{}).Where("key = ?", key).Count(&count)

	return count == 1, result.Error
}

// GetOne get one setting matching item
func (s *repository) GetOne(key string) (*models.Setting, error) {
	var setting models.Setting

	var result = s.db.Where("key = ?", key).First(&setting)

	return &setting, result.Error
}

// GetAll get all settings
func (s *repository) GetAll() ([]models.Setting, error) {
	var settings []models.Setting

	var result = s.db.Find(&settings)

	return settings, result.Error
}

// AddOrUpdate creates setting if it doesnt exist or updates if it does
func (s *repository) AddOrUpdate(key string, value string) error {
	var exist, err = s.Exist(key)

	utils.ErrorHandler(err)

	if exist {
		var setting, err = s.GetOne(key)

		s.db.Model(&setting).UpdateColumn("value", value)

		return err
	}

	var setting = &models.Setting{Key: key, Value: value}

	var result = s.db.Create(&setting)

	return result.Error
}

//DeleteAll deletes all settings
func (s *repository) DeleteAll() error {
	var result = s.db.Delete(&models.Setting{})
	return result.Error
}
