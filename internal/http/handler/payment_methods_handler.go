package handler

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/storage"
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

	// ---- 1) Parse form fields ----
	req.Type = c.FormValue("type")
	req.Name = c.FormValue("name")
	req.Provider = c.FormValue("provider")
	req.ProviderCode = c.FormValue("provider_code")

	if v := c.FormValue("fee"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "fee must be a number", nil)
		}
		req.Fee = &f
	}
	if v := c.FormValue("percent"); v != "" {
		p, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "percent must be a number", nil)
		}
		req.Percent = &p
	}

	// ---- 2) File (required) ----
	file, err := c.FormFile("image")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "image is required", err.Error())
	}

	const maxBytes = 2 * 1024 * 1024 // 2MB
	if file.Size > maxBytes {
		return response.Error(c, fiber.StatusBadRequest, "file too large (max 2MB)", nil)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return response.Error(c, fiber.StatusBadRequest, "unsupported image type (jpg/jpeg/png/webp only)", nil)
	}

	// ---- 3) Validate fields ----
	if err := h.validator.ValidateBody(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	// ---- 4) Save file ----
	filename, err := h.storage.Save(file)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to save file", err)
	}
	req.ImgUrl = "/uploads/" + filename

	// ---- 5) Service call ----
	res, err := h.service.Create(&req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create payment method", err.Error())
	}

	return response.Success(c, "payment method created successfully", res)
}

func (h *PaymentMethodsHandler) GetAll(c *fiber.Ctx) error {
	res, err := h.service.FindAll()

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get payment methods", err.Error())
	}

	return response.Success(c, "payment methods retrieved successfully", res)
}
func (h *PaymentMethodsHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id parameter", err.Error())
	}
	res, err := h.service.FindByID(id)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get payment method", err.Error())
	}

	return response.Success(c, "payment method retrieved successfully", res)
}

func (h *PaymentMethodsHandler) Update(c *fiber.Ctx) error {
	// ---- 0) Parse ID ----
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id parameter", err.Error())
	}

	var req dto.UpdatePaymentMethodRequest
	var hasAnyField bool

	// ---- 1) Parse form fields ----
	if v := c.FormValue("type"); v != "" {
		req.Type = v
		hasAnyField = true
	}
	if v := c.FormValue("name"); v != "" {
		req.Name = v
		hasAnyField = true
	}
	if v := c.FormValue("provider"); v != "" {
		req.Provider = v
		hasAnyField = true
	}
	if v := c.FormValue("provider_code"); v != "" {
		req.ProviderCode = v
		hasAnyField = true
	}
	if v := c.FormValue("fee"); v != "" {
		f, perr := strconv.ParseFloat(v, 64)
		if perr != nil {
			return response.Error(c, fiber.StatusBadRequest, "fee must be a number", nil)
		}
		req.Fee = &f
		hasAnyField = true
	}
	if v := c.FormValue("percent"); v != "" {
		p, perr := strconv.ParseFloat(v, 64)
		if perr != nil {
			return response.Error(c, fiber.StatusBadRequest, "percent must be a number", nil)
		}
		req.Percent = &p
		hasAnyField = true
	}

	// ---- 2) Optional image ----
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		hasAnyField = true

		const maxBytes = 2 * 1024 * 1024 // 2MB
		if file.Size > maxBytes {
			return response.Error(c, fiber.StatusBadRequest, "file too large (max 2MB)", nil)
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
		default:
			return response.Error(c, fiber.StatusBadRequest, "unsupported image type (jpg/jpeg/png/webp only)", nil)
		}

		filename, err := h.storage.Save(file)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to save file", err)
		}
		req.ImgUrl = "/uploads/" + filename
	}

	// ---- 3) Require at least one field ----
	if !hasAnyField {
		return response.Error(c, fiber.StatusBadRequest, "no fields to update", nil)
	}

	// ---- 4) Validate ----
	if err := h.validator.ValidateBody(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	// ---- 5) Service call ----
	res, err := h.service.Update(id, &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to update payment method", err.Error())
	}

	return response.Success(c, "payment method updated successfully", res)
}

func (h *PaymentMethodsHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id parameter", err.Error())
	}

	res, err := h.service.Delete(id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete payment method", err.Error())
	}

	return response.Success(c, "payment method deleted successfully", res)
}
