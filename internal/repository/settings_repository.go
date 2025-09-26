package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type SettingsRepository interface {
	Create(ctx context.Context, req *entity.Settings) error
	FindByName(ctx context.Context, name string) (*entity.Settings, error)
	Update(ctx context.Context, settings *entity.Settings) error
	Delete(ctx context.Context, settings *entity.Settings) error
	FindAll(ctx context.Context) ([]*entity.Settings, error)
	FindByID(ctx context.Context, id int) (*entity.Settings, error)
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
func (s *settingsRepository) Create(ctx context.Context, req *entity.Settings) error {
	return s.db.WithContext(ctx).Create(req).Error
}

// Delete implements SettingsRepository.
func (s *settingsRepository) Delete(ctx context.Context, settings *entity.Settings) error {
	return s.db.WithContext(ctx).Delete(settings).Error
}

// FindAll implements SettingsRepository.
func (s *settingsRepository) FindAll(ctx context.Context) ([]*entity.Settings, error) {
	var settings []*entity.Settings
	if err := s.db.WithContext(ctx).Find(&settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

// FindByID implements SettingsRepository.
func (s *settingsRepository) FindByID(ctx context.Context, id int) (*entity.Settings, error) {
	var settings entity.Settings
	if err := s.db.WithContext(ctx).First(&settings, id).Error; err != nil {
		return nil, err
	}

	return &settings, nil
}

// FindByName implements SettingsRepository.
func (s *settingsRepository) FindByName(ctx context.Context, name string) (*entity.Settings, error) {
	var settings entity.Settings
	if err := s.db.WithContext(ctx).Where("name = ?", name).First(&settings).Error; err != nil {
		return nil, err
	}

	return &settings, nil
}

// Update implements SettingsRepository.
func (s *settingsRepository) Update(ctx context.Context, settings *entity.Settings) error {
	return s.db.WithContext(ctx).Save(settings).Error
}
