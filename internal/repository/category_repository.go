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

type CategoryRepository interface {
	Create(ctx context.Context, req *entity.Category) error
	FindAll(ctx context.Context, q dto.CategoryListQuery) (items []entity.Category, meta pagination.Meta, err error)
	FindByID(ctx context.Context, id int64) (*entity.Category, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Category, error)
	Update(ctx context.Context, req *entity.Category) error
	Delete(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create implements CategoryRepository.
func (c *categoryRepository) Create(ctx context.Context, req *entity.Category) error {
	err := c.db.WithContext(ctx).Create(req).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperror.ErrConflict
	}

	return err
}

// Delete implements CategoryRepository.
func (c *categoryRepository) Delete(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
}

// FindAll implements CategoryRepository.
func (c *categoryRepository) FindAll(ctx context.Context, q dto.CategoryListQuery) (items []entity.Category, meta pagination.Meta, err error) {
	q.Normalize() // Normalisasi page dan limit

	base := c.db.WithContext(ctx).Model(&entity.Category{})

	// Tentukan kolom yang boleh di-sort
	allowedSort := map[string]struct{}{"created_at": {}, "name": {}, "id": {}}

	// Terapkan filter dan search
	filtered := base.
		Scopes(
			CategoryFilters(q.MenuID, q.ProviderID, q.Status, q.Type),  // Filter spesifik
			ILike([]string{"categories.name", "categories.slug"}, q.Q), // Search
		)

	// Hitung total data setelah filter
	var total int64
	if err = filtered.Count(&total).Error; err != nil {
		return
	}

	// Ambil data untuk halaman ini
	if err = filtered.
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

func (c *categoryRepository) FindByID(ctx context.Context, id int64) (*entity.Category, error) {
	var data entity.Category
	err := c.db.WithContext(ctx).Preload("Products").Where("id = ?", id).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}

	return &data, err
}

func (c *categoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	var data entity.Category

	err := c.db.WithContext(ctx).Preload("Products").Where("slug = ?", slug).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}
	return &data, err
}

// Update implements CategoryRepository.
func (c *categoryRepository) Update(ctx context.Context, req *entity.Category) error {
	return c.db.WithContext(ctx).Save(req).Error
}

func CategoryFilters(menuID, providerID *int64, status, catType *string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if menuID != nil {
			db = db.Where("menu_id = ?", *menuID)
		}
		if providerID != nil {
			db = db.Where("provider_id = ?", *providerID)
		}
		if status != nil {
			db = db.Where("status = ?", *status)
		}
		if catType != nil {
			db = db.Where("type = ?", *catType)
		}
		return db
	}
}
