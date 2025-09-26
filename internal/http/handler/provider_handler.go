package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
)

type ProviderHandler struct {
	service   service.ProviderService
	validator validator.Validator
}

func NewProviderHandler(s service.ProviderService, v validator.Validator) *ProviderHandler {
	return &ProviderHandler{
		service:   s,
		validator: v,
	}
}

func (h *ProviderHandler) Create(c *fiber.Ctx) error {
	var req dto.ProviderRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	provider, err := h.service.Create(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to create menu", err.Error())
	}

	return response.Success(c, "menu created sucessfully", provider)
}

func (h *ProviderHandler) GetAll(c *fiber.Ctx) error {
	providers, err := h.service.GetAll(c.Context())

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to get list prvoider", err.Error())
	}

	return response.Success(c, "success get all provider", providers)
}

func (h *ProviderHandler) Update(c *fiber.Ctx) error {
	var req dto.ProviderUpdate

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid menu ID", err.Error())
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request bddy", err.Error())
	}

	if err := h.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	provider, err := h.service.Update(c.UserContext(), int64(id), &req)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to update provider", err.Error())
	}

	return response.Success(c, "provider updated", provider)

}

func (h *ProviderHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid menu ID", err.Error())
	}

	provider, err := h.service.Delete(c.UserContext(), int64(id))

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to delete provider", err.Error())
	}

	return response.Success(c, "provider deleted", provider)
}
