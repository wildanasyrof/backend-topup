package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func SetupRouter(app *fiber.App, di *di.DI, cfg *config.Config) {
	app.Use(middleware.LoggerMiddleware(di.Logger))

	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })

	app.Static("/uploads", cfg.Server.UploadDir)

	AuthRoutes(app.Group("/auth"), di.AuthHandler)

	me := app.Group("/me")
	me.Use(middleware.Auth(di.Jwt, "admin", "user"))
	UserRoutes(me, di.UserHandler)

	menu := app.Group("/menus")
	MenuRoutes(menu, di.MenuHandler, di)

	settings := app.Group("/settings")
	settings.Use(middleware.Auth(di.Jwt, "admin"))
	SettingsRoutes(settings, di.SettingsHandler)

	paymentMethods := app.Group("/payment-methods")
	PaymentMethodsRoutes(paymentMethods, di.PaymentMethodsHandler, di)

	banner := app.Group("/banners")
	banner.Use(middleware.Auth(di.Jwt, "admin"))
	BannerRoutes(banner, di.BannerHandler)

	deposit := app.Group("/deposits")
	deposit.Use(middleware.Auth(di.Jwt, "admin", "user"))
	DepositRoutes(deposit, di.DepositHanlder)

	provider := app.Group("/providers")
	provider.Use(middleware.Auth(di.Jwt, "admin"))
	ProviderRoutes(provider, di.ProviderHandler)

	category := app.Group("/categories")
	CategoryRotues(category, di)

	product := app.Group("/products")
	ProductRouter(product, di)

	price := app.Group("/prices")
	PriceRoutes(price, di)
}
