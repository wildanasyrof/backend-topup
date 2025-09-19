package handler

import (
	"log"

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
	role, ok := c.Locals("role").(string)
	if !ok {
		log.Println("Role not found in context")
		return response.Error(c, fiber.StatusUnauthorized, "Role not found", role)
	}

	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		log.Println("User ID not found in context")
		return response.Error(c, fiber.StatusUnauthorized, "User ID not found", nil)
	}

	user, err := h.userService.GetUserByID(uid)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get user profile", err.Error())
	}

	return response.Success(c, "User profile retrieved successfully", user)
}
