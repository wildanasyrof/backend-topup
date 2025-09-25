package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func OrderRoutes(r fiber.Router, di *di.DI) {
	// Route for GUEST/unauthenticated users
	r.Post("/guest", di.OrderHandler.CreateGuest)
	r.Get("/:ref", di.OrderHandler.GetByRef)

	// Route for LOGGED-IN users (requires authentication)
	r.Use(middleware.Auth(di.Jwt))
	r.Post("/", di.OrderHandler.Create)

	r.Use(middleware.Auth(di.Jwt, "admin"))
	r.Get("/all", di.OrderHandler.GetAll)
}
