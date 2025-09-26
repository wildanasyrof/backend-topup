package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type ProviderRepository interface {
	Create(ctx context.Context, req *entity.Provider) error
	FindAll(ctx context.Context) ([]*entity.Provider, error)
	FindByID(ctx context.Context, id int64) (*entity.Provider, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Provider, error)
	Update(ctx context.Context, req *entity.Provider) error
	Delete(ctx context.Context, id int64) error
}

type providerRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) ProviderRepository {
	return &providerRepository{db: db}
}

// Create implements ProviderRepository.
func (p *providerRepository) Create(ctx context.Context, req *entity.Provider) error {
	return p.db.WithContext(ctx).Create(req).Error
}

// Delete implements ProviderRepository.
func (p *providerRepository) Delete(ctx context.Context, id int64) error {
	return p.db.WithContext(ctx).Delete(&entity.Provider{}, id).Error
}

// FindAll implements ProviderRepository.
func (p *providerRepository) FindAll(ctx context.Context) ([]*entity.Provider, error) {
	var providers []*entity.Provider
	err := p.db.WithContext(ctx).Find(&providers).Error
	if err != nil {
		return nil, err
	}

	return providers, nil
}

// FindByID implements ProviderRepository.
func (p *providerRepository) FindByID(ctx context.Context, id int64) (*entity.Provider, error) {
	var provider entity.Provider
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&provider).Error
	if err != nil {
		return nil, gorm.ErrCheckConstraintViolated
	}

	return &provider, nil
}

func (p *providerRepository) FindBySlug(ctx context.Context, slug string) (*entity.Provider, error) {
	var provider entity.Provider
	err := p.db.WithContext(ctx).Where("slug = ?", slug).First(&provider).Error
	if err != nil {
		return nil, err
	}

	return &provider, nil
}

// Update implements ProviderRepository.
func (p *providerRepository) Update(ctx context.Context, req *entity.Provider) error {
	return p.db.WithContext(ctx).Save(req).Error
}
