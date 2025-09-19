package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type SettingsRepository interface {
	Create(req *entity.Settings) error
	FindByName(name string) (*entity.Settings, error)
	Update(settings *entity.Settings) error
	Delete(settings *entity.Settings) error
	FindAll() ([]*entity.Settings, error)
	FindByID(id int) (*entity.Settings, error)
}

type settingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{
		db: db,
	}
}

// Create implements SettingsRepository.
func (s *settingsRepository) Create(req *entity.Settings) error {
	return s.db.Create(req).Error
}

// Delete implements SettingsRepository.
func (s *settingsRepository) Delete(settings *entity.Settings) error {
	return s.db.Delete(settings).Error
}

// FindAll implements SettingsRepository.
func (s *settingsRepository) FindAll() ([]*entity.Settings, error) {
	var settings []*entity.Settings
	if err := s.db.Find(&settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

// FindByID implements SettingsRepository.
func (s *settingsRepository) FindByID(id int) (*entity.Settings, error) {
	var settings entity.Settings
	if err := s.db.First(&settings, id).Error; err != nil {
		return nil, err
	}

	return &settings, nil
}

// FindByName implements SettingsRepository.
func (s *settingsRepository) FindByName(name string) (*entity.Settings, error) {
	var settings entity.Settings
	if err := s.db.Where("name = ?", name).First(&settings).Error; err != nil {
		return nil, err
	}

	return &settings, nil
}

// Update implements SettingsRepository.
func (s *settingsRepository) Update(settings *entity.Settings) error {
	return s.db.Save(settings).Error
}
