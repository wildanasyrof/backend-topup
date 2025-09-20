package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/http/handler"
)

func DepositRoutes(r fiber.Router, h *handler.DepositHandler) {
	r.Post("/", h.Create)
	r.Get("/", h.GetByDepositID)
	r.Get("/all", h.GetByUserID)
}
