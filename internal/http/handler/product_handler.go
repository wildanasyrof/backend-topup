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

type ProductHandler struct {
	service   service.ProductService
	validator validator.Validator
	storage   storage.LocalStorage
}

func NewProductHandler(service service.ProductService, validator validator.Validator, storage storage.LocalStorage) *ProductHandler {
	return &ProductHandler{
		service:   service,
		validator: validator,
		storage:   storage,
	}
}

func (p *ProductHandler) Create(c *fiber.Ctx) error {
	var req dto.ProductCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := p.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	img, err := c.FormFile("image")

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "image is reequired", err.Error())
	}

	imgUrl, err := utils.UploadImage(img, p.storage)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "file error", err.Error())
	}

	req.ImgURL = imgUrl

	product, err := p.service.Create(&req)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed creating product", err.Error())
	}

	return response.Success(c, "success creaating product", product)
}

func (p *ProductHandler) GetAll(c *fiber.Ctx) error {
	products, err := p.service.GetAll()
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed get all products", err.Error())
	}

	return response.Success(c, "success get all products", products)
}

func (p *ProductHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid parameter id", err.Error())
	}

	var req dto.ProductUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := p.validator.ValidateBody(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	img, err := c.FormFile("image")

	if err == nil && img != nil {
		imgUrl, err := utils.UploadImage(img, p.storage)

		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "file error", err.Error())
		}
		req.ImgURL = &imgUrl
	}

	product, err := p.service.Update(id, &req)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed updating product", err)
	}

	return response.Success(c, "success updating product", product)
}

func (p *ProductHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid parameter id", err.Error())
	}

	product, err := p.service.Delete(id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to delete product", err.Error())
	}

	return response.Success(c, "success deleted product", product)
}
