package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type PriceRepository interface {
	Create(req *entity.Price) error
	FindAll() ([]*entity.Price, error)
	FindByID(id int) (*entity.Price, error)
	FindByProductIDnUserLevelID(productId int, userLevelId int) (*entity.Price, error)
	Update(price *entity.Price) error
	Delete(id int) error
}

type priceRepository struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) PriceRepository {
	return &priceRepository{db: db}
}

// Create implements PriceRepository.
func (p *priceRepository) Create(req *entity.Price) error {
	return p.db.Create(req).Error
}

// Delete implements PriceRepository.
func (p *priceRepository) Delete(id int) error {
	return p.db.Delete(&entity.Price{}, id).Error
}

// FindAll implements PriceRepository.
func (p *priceRepository) FindAll() ([]*entity.Price, error) {
	var prices []*entity.Price

	if err := p.db.Find(&prices).Error; err != nil {
		return nil, err
	}

	return prices, nil
}

// FindByID implements PriceRepository.
func (p *priceRepository) FindByID(id int) (*entity.Price, error) {
	var price entity.Price

	if err := p.db.Where("id = ?", id).First(&price).Error; err != nil {
		return nil, err
	}

	return &price, nil
}

// Update implements PriceRepository.
func (p *priceRepository) Update(price *entity.Price) error {
	return p.db.Save(price).Error
}

func (p *priceRepository) FindByProductIDnUserLevelID(productId int, userLevelId int) (*entity.Price, error) {
	var price entity.Price

	err := p.db.
		Where("product_id = ? AND user_level_id = ?", productId, userLevelId).
		First(&price).Error

	return &price, err
}
