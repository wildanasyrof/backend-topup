package repository

import (
	"context"
	"errors"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(ctx context.Context, req *entity.Menu) error
	FindByID(ctx context.Context, id int64) (*entity.Menu, error)
	FindAll(ctx context.Context) ([]*entity.Menu, error)
	Update(ctx context.Context, req *entity.Menu) error
	Delete(ctx context.Context, id int64) error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (m *menuRepository) Create(ctx context.Context, req *entity.Menu) error {
	err := m.db.WithContext(ctx).Create(req).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperror.ErrConflict
	}

	return err
}

func (m *menuRepository) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&entity.Menu{}, id).Error
}

func (m *menuRepository) FindAll(ctx context.Context) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	err := m.db.WithContext(ctx).
		Preload("Categories").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (m *menuRepository) FindByID(ctx context.Context, id int64) (*entity.Menu, error) {
	var menu entity.Menu
	err := m.db.WithContext(ctx).Preload("Categories").Where("id = ?", id).First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}
	return &menu, nil
}

func (m *menuRepository) Update(ctx context.Context, req *entity.Menu) error {
	return m.db.WithContext(ctx).Save(req).Error
}
