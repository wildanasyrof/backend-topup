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

type SettingsHandler struct {
	settingsService service.SettingsService
	validator       validator.Validator
}

func NewSettingsHandler(settingsService service.SettingsService, validator validator.Validator) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		validator:       validator,
	}
}

func (h *SettingsHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateSettingsRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	settings, err := h.settingsService.Create(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return response.Created(c, settings)
}

func (h *SettingsHandler) FindAll(c *fiber.Ctx) error {
	// 1. Definisikan var untuk query DTO
	var req dto.SettingsListQuery

	// 2. Parse query parameters (e.g., ?page=1&q=site_name)
	if err := c.QueryParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid query parameters", err)
	}

	// 3. Panggil service, sekarang mengembalikan 3 nilai
	// Gunakan c.UserContext() untuk menghormati timeout middleware
	items, meta, err := h.settingsService.FindAll(c.UserContext(), req)
	if err != nil {
		return err
	}

	// 4. Kembalikan response dengan data dan meta
	return response.OK(c, items, meta)
}

func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	var req dto.UpdateSettingsRequest

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

	settings, err := h.settingsService.Update(c.UserContext(), id, &req)
	if err != nil {
		return err
	}

	return response.OK(c, settings)
}

func (h *SettingsHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	settings, err := h.settingsService.Delete(c.UserContext(), id)
	if err != nil {
		return err
	}

	return response.OK(c, settings)
}
