package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type PaymentMethodsRepository interface {
	FindAll(ctx context.Context) ([]*entity.PaymentMethod, error)
	Create(ctx context.Context, paymentMethod *entity.PaymentMethod) error
	Update(ctx context.Context, paymentMethod *entity.PaymentMethod) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*entity.PaymentMethod, error)
}

type paymentMethodsRepository struct {
	db *gorm.DB
}

func NewPaymentMethodsRepository(db *gorm.DB) PaymentMethodsRepository {
	return &paymentMethodsRepository{db: db}
}

func (p *paymentMethodsRepository) Create(ctx context.Context, paymentMethod *entity.PaymentMethod) error {
	return p.db.WithContext(ctx).Create(paymentMethod).Error
}

func (p *paymentMethodsRepository) Delete(ctx context.Context, id uint64) error {
	return p.db.WithContext(ctx).Delete(&entity.PaymentMethod{}, id).Error
}

func (p *paymentMethodsRepository) FindAll(ctx context.Context) ([]*entity.PaymentMethod, error) {
	var data []*entity.PaymentMethod
	if err := p.db.WithContext(ctx).Preload("Provider").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (p *paymentMethodsRepository) FindByID(ctx context.Context, id uint64) (*entity.PaymentMethod, error) {
	var data entity.PaymentMethod
	if err := p.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (p *paymentMethodsRepository) Update(ctx context.Context, paymentMethod *entity.PaymentMethod) error {
	return p.db.WithContext(ctx).Save(paymentMethod).Error
}
