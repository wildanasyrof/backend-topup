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

type MenuHandler struct {
	menuService service.MenuService
	validator   validator.Validator
}

func NewMenuHandler(menuService service.MenuService, validator validator.Validator) *MenuHandler {
	return &MenuHandler{menuService: menuService, validator: validator}
}

func (h *MenuHandler) GetAll(c *fiber.Ctx) error {
	// 1. Definisikan var untuk query DTO
	var req dto.MenuListQuery

	// 2. Parse query parameters (e.g., ?page=1&limit=10&q=game)
	if err := c.QueryParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid query parameters", err)
	}

	// 3. Panggil service, sekarang mengembalikan 3 nilai
	items, meta, err := h.menuService.GetAll(c.UserContext(), req)
	if err != nil {
		return err
	}

	// 4. Kembalikan response dengan data dan meta
	return response.OK(c, items, meta)
}

func (h *MenuHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}
	menu, err := h.menuService.GetByID(c.UserContext(), id)

	if err != nil {
		return err
	}

	return response.OK(c, menu)
}

func (h *MenuHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateMenuRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	menu, err := h.menuService.Create(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return response.OK(c, menu)
}

func (h *MenuHandler) Update(c *fiber.Ctx) error {
	var req dto.CreateMenuRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	menu, errc := h.menuService.Update(c.UserContext(), id, &req)

	if errc != nil {
		return err
	}

	return response.OK(c, menu)
}

func (h *MenuHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	menu, err := h.menuService.Delete(c.UserContext(), id)

	if err != nil {
		return err
	}

	return response.OK(c, menu)
}
