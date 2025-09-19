package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
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
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	settings, err := h.settingsService.Create(&req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create settings", err.Error())
	}

	return response.Success(c, "Settings created successfully", settings)

}

func (h *SettingsHandler) FindAll(c *fiber.Ctx) error {
	settings, err := h.settingsService.FindAll()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch settings", err.Error())
	}

	return response.Success(c, "Settings fetched successfully", settings)
}

func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	var req dto.UpdateSettingsRequest

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid menu ID", err.Error())
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	settings, err := h.settingsService.Update(id, &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update settings", err.Error())
	}

	return response.Success(c, "Settings updated successfully", settings)
}

func (h *SettingsHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid Settings ID", err.Error())
	}

	settings, err := h.settingsService.Delete(id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete settings", err.Error())
	}

	return response.Success(c, "Settings deleted successfully", settings)
}
