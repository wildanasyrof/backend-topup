package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/router"
)

func main() {
	app := fiber.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error load config:%v:", err)
	}

	di := di.InitDI(cfg)
	router.SetupRouter(app, di, cfg)

	app.Listen(":" + cfg.Server.Port)
}
