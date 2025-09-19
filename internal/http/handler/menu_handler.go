package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
)

type MenuHandler struct {
	menuService service.MenuService
	validator   validator.Validator
}

func NewMenuHandler(menuService service.MenuService, validator validator.Validator) *MenuHandler {
	return &MenuHandler{menuService: menuService, validator: validator}
}

func (h *MenuHandler) GetAll(c *fiber.Ctx) error {
	menus, err := h.menuService.GetAll()

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get menus", err.Error())
	}

	return response.Success(c, "Menus retrieved successfully", menus)
}

func (h *MenuHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid menu ID", err.Error())
	}
	menu, err := h.menuService.GetByID(id)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get menu", err.Error())
	}

	return response.Success(c, "Menu retrieved successfully", menu)
}

func (h *MenuHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateMenuRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation error", err)
	}

	menu, err := h.menuService.Create(&req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create menu", err.Error())
	}

	return response.Success(c, "Menu created successfully", menu)
}

func (h *MenuHandler) Update(c *fiber.Ctx) error {
	var req dto.CreateMenuRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation error", err)
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid menu ID", err.Error())
	}

	menu, errc := h.menuService.Update(id, &req)

	if errc != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update menu", errc.Error())
	}

	return response.Success(c, "Menu updated successfully", menu)
}

func (h *MenuHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid menu ID", err.Error())
	}

	menu, err := h.menuService.Delete(id)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete menu", err.Error())
	}

	return response.Success(c, "Menu deleted successfully", menu)
}
