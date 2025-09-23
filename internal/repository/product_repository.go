package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(req *entity.Product) error
	FindAll() ([]*entity.Product, error)
	FindByID(id int) (*entity.Product, error)
	Update(req *entity.Product) error
	Delete(id int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create implements ProductRepository.
func (p *productRepository) Create(req *entity.Product) error {
	return p.db.Create(req).Error
}

// Delete implements ProductRepository.
func (p *productRepository) Delete(id int) error {
	return p.db.Delete(&entity.Product{}, id).Error
}

// FindAll implements ProductRepository.
func (p *productRepository) FindAll() ([]*entity.Product, error) {
	var products []*entity.Product

	if err := p.db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productRepository) FindByID(id int) (*entity.Product, error) {
	var product entity.Product

	if err := p.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// Update implements ProductRepository.
func (p *productRepository) Update(req *entity.Product) error {
	return p.db.Save(req).Error
}
