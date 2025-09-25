package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/storage"
	"github.com/wildanasyrof/backend-topup/pkg/utils"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
)

type CategoryHandler struct {
	service   service.CategoryService
	validator validator.Validator
	storage   storage.LocalStorage
}

func NewCategoryHandler(service service.CategoryService, validator validator.Validator, storage storage.LocalStorage) *CategoryHandler {
	return &CategoryHandler{
		service:   service,
		validator: validator,
		storage:   storage,
	}
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	img, err := c.FormFile("image")

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "image is reequired", err.Error())
	}

	imgUrl, err := utils.UploadImage(img, h.storage)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "file error", err.Error())
	}

	req.ImgUrl = imgUrl

	category, err := h.service.Create(&req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed create category", err.Error())
	}

	return response.Success(c, "success creating category", category)
}

func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {

	category, err := h.service.GetAll()
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to get category", err.Error())
	}

	return response.Success(c, "success get all category", category)
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	var req dto.UpdateCategoryRequest

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid category id", err.Error())
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	img, err := c.FormFile("image")

	if err == nil && img != nil {
		imgUrl, err := utils.UploadImage(img, h.storage)

		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "file error", err.Error())
		}
		req.ImgUrl = imgUrl
	}

	category, err := h.service.Update(int64(id), &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "error updating category", err.Error())
	}

	return response.Success(c, "success update category", category)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid category id", err.Error())
	}

	category, err := h.service.Delete(int64(id))

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "error deleting category", err.Error())
	}

	return response.Success(c, "success deleting category", category)
}

func (h *CategoryHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	category, err := h.service.GetBySlug(slug)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to get category", err)
	}

	return response.Success(c, "success get category", category)
}
