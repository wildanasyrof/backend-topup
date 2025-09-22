package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type ProviderRepository interface {
	Create(req *entity.Provider) error
	FindAll() ([]*entity.Provider, error)
	FindByID(id int64) (*entity.Provider, error)
	FindBySlug(slug string) (*entity.Provider, error)
	Update(req *entity.Provider) error
	Delete(id int64) error
}

type providerRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) ProviderRepository {
	return &providerRepository{db: db}
}

// Create implements ProviderRepository.
func (p *providerRepository) Create(req *entity.Provider) error {
	return p.db.Create(req).Error
}

// Delete implements ProviderRepository.
func (p *providerRepository) Delete(id int64) error {
	return p.db.Delete(&entity.Provider{}, id).Error
}

// FindAll implements ProviderRepository.
func (p *providerRepository) FindAll() ([]*entity.Provider, error) {
	var providers []*entity.Provider
	err := p.db.Find(&providers).Error
	if err != nil {
		return nil, err
	}

	return providers, nil
}

// FindByID implements ProviderRepository.
func (p *providerRepository) FindByID(id int64) (*entity.Provider, error) {
	var provider entity.Provider
	err := p.db.Where("id = ?", id).First(&provider).Error
	if err != nil {
		return nil, gorm.ErrCheckConstraintViolated
	}

	return &provider, nil
}

func (p *providerRepository) FindBySlug(slug string) (*entity.Provider, error) {
	var provider entity.Provider
	err := p.db.Where("slug = ?", slug).First(&provider).Error
	if err != nil {
		return nil, err
	}

	return &provider, nil
}

// Update implements ProviderRepository.
func (p *providerRepository) Update(req *entity.Provider) error {
	return p.db.Save(req).Error
}
