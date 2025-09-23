package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func PriceRoutes(r fiber.Router, di *di.DI) {
	r.Get("/", di.PriceHandler.GetAll)

	r.Use(middleware.Auth(di.Jwt, "admin"))
	r.Post("/", di.PriceHandler.Create)
	r.Put("/:id", di.PriceHandler.Update)
	r.Delete("/:id", di.PriceHandler.Delete)
}
