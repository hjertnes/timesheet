package settings

import (
	"github.com/hjertnes/timesheet/models"
	"github.com/hjertnes/timesheet/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ISettingsRepository interface {
	Exist(key string) (bool, error)
	GetOne(key string) (*models.Setting, error)
	GetAll() ([]models.Setting, error)
	Add(key string, value string) error
	Update(key string, value string) error
	Delete(key string) error
	DeleteAll() error
}

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (s *SettingsRepository) Exist(key string) (bool, error) {
	var count int
	var result = s.db.Model(&models.Setting{}).Where("key = ?", key).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count == 1, result.Error
}

func (s *SettingsRepository) GetOne(key string) (*models.Setting, error) {
	var setting models.Setting
	var result = s.db.Where("key = ?", key).First(&setting)
	return &setting, result.Error
}

func (s *SettingsRepository) GetAll() ([]models.Setting, error) {
	var settings []models.Setting
	var result = s.db.Find(&settings)
	return settings, result.Error
}

func (s *SettingsRepository) Add(key string, value string) error {
	var setting = &models.Setting{Key: key, Value: value}
	var result = s.db.Create(&setting)
	return result.Error
}

func (s *SettingsRepository) Update(key string, value string) error {
	var setting, err = s.GetOne(key)
	s.db.Model(&setting).UpdateColumn("value", value)
	return err
}

func (s *SettingsRepository) Delete(key string) error {
	var exist, err = s.Exist(key)
	utils.ErrorHandler(err)
	if exist {
		var setting, err = s.GetOne(key)
		if err != nil {
			return err
		}
		var result = s.db.Delete(&setting)
		return result.Error
	}
	return nil
}

func (s *SettingsRepository) DeleteAll() error {
	var result = s.db.Delete(&models.Setting{})
	return result.Error
}
