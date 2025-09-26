package handler

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/storage"
)

type BannerHandler struct {
	bannerSvc service.BannerService
	storage   storage.LocalStorage
}

func NewBannerHandler(s service.BannerService, storage storage.LocalStorage) *BannerHandler {
	return &BannerHandler{
		bannerSvc: s,
		storage:   storage,
	}
}

func (h *BannerHandler) Create(c *fiber.Ctx) error {
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

	// ---- 4) Save file ----
	filename, err := h.storage.Save(file)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to save file", err)
	}
	imgUrl := "/uploads/" + filename

	// ---- 5) Service call ----
	res, err := h.bannerSvc.Create(c.UserContext(), imgUrl)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to upload banner", err.Error())
	}

	return response.Success(c, "banner uploaded successfully", res)
}

func (h *BannerHandler) GetAll(c *fiber.Ctx) error {
	res, err := h.bannerSvc.FindAll(c.Context())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get banners", err.Error())
	}
	return response.Success(c, "banners retrieved successfully", res)
}

func (h *BannerHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid banner id", err.Error())
	}

	banner, err := h.bannerSvc.Delete(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete banner", err.Error())
	}

	return response.Success(c, "banner deleted successfully", banner)
}

func (h *BannerHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid banner id", err.Error())
	}

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

	// ---- 4) Save file ----
	filename, err := h.storage.Save(file)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to save file", err)
	}
	imgUrl := "/uploads/" + filename

	// ---- 5) Service call ----
	res, err := h.bannerSvc.Update(c.UserContext(), id, imgUrl)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to upload banner", err.Error())
	}

	return response.Success(c, "banner uploaded successfully", res)
}
