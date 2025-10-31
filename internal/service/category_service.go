package service

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
)

type CategoryService interface {
	Create(ctx context.Context, req *dto.CreateCategoryRequest) (*entity.Category, error)
	GetAll(ctx context.Context, q dto.CategoryListQuery) ([]entity.Category, pagination.Meta, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	Update(ctx context.Context, id int64, req *dto.UpdateCategoryRequest) (*entity.Category, error)
	Delete(ctx context.Context, id int64) (*entity.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// Create implements CategoryService.
func (c *categoryService) Create(ctx context.Context, req *dto.CreateCategoryRequest) (*entity.Category, error) {
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

	if err := c.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete implements CategoryService.
func (c *categoryService) Delete(ctx context.Context, id int64) (*entity.Category, error) {
	category, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := c.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	return category, nil
}

// GetAll implements CategoryService.
func (c *categoryService) GetAll(ctx context.Context, q dto.CategoryListQuery) ([]entity.Category, pagination.Meta, error) {
	return c.repo.FindAll(ctx, q)
}

func (c *categoryService) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	category, err := c.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// Update implements CategoryService.
func (c *categoryService) Update(ctx context.Context, id int64, req *dto.UpdateCategoryRequest) (*entity.Category, error) {
	category, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	req.UpdateEntity(category)

	if err := c.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}
