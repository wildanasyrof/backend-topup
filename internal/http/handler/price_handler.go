package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
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
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := p.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	price, err := p.service.Create(c.UserContext(), &req)

	if err != nil {
		return err
	}

	return response.Created(c, price)
}

func (p *PriceHandler) GetAll(c *fiber.Ctx) error {
	prices, err := p.service.GetAll(c.Context())

	if err != nil {
		return err
	}

	return response.OK(c, prices)
}

func (p *PriceHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	var req dto.UpdatePrice

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := p.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	price, err := p.service.Update(c.UserContext(), id, &req)

	if err != nil {
		return err
	}

	return response.OK(c, price)
}

func (p *PriceHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	price, err := p.service.Delete(c.UserContext(), id)

	if err != nil {
		return err
	}

	return response.OK(c, price)
}
