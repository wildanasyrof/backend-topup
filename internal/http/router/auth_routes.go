package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/http/handler"
)

func AuthRoutes(r fiber.Router, a *handler.AuthHandler) {
	r.Post("/register", a.Register)
	r.Post("/login", a.Login)
	r.Get("/google/login", a.GoogleLogin)
	r.Get("/google/callback", a.GoogleCallback)
}
