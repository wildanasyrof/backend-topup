package service

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type ProductService interface {
	Create(req *dto.ProductCreateRequest) (*entity.Product, error)
	GetAll() ([]*entity.Product, error)
	Update(id int, req *dto.ProductUpdateRequest) (*entity.Product, error)
	Delete(id int) (*entity.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductRepository(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// Create implements ProductService.
func (p *productService) Create(req *dto.ProductCreateRequest) (*entity.Product, error) {
	product := &entity.Product{
		Name:        req.Name,
		CategoryID:  int(req.CategoryID),
		ProviderID:  int64(req.ProviderID),
		Status:      entity.CatStatus(req.Status),
		Description: req.Description,
		ImgUrl:      req.ImgURL,
	}

	if err := p.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// Delete implements ProductService.
func (p *productService) Delete(id int) (*entity.Product, error) {
	product, err := p.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	if err := p.repo.Delete(id); err != nil {
		return nil, err
	}

	return product, err
}

// GetAll implements ProductService.
func (p *productService) GetAll() ([]*entity.Product, error) {
	return p.repo.FindAll()
}

// Update implements ProductService.
func (p *productService) Update(id int, req *dto.ProductUpdateRequest) (*entity.Product, error) {
	product, err := p.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	req.ToEntity(product)

	if err := p.repo.Update(product); err != nil {
		return nil, err
	}

	return product, err
}
