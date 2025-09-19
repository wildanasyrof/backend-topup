package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
)

type PaymentMethodsHandler struct {
	service   service.PaymentMethodsService
	validator validator.Validator
}

func NewPaymentMethodsHandler(service service.PaymentMethodsService, validator validator.Validator) *PaymentMethodsHandler {
	return &PaymentMethodsHandler{
		service:   service,
		validator: validator,
	}
}

func (h *PaymentMethodsHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePaymentMethodRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	res, err := h.service.Create(&req)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create payment method", err.Error())
	}

	return response.Success(c, "payment method created successfully", res)
}

func (h *PaymentMethodsHandler) GetAll(c *fiber.Ctx) error {
	res, err := h.service.FindAll()

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get payment methods", err.Error())
	}

	return response.Success(c, "payment methods retrieved successfully", res)
}
func (h *PaymentMethodsHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id parameter", err.Error())
	}
	res, err := h.service.FindByID(id)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get payment method", err.Error())
	}

	return response.Success(c, "payment method retrieved successfully", res)
}

func (h *PaymentMethodsHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	var req dto.UpdatePaymentMethodRequest

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id parameter", err.Error())
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	res, err := h.service.Update(id, &req)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to update payment method", err.Error())
	}

	return response.Success(c, "payment method updated successfully", res)
}

func (h *PaymentMethodsHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id parameter", err.Error())
	}

	res, err := h.service.Delete(id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete payment method", err.Error())
	}

	return response.Success(c, "payment method deleted successfully", res)
}
