package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/service"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/storage"
	"github.com/wildanasyrof/backend-topup/pkg/utils"
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

	img, err := c.FormFile("image")

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "image is required", err)
	}

	imgUrl, err := utils.UploadImage(img, h.storage)

	if err != nil {
		return err
	}

	// ---- 5) Service call ----
	res, err := h.bannerSvc.Create(c.UserContext(), imgUrl)
	if err != nil {
		return err
	}

	return response.OK(c, res)
}

func (h *BannerHandler) GetAll(c *fiber.Ctx) error {
	res, err := h.bannerSvc.FindAll(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *BannerHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	banner, err := h.bannerSvc.Delete(c.UserContext(), id)
	if err != nil {
		return err
	}

	return response.OK(c, banner)
}

func (h *BannerHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid request param", err)
	}

	file, err := c.FormFile("image")
	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "image is required", err)
	}

	imgUrl, err := utils.UploadImage(file, h.storage)

	if err != nil {
		return err
	}

	// ---- 5) Service call ----
	res, err := h.bannerSvc.Update(c.UserContext(), id, imgUrl)
	if err != nil {
		return err
	}

	return response.OK(c, res)
}
