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

type PaymentMethodsHandler struct {
	service   service.PaymentMethodsService
	validator validator.Validator
	storage   storage.LocalStorage
}

func NewPaymentMethodsHandler(service service.PaymentMethodsService, validator validator.Validator, storage storage.LocalStorage) *PaymentMethodsHandler {
	return &PaymentMethodsHandler{
		service:   service,
		validator: validator,
		storage:   storage,
	}
}

func (h *PaymentMethodsHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePaymentMethodRequest

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

	// ---- 5) Service call ----
	res, err := h.service.Create(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return response.Created(c, res)
}

func (h *PaymentMethodsHandler) GetAll(c *fiber.Ctx) error {
	res, err := h.service.FindAll(c.Context())

	if err != nil {
		return err
	}

	return response.OK(c, res)
}

func (h *PaymentMethodsHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}
	res, err := h.service.FindByID(c.UserContext(), id)

	if err != nil {
		return err
	}

	return response.OK(c, res)
}

func (h *PaymentMethodsHandler) Update(c *fiber.Ctx) error {
	// ---- 0) Parse ID ----
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	var req dto.UpdatePaymentMethodRequest
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

	// ---- 5) Service call ----
	res, err := h.service.Update(c.UserContext(), id, &req)
	if err != nil {
		return err
	}

	return response.OK(c, res)
}

func (h *PaymentMethodsHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	res, err := h.service.Delete(c.UserContext(), id)
	if err != nil {
		return err
	}

	return response.OK(c, res)
}
