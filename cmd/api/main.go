package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Aplikasi berjalan. Tekan CTRL+C untuk berhenti.")
	<-quit // Menunggu sinyal SIGINT/SIGTERM untuk keluar
	log.Println("Menerima sinyal shutdown...")
}
