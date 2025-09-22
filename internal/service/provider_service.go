package service

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type ProviderService interface {
	Create(req *dto.ProviderRequest) (*entity.Provider, error)
	GetAll() ([]*entity.Provider, error)
	Update(id int64, req *dto.ProviderUpdate) (*entity.Provider, error)
	Delete(id int64) (*entity.Provider, error)
}

type providerService struct {
	providerRepo repository.ProviderRepository
}

func NewProviderService(p repository.ProviderRepository) ProviderService {
	return &providerService{providerRepo: p}
}

// Create implements ProviderService.
func (p *providerService) Create(req *dto.ProviderRequest) (*entity.Provider, error) {
	provider := &entity.Provider{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := p.providerRepo.Create(provider); err != nil {
		return nil, err
	}

	return provider, nil
}

// Delete implements ProviderService.
func (p *providerService) Delete(id int64) (*entity.Provider, error) {
	provider, err := p.providerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := p.providerRepo.Delete(id); err != nil {
		return nil, err
	}

	return provider, nil
}

// GetAll implements ProviderService.
func (p *providerService) GetAll() ([]*entity.Provider, error) {
	return p.providerRepo.FindAll()
}

// Update implements ProviderService.
func (p *providerService) Update(id int64, req *dto.ProviderUpdate) (*entity.Provider, error) {

	provider, err := p.providerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		provider.Name = req.Name
	}

	if req.Slug != "" {
		provider.Slug = req.Slug
	}

	if err := p.providerRepo.Update(provider); err != nil {
		return nil, err
	}

	return provider, nil
}
