package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/di"
)

// SessionRoutes mengatur rute untuk manajemen sesi
func SessionRoutes(r fiber.Router, di *di.DI) {
	// Semua rute di sini memerlukan autentikasi (sudah diatur di routes.go)
	r.Get("/", di.SessionHandler.GetMySessions)
	r.Delete("/:id", di.SessionHandler.RevokeSession)
}
