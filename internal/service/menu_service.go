package service

import (
	"errors"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type MenuService interface {
	Create(req *dto.CreateMenuRequest) (*entity.Menu, error)
	GetAll() ([]*entity.Menu, error)
	GetByID(id int) (*entity.Menu, error)
	Update(id int, req *dto.CreateMenuRequest) (*entity.Menu, error)
	Delete(id int) (*entity.Menu, error)
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
func (m *menuService) Create(req *dto.CreateMenuRequest) (*entity.Menu, error) {
	menu := &entity.Menu{
		Name: req.Name,
	}

	if err := m.menuRepository.Create(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// Delete implements MenuService.
func (m *menuService) Delete(id int) (*entity.Menu, error) {
	menu, err := m.menuRepository.FindByID(int64(id))
	if err != nil {
		return nil, errors.New("menu not found")
	}

	if err := m.menuRepository.Delete(int64(id)); err != nil {
		return nil, err
	}

	return menu, err
}

// GetAll implements MenuService.
func (m *menuService) GetAll() ([]*entity.Menu, error) {
	return m.menuRepository.FindAll()
}

// GetByID implements MenuService.
func (m *menuService) GetByID(id int) (*entity.Menu, error) {
	return m.menuRepository.FindByID(int64(id))
}

// Update implements MenuService.
func (m *menuService) Update(id int, req *dto.CreateMenuRequest) (*entity.Menu, error) {
	menu, err := m.menuRepository.FindByID(int64(id))
	if err != nil {
		return nil, errors.New("menu not found")
	}

	menu.Name = req.Name

	if err := m.menuRepository.Update(menu); err != nil {
		return nil, err
	}

	return menu, nil
}
