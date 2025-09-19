package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
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
		log.Println("User ID not found in context")
		return response.Error(c, fiber.StatusUnauthorized, "User ID not found", nil)
	}

	user, err := h.userService.GetUserByID(uid)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get user profile", err.Error())
	}

	return response.Success(c, "User profile retrieved successfully", user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var req dto.UpdateUserRequest

	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		log.Println("User ID not found in context")
		return response.Error(c, fiber.StatusUnauthorized, "User ID not found", nil)
	}

	if err := c.BodyParser(&req); err != nil {
		log.Println("Failed to parse request body:", err)
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(&req); err != nil {
		log.Println("Validation error:", err)
		return response.Error(c, fiber.StatusBadRequest, "Validation error", err)
	}

	user, err := h.userService.Update(uid, &req)
	if err != nil {
		log.Println("Failed to update user:", err)
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update user", err.Error())
	}

	return response.Success(c, "User updated successfully", user)

}
