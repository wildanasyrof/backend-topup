package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/handler"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func PaymentMethodsRoutes(r fiber.Router, h *handler.PaymentMethodsHandler, di *di.DI) {
	r.Get("/", h.GetAll)
	r.Get("/:id", h.GetByID)

	r.Use(middleware.Auth(di.Jwt, "admin"))
	r.Post("/", h.Create)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}
