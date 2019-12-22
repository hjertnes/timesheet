package settings

import (
	"github.com/hjertnes/timesheet/models"
	"github.com/hjertnes/timesheet/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ISettingsRepository interface {
	GetOne(key string) (*models.Setting, error)
	GetAll() ([]models.Setting, error)
	AddOrUpdate(key string, value string) error
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

func (s *SettingsRepository) AddOrUpdate(key string, value string) error {
	var exist, err = s.Exist(key)
	utils.ErrorHandler(err)
	if exist {
		var setting, err = s.GetOne(key)
		s.db.Model(&setting).UpdateColumn("value", value)
		return err
	} else {
		var setting = &models.Setting{Key: key, Value: value}
		var result = s.db.Create(&setting)
		return result.Error
	}

}

func (s *SettingsRepository) DeleteAll() error {
	var result = s.db.Delete(&models.Setting{})
	return result.Error
}
