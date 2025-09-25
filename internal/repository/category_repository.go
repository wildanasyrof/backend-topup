package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(req *entity.Category) error
	FindAll() ([]*entity.Category, error)
	FindByID(id int64) (*entity.Category, error)
	FindBySlug(slug string) (*entity.Category, error)
	Update(req *entity.Category) error
	Delete(id int64) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create implements CategoryRepository.
func (c *categoryRepository) Create(req *entity.Category) error {
	return c.db.Create(req).Error
}

// Delete implements CategoryRepository.
func (c *categoryRepository) Delete(id int64) error {
	return c.db.Delete(&entity.Category{}, id).Error
}

// FindAll implements CategoryRepository.
func (c *categoryRepository) FindAll() ([]*entity.Category, error) {
	var data []*entity.Category
	if err := c.db.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (c *categoryRepository) FindByID(id int64) (*entity.Category, error) {
	var data entity.Category
	if err := c.db.Preload("Products").Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *categoryRepository) FindBySlug(slug string) (*entity.Category, error) {
	var data entity.Category

	err := c.db.Preload("Products").Where("slug = ?", slug).First(&data).Error

	return &data, err
}

// Update implements CategoryRepository.
func (c *categoryRepository) Update(req *entity.Category) error {
	return c.db.Save(req).Error
}
