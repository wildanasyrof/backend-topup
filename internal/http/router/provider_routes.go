package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/http/handler"
)

func ProviderRoutes(r fiber.Router, h *handler.ProviderHandler) {
	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}
