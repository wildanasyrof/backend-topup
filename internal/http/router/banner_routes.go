package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func BannerRoutes(r fiber.Router, di *di.DI) {
	r.Get("/", di.BannerHandler.GetAll)

	r.Use(middleware.Auth(di.Jwt, "admin"))
	r.Post("/", di.BannerHandler.Create)
	r.Put("/:id", di.BannerHandler.Update)
	r.Delete("/:id", di.BannerHandler.Delete)
}
