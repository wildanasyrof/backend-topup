package service

import (
	"context"
	"errors"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
)

type MenuService interface {
	Create(ctx context.Context, req *dto.CreateMenuRequest) (*entity.Menu, error)
	GetAll(ctx context.Context, q dto.MenuListQuery) ([]*entity.Menu, pagination.Meta, error)
	GetByID(ctx context.Context, id int) (*entity.Menu, error)
	Update(ctx context.Context, id int, req *dto.CreateMenuRequest) (*entity.Menu, error)
	Delete(ctx context.Context, id int) (*entity.Menu, error)
}

type menuService struct {
	menuRepository repository.MenuRepository
}

func NewMenuService(menuRepo repository.MenuRepository) MenuService {
	return &menuService{
		menuRepository: menuRepo,
	}
}

// Create implements MenuService.
func (m *menuService) Create(ctx context.Context, req *dto.CreateMenuRequest) (*entity.Menu, error) {
	menu := &entity.Menu{
		Name: req.Name,
	}

	if err := m.menuRepository.Create(ctx, menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// Delete implements MenuService.
func (m *menuService) Delete(ctx context.Context, id int) (*entity.Menu, error) {
	menu, err := m.menuRepository.FindByID(ctx, int64(id))
	if err != nil {
		return nil, errors.New("menu not found")
	}

	if err := m.menuRepository.Delete(ctx, int64(id)); err != nil {
		return nil, err
	}

	return menu, err
}

// GetAll implements MenuService.
func (m *menuService) GetAll(ctx context.Context, q dto.MenuListQuery) ([]*entity.Menu, pagination.Meta, error) {
	return m.menuRepository.FindAll(ctx, q)
}

// GetByID implements MenuService.
func (m *menuService) GetByID(ctx context.Context, id int) (*entity.Menu, error) {
	return m.menuRepository.FindByID(ctx, int64(id))
}

// Update implements MenuService.
func (m *menuService) Update(ctx context.Context, id int, req *dto.CreateMenuRequest) (*entity.Menu, error) {
	menu, err := m.menuRepository.FindByID(ctx, int64(id))
	if err != nil {
		return nil, errors.New("menu not found")
	}

	menu.Name = req.Name

	if err := m.menuRepository.Update(ctx, menu); err != nil {
		return nil, err
	}

	return menu, nil
}
