package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
)

type AuthHandler struct {
	authService service.AuthService
	validator   validator.Validator
}

func NewAuthHandler(authService service.AuthService, validator validator.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	user, err := h.authService.Register(c.UserContext(), &req)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "uni_users_email") {
			return response.Error(c, fiber.StatusConflict, "Registration failed", fiber.Map{
				"message": "Email already exists",
				"field":   "email",
			})
		}

		return response.Error(c, fiber.StatusInternalServerError, "Registration failed", err.Error())

	}

	return response.Success(c, "user registered succesfully", user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginUserRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	user, token, err := h.authService.Login(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "login failed", err.Error())
	}

	return response.Success(c, "login success", fiber.Map{
		"user":  user,
		"token": token,
	})
}
