package service

import (
	"context" // Import the context package

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

// ProviderService interface updated to include context.Context
type ProviderService interface {
	Create(ctx context.Context, req *dto.ProviderRequest) (*entity.Provider, error)
	GetAll(ctx context.Context) ([]*entity.Provider, error)
	Update(ctx context.Context, id int64, req *dto.ProviderUpdate) (*entity.Provider, error)
	Delete(ctx context.Context, id int64) (*entity.Provider, error)
}

type providerService struct {
	providerRepo repository.ProviderRepository
}

func NewProviderService(p repository.ProviderRepository) ProviderService {
	return &providerService{providerRepo: p}
}

// Create implements ProviderService.
func (p *providerService) Create(ctx context.Context, req *dto.ProviderRequest) (*entity.Provider, error) {
	provider := &entity.Provider{
		Name: req.Name,
		Ref:  req.Ref,
	}

	// Pass ctx to the repository call
	if err := p.providerRepo.Create(ctx, provider); err != nil {
		return nil, err
	}

	return provider, nil
}

// Delete implements ProviderService.
func (p *providerService) Delete(ctx context.Context, id int64) (*entity.Provider, error) {
	// Pass ctx to the repository call
	provider, err := p.providerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Pass ctx to the repository call
	if err := p.providerRepo.Delete(ctx, id); err != nil {
		return nil, err
	}

	return provider, nil
}

// GetAll implements ProviderService.
func (p *providerService) GetAll(ctx context.Context) ([]*entity.Provider, error) {
	// Pass ctx to the repository call
	return p.providerRepo.FindAll(ctx)
}

// Update implements ProviderService.
func (p *providerService) Update(ctx context.Context, id int64, req *dto.ProviderUpdate) (*entity.Provider, error) {
	// Pass ctx to the repository call
	provider, err := p.providerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		provider.Name = req.Name
	}

	if req.Ref != "" {
		provider.Ref = req.Ref
	}

	// Pass ctx to the repository call
	if err := p.providerRepo.Update(ctx, provider); err != nil {
		return nil, err
	}

	return provider, nil
}
