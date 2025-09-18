package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/di"
)

func SetupRouter(app *fiber.App, di *di.DI, cfg *config.Config) {

	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })

	AuthRoutes(app.Group("/auth"), di.AuthHandler)
}
