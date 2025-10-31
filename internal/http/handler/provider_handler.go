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
	// 1. Definisikan var untuk query DTO
	var req dto.ProviderListQuery

	// 2. Parse query parameters (e.g., ?page=1&limit=10&q=digiflazz)
	if err := c.QueryParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid query parameters", err)
	}

	// 3. Panggil service, sekarang mengembalikan 3 nilai
	items, meta, err := h.service.GetAll(c.UserContext(), req)
	if err != nil {
		return err
	}

	// 4. Kembalikan response dengan data dan meta
	return response.OK(c, items, meta)
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
