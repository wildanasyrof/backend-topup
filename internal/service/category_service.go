package service

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type CategoryService interface {
	Create(req *dto.CreateCategoryRequest) (*entity.Category, error)
	GetAll() ([]*entity.Category, error)
	GetBySlug(slug string) (*entity.Category, error)
	Update(id int64, req *dto.UpdateCategoryRequest) (*entity.Category, error)
	Delete(id int64) (*entity.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// Create implements CategoryService.
func (c *categoryService) Create(req *dto.CreateCategoryRequest) (*entity.Category, error) {
	category := &entity.Category{
		Name:        req.Name,
		Type:        entity.Type(req.Type),
		MenuID:      req.MenuID,
		ProviderID:  req.ProviderID,
		Slug:        req.Slug,
		Status:      entity.CatStatus(req.Status),
		Description: req.Description,
		InputType:   req.InputType,
		ImgUrl:      req.ImgUrl,
		IsLogin:     req.IsLogin,
	}

	if err := c.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete implements CategoryService.
func (c *categoryService) Delete(id int64) (*entity.Category, error) {
	category, err := c.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := c.repo.Delete(id); err != nil {
		return nil, err
	}

	return category, nil
}

// GetAll implements CategoryService.
func (c *categoryService) GetAll() ([]*entity.Category, error) {
	return c.repo.FindAll()
}

func (c *categoryService) GetBySlug(slug string) (*entity.Category, error) {
	category, err := c.repo.FindBySlug(slug)

	if err != nil {
		return nil, err
	}

	return category, nil
}

// Update implements CategoryService.
func (c *categoryService) Update(id int64, req *dto.UpdateCategoryRequest) (*entity.Category, error) {
	category, err := c.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	req.UpdateEntity(category)

	if err := c.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil

}
