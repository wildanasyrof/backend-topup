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

type ProductHandler struct {
	service   service.ProductService
	exSvc     service.ExternalService
	validator validator.Validator
	storage   storage.LocalStorage
}

func NewProductHandler(service service.ProductService, validator validator.Validator, storage storage.LocalStorage, exSvc service.ExternalService) *ProductHandler {
	return &ProductHandler{
		service:   service,
		validator: validator,
		storage:   storage,
		exSvc:     exSvc,
	}
}

func (p *ProductHandler) Create(c *fiber.Ctx) error {
	var req dto.ProductCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := p.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	img, err := c.FormFile("image")

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "image is required", err)
	}

	imgUrl, err := utils.UploadImage(img, p.storage)

	if err != nil {
		return err
	}

	req.ImgURL = imgUrl

	product, err := p.service.Create(c.UserContext(), &req)

	if err != nil {
		return err
	}

	return response.Created(c, product)
}

func (p *ProductHandler) GetAll(c *fiber.Ctx) error {
	var req dto.ProductListQuery

	if err := c.QueryParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	req.Normalize()

	items, meta, err := p.service.GetAll(c.UserContext(), req)
	if err != nil {
		return err
	}

	return response.OK(c, items, meta)
}

func (p *ProductHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	var req dto.ProductUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := p.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	img, err := c.FormFile("image")

	if err == nil && img != nil {
		imgUrl, err := utils.UploadImage(img, p.storage)

		if err != nil {
			return err
		}
		req.ImgURL = &imgUrl
	}

	product, err := p.service.Update(c.UserContext(), id, &req)

	if err != nil {
		return err
	}

	return response.OK(c, product)
}

func (p *ProductHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	product, err := p.service.Delete(c.UserContext(), id)
	if err != nil {
		return err
	}

	return response.OK(c, product)
}

func (p *ProductHandler) DFUpdate(c *fiber.Ctx) error {
	data, err := p.exSvc.DFSaveProductList(c.UserContext())

	if err != nil {
		return err
	}

	return response.OK(c, data)
}
