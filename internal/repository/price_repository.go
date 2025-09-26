package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type PriceRepository interface {
	Create(ctx context.Context, req *entity.Price) error
	FindAll(ctx context.Context) ([]*entity.Price, error)
	FindByID(ctx context.Context, id int) (*entity.Price, error)
	FindByProductIDnUserLevelID(ctx context.Context, productId int, userLevelId int) (*entity.Price, error)
	Update(ctx context.Context, price *entity.Price) error
	Delete(ctx context.Context, id int) error
}

type priceRepository struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) PriceRepository {
	return &priceRepository{db: db}
}

// Create implements PriceRepository.
func (p *priceRepository) Create(ctx context.Context, req *entity.Price) error {
	return p.db.WithContext(ctx).Create(req).Error
}

// Delete implements PriceRepository.
func (p *priceRepository) Delete(ctx context.Context, id int) error {
	return p.db.WithContext(ctx).Delete(&entity.Price{}, id).Error
}

// FindAll implements PriceRepository.
func (p *priceRepository) FindAll(ctx context.Context) ([]*entity.Price, error) {
	var prices []*entity.Price

	if err := p.db.WithContext(ctx).Find(&prices).Error; err != nil {
		return nil, err
	}

	return prices, nil
}

// FindByID implements PriceRepository.
func (p *priceRepository) FindByID(ctx context.Context, id int) (*entity.Price, error) {
	var price entity.Price

	if err := p.db.WithContext(ctx).Where("id = ?", id).First(&price).Error; err != nil {
		return nil, err
	}

	return &price, nil
}

// Update implements PriceRepository.
func (p *priceRepository) Update(ctx context.Context, price *entity.Price) error {
	return p.db.WithContext(ctx).Save(price).Error
}

func (p *priceRepository) FindByProductIDnUserLevelID(ctx context.Context, productId int, userLevelId int) (*entity.Price, error) {
	var price entity.Price

	err := p.db.WithContext(ctx).
		Where("product_id = ? AND user_level_id = ?", productId, userLevelId).
		First(&price).Error

	return &price, err
}
