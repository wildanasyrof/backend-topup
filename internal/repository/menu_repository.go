package repository

import (
	"context"
	"errors"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(ctx context.Context, req *entity.Menu) error
	FindByID(ctx context.Context, id int64) (*entity.Menu, error)
	FindAll(ctx context.Context, q dto.MenuListQuery) (items []*entity.Menu, meta pagination.Meta, err error)
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

func (m *menuRepository) FindAll(ctx context.Context, q dto.MenuListQuery) (items []*entity.Menu, meta pagination.Meta, err error) {
	q.Normalize() // Terapkan DefaultPage dan DefaultLimit

	base := m.db.WithContext(ctx).Model(&entity.Menu{})

	// Tentukan kolom yang boleh di-sort
	allowedSort := map[string]struct{}{"created_at": {}, "name": {}, "id": {}}

	// Terapkan search
	filtered := base.
		Scopes(
			ILike([]string{"menus.name"}, q.Q), // Search by name
		)

	// Hitung total data
	var total int64
	if err = filtered.Count(&total).Error; err != nil {
		return
	}

	// Ambil data untuk halaman ini (termasuk Preload yang sudah ada)
	if err = filtered.
		Preload("Categories"). // <-- Tetap Preload relasi
		Scopes(
			func(db *gorm.DB) *gorm.DB { return pagination.ScopeSort(db, q.Sort, allowedSort) },
			func(db *gorm.DB) *gorm.DB { return pagination.ScopePaginate(db, q.Page, q.Limit) },
		).
		Find(&items).Error; err != nil {
		return
	}

	// Hitung metadata paginasi
	meta = pagination.CalcMeta(int(total), q.Page, q.Limit)
	return
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
