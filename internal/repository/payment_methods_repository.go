package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type PaymentMethodsRepository interface {
	FindAll() ([]*entity.PaymentMethod, error)
	Create(paymentMethod *entity.PaymentMethod) error
	Update(paymentMethod *entity.PaymentMethod) error
	Delete(id uint64) error
	FindByID(id uint64) (*entity.PaymentMethod, error)
}

type paymentMethodsRepository struct {
	db *gorm.DB
}

func NewPaymentMethodsRepository(db *gorm.DB) PaymentMethodsRepository {
	return &paymentMethodsRepository{db: db}
}

// Create implements PaymentMethodsRepisitory.
func (p *paymentMethodsRepository) Create(paymentMethod *entity.PaymentMethod) error {
	return p.db.Create(paymentMethod).Error
}

// Delete implements PaymentMethodsRepisitory.
func (p *paymentMethodsRepository) Delete(id uint64) error {
	return p.db.Delete(&entity.PaymentMethod{}, id).Error
}

// FindAll implements PaymentMethodsRepisitory.
func (p *paymentMethodsRepository) FindAll() ([]*entity.PaymentMethod, error) {
	var data []*entity.PaymentMethod
	if err := p.db.Preload("Provider").Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// FindByID implements PaymentMethodsRepisitory.
func (p *paymentMethodsRepository) FindByID(id uint64) (*entity.PaymentMethod, error) {
	var data entity.PaymentMethod
	if err := p.db.First(&data, id).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

// Update implements PaymentMethodsRepisitory.
func (p *paymentMethodsRepository) Update(paymentMethod *entity.PaymentMethod) error {
	return p.db.Save(paymentMethod).Error
}
