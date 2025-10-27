package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
)

type UserHandler struct {
	userService service.UserService
	validator   validator.Validator
}

func NewUserHandler(userService service.UserService, validator validator.Validator) *UserHandler {
	return &UserHandler{userService: userService, validator: validator}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		return apperror.ErrUnauthorized
	}

	user, err := h.userService.GetUserByID(c.UserContext(), uid)

	if err != nil {
		return err
	}
	return response.OK(c, user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var req dto.UpdateUserRequest

	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		return apperror.ErrUnauthorized
	}

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	user, err := h.userService.Update(c.UserContext(), uid, &req)
	if err != nil {
		return err
	}

	return response.OK(c, user)

}
