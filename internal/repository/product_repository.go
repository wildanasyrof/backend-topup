package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, req *entity.Product) error
	FindAll(ctx context.Context) ([]*entity.Product, error)
	FindByID(ctx context.Context, id int) (*entity.Product, error)
	Update(ctx context.Context, req *entity.Product) error
	Delete(ctx context.Context, id int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create implements ProductRepository.
func (p *productRepository) Create(ctx context.Context, req *entity.Product) error {
	return p.db.WithContext(ctx).Create(req).Error
}

// Delete implements ProductRepository.
func (p *productRepository) Delete(ctx context.Context, id int) error {
	return p.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

// FindAll implements ProductRepository.
func (p *productRepository) FindAll(ctx context.Context) ([]*entity.Product, error) {
	var products []*entity.Product
	err := p.db.WithContext(ctx).
		Preload("Prices").
		Find(&products).Error
	return products, err
}

func (p *productRepository) FindByID(ctx context.Context, id int) (*entity.Product, error) {
	var product entity.Product

	if err := p.db.WithContext(ctx).Preload("Prices").First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// Update implements ProductRepository.
func (p *productRepository) Update(ctx context.Context, req *entity.Product) error {
	return p.db.WithContext(ctx).Save(req).Error
}
