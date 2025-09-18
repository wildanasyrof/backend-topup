package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/response"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	user, err := h.userService.GetUserByID(userID)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get user profile", err.Error())
	}

	return response.Success(c, "User profile retrieved successfully", user)
}
