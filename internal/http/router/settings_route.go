package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/http/handler"
)

func SettingsRoutes(r fiber.Router, h *handler.SettingsHandler) {
	r.Post("/", h.Create)
	r.Get("/", h.FindAll)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}
