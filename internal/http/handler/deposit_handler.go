package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/logger"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
)

type DepositHandler struct {
	DepositSvc service.DepositService
	validator  validator.Validator
	Logger     logger.Logger
}

func NewDepositHandler(depoSvc service.DepositService, validator validator.Validator, logger logger.Logger) *DepositHandler {
	return &DepositHandler{
		DepositSvc: depoSvc,
		validator:  validator,
		Logger:     logger,
	}
}

func (h *DepositHandler) Create(c *fiber.Ctx) error {
	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		return apperror.ErrUnauthorized
	}

	var req dto.DepositRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	deposit, err := h.DepositSvc.Create(c.UserContext(), uid, &req)
	if err != nil {
		return err
	}

	return response.OK(c, deposit)
}

func (h *DepositHandler) GetByUserID(c *fiber.Ctx) error {
	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		return apperror.ErrUnauthorized
	}

	deposits, err := h.DepositSvc.GetByUserID(c.UserContext(), uid)
	if err != nil {
		return err
	}

	return response.OK(c, deposits)
}

func (h *DepositHandler) GetByDepositID(c *fiber.Ctx) error {
	depositID := c.Query("id")

	h.Logger.Debug(depositID)

	deposit, err := h.DepositSvc.GetByDepositID(c.UserContext(), depositID)
	if err != nil {
		return err
	}

	return response.OK(c, deposit)
}
