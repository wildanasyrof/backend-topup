package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/handler"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func MenuRoutes(r fiber.Router, menuHandler *handler.MenuHandler, di *di.DI) {
	r.Get("/", menuHandler.GetAll)
	r.Get("/:id", menuHandler.GetByID)

	r.Use(middleware.Auth(di.Jwt, "admin"))
	r.Post("/", menuHandler.Create)
	r.Put("/:id", menuHandler.Update)
	r.Delete("/:id", menuHandler.Delete)
}
