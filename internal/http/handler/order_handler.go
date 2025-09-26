package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
)

type OrderHandler struct {
	service   service.OrderService
	validator validator.Validator
}

func NewOrderHandler(service service.OrderService, validator validator.Validator) *OrderHandler {
	return &OrderHandler{
		service:   service,
		validator: validator,
	}
}

func (o *OrderHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateOrder
	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "user not found", nil)
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := o.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	order, err := o.service.Create(c.UserContext(), uid, &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed create order", err.Error())
	}

	return response.Success(c, "success create order", order)
}

func (o *OrderHandler) CreateGuest(c *fiber.Ctx) error {
	var req dto.CreateOrder

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := o.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	order, err := o.service.Create(c.UserContext(), 2, &req)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed create order", err.Error())
	}

	return response.Success(c, "order created", order)
}

func (o *OrderHandler) GetByUserID(c *fiber.Ctx) error {
	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "User ID not found", nil)
	}

	orders, err := o.service.GetByUserID(c.UserContext(), uid)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "error get orders list", err)
	}

	return response.Success(c, "success get orders list", orders)
}

func (o *OrderHandler) GetByRef(c *fiber.Ctx) error {
	ref := c.Params("ref")

	order, err := o.service.GetByRef(c.UserContext(), ref)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to get order by ref", err.Error())
	}

	return response.Success(c, "success get order by ref", order)
}

func (o *OrderHandler) GetAll(c *fiber.Ctx) error {
	orders, err := o.service.GetAll(c.Context())

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed get orders", err.Error())
	}

	return response.Success(c, "success get all orders", orders)
}
