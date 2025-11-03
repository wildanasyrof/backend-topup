package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/middleware"
)

func SetupRouter(app *fiber.App, di *di.DI, cfg *config.Config) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", // <-- URL Frontend Anda
		AllowCredentials: true,                    // <-- WAJIB untuk cookie
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(requestid.New(requestid.Config{
		Header:     "X-Request-ID",
		ContextKey: "requestid",
	}))
	app.Use(middleware.LoggerMiddleware(di.Logger))
	app.Use(middleware.TimeoutMiddleware(time.Duration(cfg.Server.RequestTimeOut) * time.Second))

	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })

	app.Static("/uploads", cfg.Server.UploadDir)

	AuthRoutes(app.Group("/auth"), di.AuthHandler) // <-- Rute /auth

	menu := app.Group("/menus")
	MenuRoutes(menu, di.MenuHandler, di)

	app.Get("/menu", di.MenuHandler.GetAll)
	app.Get("/categories", di.CategoryHandler.GetAll)
	app.Get("/categories/:slug", di.CategoryHandler.GetBySlug)
	me := app.Group("/me")
	me.Use(middleware.Auth(di.Jwt, "admin", "user"))
	UserRoutes(me, di.UserHandler)

	// --- TAMBAHKAN INI ---
	// Grup /sessions untuk manajemen sesi (remote logout)
	sessions := app.Group("/sessions")
	sessions.Use(middleware.Auth(di.Jwt, "admin", "user")) // Wajib login
	SessionRoutes(sessions, di)
	// ---------------------

	settings := app.Group("/settings")
	settings.Use(middleware.Auth(di.Jwt, "admin"))
	SettingsRoutes(settings, di.SettingsHandler)

	paymentMethods := app.Group("/payment-methods")
	PaymentMethodsRoutes(paymentMethods, di.PaymentMethodsHandler, di)

	banner := app.Group("/banners")
	BannerRoutes(banner, di)

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

	order := app.Group("/orders")
	OrderRoutes(order, di)
}
