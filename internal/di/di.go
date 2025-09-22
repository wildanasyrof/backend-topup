package di

import (
	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/db"
	"github.com/wildanasyrof/backend-topup/internal/http/handler"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/jwt"
	logger "github.com/wildanasyrof/backend-topup/pkg/logger"
	"github.com/wildanasyrof/backend-topup/pkg/storage"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
	"gorm.io/gorm"
)

type DI struct {
	Logger                logger.Logger
	DB                    *gorm.DB
	Jwt                   jwt.JWTService
	Storage               storage.LocalStorage
	AuthHandler           *handler.AuthHandler
	UserHandler           *handler.UserHandler
	MenuHandler           *handler.MenuHandler
	SettingsHandler       *handler.SettingsHandler
	PaymentMethodsHandler *handler.PaymentMethodsHandler
	BannerHandler         *handler.BannerHandler
	DepositHanlder        *handler.DepositHandler
	ProviderHandler       *handler.ProviderHandler
}

func InitDI(cfg *config.Config) *DI {
	logger := logger.NewZerologLogger(cfg.Server.Env)
	DB := db.Connect(cfg, logger)
	validator := validator.NewValidator()
	jwt := jwt.NewJWTService(cfg)
	storage := storage.NewLocalStorage(cfg)

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

	paymentMethodRepo := repository.NewPaymentMethodsRepository(DB)
	paymentMethodService := service.NewPaymentMethodsService(paymentMethodRepo)
	paymentMethodsHandler := handler.NewPaymentMethodsHandler(paymentMethodService, validator, storage)

	bannerRepo := repository.NewBannerRepository(DB)
	bannerService := service.NewBannerService(bannerRepo)
	bannerHandler := handler.NewBannerHandler(bannerService, storage)

	depositRepo := repository.NewDepositRepository(DB)
	depositService := service.NewDepositService(depositRepo)
	depositHandler := handler.NewDepositHandler(depositService, validator, logger)

	providerRepo := repository.NewProviderRepository(DB)
	providerService := service.NewProviderService(providerRepo)
	providerHandler := handler.NewProviderHandler(providerService, validator)

	return &DI{
		Logger:                logger,
		DB:                    DB,
		Jwt:                   jwt,
		Storage:               storage,
		AuthHandler:           authHandler,
		UserHandler:           userHandler,
		MenuHandler:           menuHandler,
		SettingsHandler:       settingsHandler,
		PaymentMethodsHandler: paymentMethodsHandler,
		BannerHandler:         bannerHandler,
		DepositHanlder:        depositHandler,
		ProviderHandler:       providerHandler,
	}
}
