package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func SetupRouter(app *fiber.App, di *di.DI, cfg *config.Config) {

	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })

	AuthRoutes(app.Group("/auth"), di.AuthHandler)

	me := app.Group("/me")
	me.Use(middleware.Auth(di.Jwt), "silver", "gold", "admin")

	UserRoutes(me, di.UserHandler)

	banner := app.Group("/banners")
	banner.Use(middleware.Auth(di.Jwt, "admin"))
	banner.Get("/", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
}
