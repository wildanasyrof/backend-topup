package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, req *entity.Order) error
	FindAll(ctx context.Context) ([]*entity.Order, error)
	FindByID(ctx context.Context, id int) (*entity.Order, error)
	FindByRef(ctx context.Context, ref string) (*entity.Order, error)
	FindByUserID(ctx context.Context, userId int64) ([]*entity.Order, error)
	Update(ctx context.Context, req *entity.Order) error
	Delete(ctx context.Context, id int) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create implements OrderRepository.
func (o *orderRepository) Create(ctx context.Context, req *entity.Order) error {
	return o.db.WithContext(ctx).Create(req).Error
}

// Delete implements OrderRepository.
func (o *orderRepository) Delete(ctx context.Context, id int) error {
	return o.db.WithContext(ctx).Delete(&entity.Order{}, id).Error
}

// FindAll implements OrderRepository.
func (o *orderRepository) FindAll(ctx context.Context) ([]*entity.Order, error) {
	var orders []*entity.Order

	err := o.db.WithContext(ctx).
		Preload("Product").Find(&orders).Error

	return orders, err
}

// FindByID implements OrderRepository.
func (o *orderRepository) FindByID(ctx context.Context, id int) (*entity.Order, error) {
	var order entity.Order
	err := o.db.WithContext(ctx).Where("id = ?", id).First(&order).Error

	return &order, err
}

// FindByRef implements OrderRepository.
func (o *orderRepository) FindByRef(ctx context.Context, ref string) (*entity.Order, error) {
	var order entity.Order
	err := o.db.WithContext(ctx).
		Preload("Product").
		Where("order_ref = ?", ref).First(&order).Error

	return &order, err
}

func (o *orderRepository) FindByUserID(ctx context.Context, userId int64) ([]*entity.Order, error) {
	var orders []*entity.Order

	err := o.db.WithContext(ctx).Preload("Product").Where("user_id = ?", userId).Find(&orders).Error

	return orders, err
}

// Update implements OrderRepository.
func (o *orderRepository) Update(ctx context.Context, req *entity.Order) error {
	return o.db.WithContext(ctx).Save(req).Error
}
