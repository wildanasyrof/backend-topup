package di

import (
	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/db"
	"github.com/wildanasyrof/backend-topup/internal/http/handler"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/jwt"
	logger "github.com/wildanasyrof/backend-topup/pkg/logger"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
	"gorm.io/gorm"
)

type DI struct {
	Logger          logger.Logger
	DB              *gorm.DB
	Jwt             jwt.JWTService
	AuthHandler     *handler.AuthHandler
	UserHandler     *handler.UserHandler
	MenuHandler     *handler.MenuHandler
	SettingsHandler *handler.SettingsHandler
}

func InitDI(cfg *config.Config) *DI {
	logger := logger.NewZerologLogger(cfg.Server.Env)
	DB := db.Connect(cfg, logger)
	validator := validator.NewValidator()
	jwt := jwt.NewJWTService(cfg)

	userRepo := repository.NewUserRepository(DB)
	authService := service.NewAuthService(userRepo, jwt)
	authHandler := handler.NewAuthHandler(authService, validator)

	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, validator)

	menuRepo := repository.NewMenuRepository(DB)
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService, validator)

	settingsRepo := repository.NewSettingsRepository(DB)
	settingsService := service.NewSettingsService(settingsRepo)
	settingsHandler := handler.NewSettingsHandler(settingsService, validator)

	return &DI{
		Logger:          logger,
		DB:              DB,
		Jwt:             jwt,
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		MenuHandler:     menuHandler,
		SettingsHandler: settingsHandler,
	}
}
