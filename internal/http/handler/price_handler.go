package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
)

type PriceHandler struct {
	service   service.PriceService
	validator validator.Validator
}

func NewPriceHandler(service service.PriceService, validator validator.Validator) *PriceHandler {
	return &PriceHandler{service: service, validator: validator}
}

func (p *PriceHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePrice

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := p.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	price, err := p.service.Create(&req)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to create price", err.Error())
	}

	return response.Success(c, "success createing price", price)
}

func (p *PriceHandler) GetAll(c *fiber.Ctx) error {
	prices, err := p.service.GetAll()

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to get prices", err.Error())
	}

	return response.Success(c, "success get all price", prices)
}

func (p *PriceHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid params id", err.Error())
	}

	var req dto.UpdatePrice

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := p.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	price, err := p.service.Update(id, &req)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to update price", err.Error())
	}

	return response.Success(c, "success updating price", price)
}

func (p *PriceHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid params id", err.Error())
	}

	price, err := p.service.Delete(id)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed deleting price", err.Error())
	}

	return response.Success(c, "success deleting price", price)
}
