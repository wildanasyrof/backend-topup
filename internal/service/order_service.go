package service

import (
	"errors"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/logger"
	"github.com/wildanasyrof/backend-topup/pkg/utils"
)

type OrderService interface {
	Create(userId uint64, req *dto.CreateOrder) (*entity.Order, error)
	GetAll() ([]*entity.Order, error)
	GetByRef(ref string) (*entity.Order, error)
	Update(ref string, req *dto.UpdateOrder) (*entity.Order, error)
	GetByUserID(userId uint64) ([]*entity.Order, error)
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
func (o *orderService) Create(userId uint64, req *dto.CreateOrder) (*entity.Order, error) {

	user, err := o.userRepo.GetByID(userId)
	if err != nil {
		return nil, errors.New("error processing user")
	}

	price, err := o.priceRepo.FindByProductIDnUserLevelID(req.ProductID, user.UserLevelID)
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

	if err := o.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetAll implements OrderService.
func (o *orderService) GetAll() ([]*entity.Order, error) {
	return o.orderRepo.FindAll()
}

// GetAll implements OrderService.
func (o *orderService) GetByUserID(userId uint64) ([]*entity.Order, error) {
	return o.orderRepo.FindByUserID(int64(userId))
}

// GetByRef implements OrderService.
func (o *orderService) GetByRef(ref string) (*entity.Order, error) {
	return o.orderRepo.FindByRef(ref)
}

// Update implements OrderService.
func (o *orderService) Update(ref string, req *dto.UpdateOrder) (*entity.Order, error) {
	panic("unimplemented")
}
