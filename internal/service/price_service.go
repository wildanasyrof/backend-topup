package service

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type PriceService interface {
	Create(req *dto.CreatePrice) (*entity.Price, error)
	GetAll() ([]*entity.Price, error)
	Update(id int, req *dto.UpdatePrice) (*entity.Price, error)
	Delete(id int) (*entity.Price, error)
}

type priceService struct {
	repo repository.PriceRepository
}

func NewPriceService(repo repository.PriceRepository) PriceService {
	return &priceService{repo: repo}
}

// Create implements PriceService.
func (p *priceService) Create(req *dto.CreatePrice) (*entity.Price, error) {
	price := &entity.Price{
		ProductID:   req.ProductID,
		UserLevelID: req.UserLevelID,
		Price:       req.Price,
	}

	if err := p.repo.Create(price); err != nil {
		return nil, err
	}

	return price, nil
}

// Delete implements PriceService.
func (p *priceService) Delete(id int) (*entity.Price, error) {
	price, err := p.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	if err := p.repo.Delete(id); err != nil {
		return nil, err
	}

	return price, nil
}

// GetAll implements PriceService.
func (p *priceService) GetAll() ([]*entity.Price, error) {
	return p.repo.FindAll()
}

// Update implements PriceService.
func (p *priceService) Update(id int, req *dto.UpdatePrice) (*entity.Price, error) {
	price, err := p.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	req.ToEntity(price)

	if err := p.repo.Update(price); err != nil {
		return nil, err
	}

	return price, nil

}
