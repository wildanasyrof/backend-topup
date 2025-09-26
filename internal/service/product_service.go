package service

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
)

type ProductService interface {
	Create(ctx context.Context, req *dto.ProductCreateRequest) (*entity.Product, error)
	GetAll(ctx context.Context, q dto.ProductListQuery) ([]entity.Product, pagination.Meta, error)
	Update(ctx context.Context, id int, req *dto.ProductUpdateRequest) (*entity.Product, error)
	Delete(ctx context.Context, id int) (*entity.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductRepository(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// Create implements ProductService.
func (p *productService) Create(ctx context.Context, req *dto.ProductCreateRequest) (*entity.Product, error) {
	product := &entity.Product{
		Name:        req.Name,
		CategoryID:  int(req.CategoryID),
		ProviderID:  int64(req.ProviderID),
		Status:      entity.CatStatus(req.Status),
		Description: req.Description,
		ImgUrl:      req.ImgURL,
	}

	if err := p.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// Delete implements ProductService.
func (p *productService) Delete(ctx context.Context, id int) (*entity.Product, error) {
	product, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := p.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	return product, nil
}

// GetAll implements ProductService.
func (p *productService) GetAll(ctx context.Context, q dto.ProductListQuery) ([]entity.Product, pagination.Meta, error) {
	return p.repo.FindAll(ctx, q)
}

// Update implements ProductService.
func (p *productService) Update(ctx context.Context, id int, req *dto.ProductUpdateRequest) (*entity.Product, error) {
	product, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	req.ToEntity(product)

	if err := p.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}
