package repository

import (
	"context"
	"errors"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, req *entity.Category) error
	FindAll(ctx context.Context) ([]*entity.Category, error)
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
func (c *categoryRepository) FindAll(ctx context.Context) ([]*entity.Category, error) {
	var data []*entity.Category
	if err := c.db.WithContext(ctx).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
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
