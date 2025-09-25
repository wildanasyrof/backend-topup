package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(req *entity.Order) error
	FindAll() ([]*entity.Order, error)
	FindByID(id int) (*entity.Order, error)
	FindByRef(ref string) (*entity.Order, error)
	FindByUserID(userId int64) ([]*entity.Order, error)
	Update(req *entity.Order) error
	Delete(id int) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create implements OrderRepository.
func (o *orderRepository) Create(req *entity.Order) error {
	return o.db.Create(req).Error
}

// Delete implements OrderRepository.
func (o *orderRepository) Delete(id int) error {
	return o.db.Delete(&entity.Order{}, id).Error
}

// FindAll implements OrderRepository.
func (o *orderRepository) FindAll() ([]*entity.Order, error) {
	var orders []*entity.Order

	err := o.db.
		Preload("Product").Find(&orders).Error

	return orders, err
}

// FindByID implements OrderRepository.
func (o *orderRepository) FindByID(id int) (*entity.Order, error) {
	var order entity.Order
	err := o.db.Where("id = ?", id).First(&order).Error

	return &order, err
}

// FindByRef implements OrderRepository.
func (o *orderRepository) FindByRef(ref string) (*entity.Order, error) {
	var order entity.Order
	err := o.db.
		Preload("Product").
		Where("order_ref = ?", ref).First(&order).Error

	return &order, err
}

func (o *orderRepository) FindByUserID(userId int64) ([]*entity.Order, error) {
	var orders []*entity.Order

	err := o.db.Preload("Product").Where("user_id = ?", userId).Find(&orders).Error

	return orders, err
}

// Update implements OrderRepository.
func (o *orderRepository) Update(req *entity.Order) error {
	return o.db.Save(req).Error
}
