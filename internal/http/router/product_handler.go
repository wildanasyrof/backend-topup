package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func ProductRouter(r fiber.Router, di *di.DI) {
	r.Get("/", di.ProductHandler.GetAll)

	r.Use(middleware.Auth(di.Jwt, "admin"))
	r.Post("/", di.ProductHandler.Create)
	r.Put("/:id", di.ProductHandler.Update)
	r.Delete("/:id", di.ProductHandler.Delete)
	r.Get("/df", di.ProductHandler.DFUpdate)
}
