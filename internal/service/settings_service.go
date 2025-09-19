package service

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type SettingsService interface {
	Create(req *dto.CreateSettingsRequest) (*entity.Settings, error)
	FindByName(name string) (*entity.Settings, error)
	Update(id int, req *dto.UpdateSettingsRequest) (*entity.Settings, error)
	Delete(id int) (*entity.Settings, error)
	FindAll() ([]*entity.Settings, error)
}

type settingsService struct {
	settingsRepo repository.SettingsRepository
}

func NewSettingsService(settingsRepo repository.SettingsRepository) SettingsService {
	return &settingsService{
		settingsRepo: settingsRepo,
	}
}

// Create implements SettingsService.
func (s *settingsService) Create(req *dto.CreateSettingsRequest) (*entity.Settings, error) {
	settings := &entity.Settings{
		Name:  req.Name,
		Value: req.Value,
	}

	if err := s.settingsRepo.Create(settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// Delete implements SettingsService.
func (s *settingsService) Delete(id int) (*entity.Settings, error) {
	settings, err := s.settingsRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.settingsRepo.Delete(settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// FindAll implements SettingsService.
func (s *settingsService) FindAll() ([]*entity.Settings, error) {
	return s.settingsRepo.FindAll()
}

// GetByName implements SettingsService.
func (s *settingsService) FindByName(name string) (*entity.Settings, error) {
	return s.settingsRepo.FindByName(name)
}

// Update implements SettingsService.
func (s *settingsService) Update(id int, req *dto.UpdateSettingsRequest) (*entity.Settings, error) {
	settings, err := s.settingsRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		settings.Name = req.Name
	}

	if req.Value != "" {
		settings.Value = req.Value
	}

	if err := s.settingsRepo.Update(settings); err != nil {
		return nil, err
	}

	return settings, nil
}
