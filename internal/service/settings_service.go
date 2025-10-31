package service

import (
	"context" // Import the context package

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
)

// SettingsService interface updated to include context.Context
type SettingsService interface {
	Create(ctx context.Context, req *dto.CreateSettingsRequest) (*entity.Settings, error)
	FindByName(ctx context.Context, name string) (*entity.Settings, error)
	Update(ctx context.Context, id int, req *dto.UpdateSettingsRequest) (*entity.Settings, error)
	Delete(ctx context.Context, id int) (*entity.Settings, error)
	FindAll(ctx context.Context, q dto.SettingsListQuery) ([]*entity.Settings, pagination.Meta, error)
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
func (s *settingsService) Create(ctx context.Context, req *dto.CreateSettingsRequest) (*entity.Settings, error) {
	settings := &entity.Settings{
		Name:  req.Name,
		Value: req.Value,
	}

	// Pass ctx to the repository call
	if err := s.settingsRepo.Create(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// Delete implements SettingsService.
func (s *settingsService) Delete(ctx context.Context, id int) (*entity.Settings, error) {
	// Pass ctx to the repository call
	settings, err := s.settingsRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Pass ctx to the repository call
	if err := s.settingsRepo.Delete(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// FindAll implements SettingsService.
func (s *settingsService) FindAll(ctx context.Context, q dto.SettingsListQuery) ([]*entity.Settings, pagination.Meta, error) {
	return s.settingsRepo.FindAll(ctx, q)
}

// GetByName implements SettingsService.
func (s *settingsService) FindByName(ctx context.Context, name string) (*entity.Settings, error) {
	// Pass ctx to the repository call
	return s.settingsRepo.FindByName(ctx, name)
}

// Update implements SettingsService.
func (s *settingsService) Update(ctx context.Context, id int, req *dto.UpdateSettingsRequest) (*entity.Settings, error) {
	// Pass ctx to the repository call
	settings, err := s.settingsRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		settings.Name = req.Name
	}

	if req.Value != "" {
		settings.Value = req.Value
	}

	// Pass ctx to the repository call
	if err := s.settingsRepo.Update(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}
