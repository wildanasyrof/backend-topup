package service

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
)

type PriceService interface {
	Create(ctx context.Context, req *dto.CreatePrice) (*entity.Price, error)
	GetAll(ctx context.Context, q dto.PriceListQuery) ([]*entity.Price, pagination.Meta, error)
	Update(ctx context.Context, id int, req *dto.UpdatePrice) (*entity.Price, error)
	Delete(ctx context.Context, id int) (*entity.Price, error)
}

type priceService struct {
	repo repository.PriceRepository
}

func NewPriceService(repo repository.PriceRepository) PriceService {
	return &priceService{repo: repo}
}

// Create implements PriceService.
func (p *priceService) Create(ctx context.Context, req *dto.CreatePrice) (*entity.Price, error) {
	price := &entity.Price{
		ProductID:   req.ProductID,
		UserLevelID: req.UserLevelID,
		Price:       req.Price,
	}

	if err := p.repo.Create(ctx, price); err != nil {
		return nil, err
	}

	return price, nil
}

// Delete implements PriceService.
func (p *priceService) Delete(ctx context.Context, id int) (*entity.Price, error) {
	price, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := p.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	return price, nil
}

// GetAll implements PriceService.
func (p *priceService) GetAll(ctx context.Context, q dto.PriceListQuery) ([]*entity.Price, pagination.Meta, error) {
	return p.repo.FindAll(ctx, q)
}

// Update implements PriceService.
func (p *priceService) Update(ctx context.Context, id int, req *dto.UpdatePrice) (*entity.Price, error) {
	price, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	req.ToEntity(price)

	if err := p.repo.Update(ctx, price); err != nil {
		return nil, err
	}

	return price, nil
}
