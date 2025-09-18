package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/http/handler"
)

func UserRoutes(r fiber.Router, userHandler *handler.UserHandler) {
	r.Get("/", userHandler.GetProfile)
}
