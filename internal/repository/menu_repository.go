package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type MenuRepository interface {
	// Define your menu repository methods here
	Create(req *entity.Menu) error
	FindByID(id int64) (*entity.Menu, error)
	FindAll() ([]*entity.Menu, error)
	Update(req *entity.Menu) error
	Delete(id int64) error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

// Create implements MenuRepository.
func (m *menuRepository) Create(req *entity.Menu) error {
	return m.db.Create(req).Error
}

// Delete implements MenuRepository.
func (m *menuRepository) Delete(id int64) error {
	return m.db.Delete(&entity.Menu{}, id).Error
}

// FindAll implements MenuRepository.
func (m *menuRepository) FindAll() ([]*entity.Menu, error) {
	var menus []*entity.Menu
	err := m.db.Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

// FindByID implements MenuRepository.
func (m *menuRepository) FindByID(id int64) (*entity.Menu, error) {
	var menu entity.Menu
	err := m.db.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// Update implements MenuRepository.
func (m *menuRepository) Update(req *entity.Menu) error {
	return m.db.Save(req).Error
}
