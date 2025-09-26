package service

import (
	"context"
	"errors"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/logger"
	"github.com/wildanasyrof/backend-topup/pkg/utils"
)

type OrderService interface {
	Create(ctx context.Context, userId uint64, req *dto.CreateOrder) (*entity.Order, error)
	GetAll(ctx context.Context) ([]*entity.Order, error)
	GetByRef(ctx context.Context, ref string) (*entity.Order, error)
	Update(ctx context.Context, ref string, req *dto.UpdateOrder) (*entity.Order, error)
	GetByUserID(ctx context.Context, userId uint64) ([]*entity.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
	priceRepo repository.PriceRepository
	logger    logger.Logger
}

func NewOrderService(orderRepo repository.OrderRepository, logger logger.Logger, userRepo repository.UserRepository, priceRepo repository.PriceRepository) OrderService {
	return &orderService{orderRepo: orderRepo, logger: logger, userRepo: userRepo, priceRepo: priceRepo}
}

// Create implements OrderService.
func (o *orderService) Create(ctx context.Context, userId uint64, req *dto.CreateOrder) (*entity.Order, error) {
	user, err := o.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, errors.New("error processing user")
	}

	price, err := o.priceRepo.FindByProductIDnUserLevelID(ctx, req.ProductID, user.UserLevelID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	order := &entity.Order{
		OrderRef:     utils.GenerateTopupID(),
		UserID:       userId,
		ProductID:    uint64(req.ProductID),
		WA:           req.WA,
		Email:        req.Email,
		CustomerName: req.CustomerName,
		CustomerID:   req.CustomerID,
		PaymentRef:   "link referensi",
		Amount:       price.Price,
		Fee:          500,
	}

	if err := o.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetAll implements OrderService.
func (o *orderService) GetAll(ctx context.Context) ([]*entity.Order, error) {
	return o.orderRepo.FindAll(ctx)
}

// GetByUserID implements OrderService.
func (o *orderService) GetByUserID(ctx context.Context, userId uint64) ([]*entity.Order, error) {
	return o.orderRepo.FindByUserID(ctx, int64(userId))
}

// GetByRef implements OrderService.
func (o *orderService) GetByRef(ctx context.Context, ref string) (*entity.Order, error) {
	return o.orderRepo.FindByRef(ctx, ref)
}

// Update implements OrderService.
func (o *orderService) Update(ctx context.Context, ref string, req *dto.UpdateOrder) (*entity.Order, error) {
	panic("unimplemented")
}
