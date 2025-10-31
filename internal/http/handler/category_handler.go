package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
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
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	img, err := c.FormFile("image")

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "image is required", err)
	}

	imgUrl, err := utils.UploadImage(img, h.storage)

	if err != nil {
		return err
	}

	req.ImgUrl = imgUrl

	category, err := h.service.Create(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return response.OK(c, category)
}

func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	// 1. Definisikan var untuk query DTO
	var req dto.CategoryListQuery

	// 2. Parse query parameters (e.g., ?page=1&limit=10&q=game)
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

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	var req dto.UpdateCategoryRequest

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

	img, err := c.FormFile("image")

	if err == nil && img != nil {
		imgUrl, err := utils.UploadImage(img, h.storage)

		if err != nil {
			return err
		}
		req.ImgUrl = imgUrl
	}

	category, err := h.service.Update(c.UserContext(), int64(id), &req)
	if err != nil {
		return err
	}

	return response.OK(c, category)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	category, err := h.service.Delete(c.UserContext(), int64(id))

	if err != nil {
		return err
	}

	return response.OK(c, category)
}

func (h *CategoryHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	category, err := h.service.GetBySlug(c.UserContext(), slug)

	if err != nil {
		return err
	}

	return response.OK(c, category)
}
