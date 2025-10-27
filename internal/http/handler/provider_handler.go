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
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	provider, err := h.service.Create(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return response.Created(c, provider)
}

func (h *ProviderHandler) GetAll(c *fiber.Ctx) error {
	providers, err := h.service.GetAll(c.Context())

	if err != nil {
		return err
	}

	return response.OK(c, providers)
}

func (h *ProviderHandler) Update(c *fiber.Ctx) error {
	var req dto.ProviderUpdate

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	provider, err := h.service.Update(c.UserContext(), int64(id), &req)

	if err != nil {
		return err
	}

	return response.OK(c, provider)

}

func (h *ProviderHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	provider, err := h.service.Delete(c.UserContext(), int64(id))

	if err != nil {
		return err
	}

	return response.OK(c, provider)
}
