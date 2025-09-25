package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func CategoryRotues(r fiber.Router, di *di.DI) {
	r.Get("/", di.CategoryHandler.GetAll)
	r.Get("/:slug", di.CategoryHandler.GetBySlug)

	r.Use(middleware.Auth(di.Jwt, "admin"))
	r.Post("/", di.CategoryHandler.Create)
	r.Put("/:id", di.CategoryHandler.Update)
	r.Delete("/:id", di.CategoryHandler.Delete)
}
