package db

import (
	"fmt"

	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	logger "github.com/wildanasyrof/backend-topup/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config, logger logger.Logger) *gorm.DB {
	dsn := buildPostgresDSN(cfg) // Asumsikan cfg.Db adalah struct yang menampung field DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database error, " + err.Error())
	}
	if err := db.AutoMigrate(&entity.User{}, &entity.Menu{}, &entity.Settings{}, &entity.PaymentMethod{}, &entity.Banner{}, &entity.Deposit{}, &entity.Provider{}, &entity.Category{}, &entity.UserLevel{}, &entity.Product{}, &entity.Price{}, &entity.Order{}); err != nil {
		logger.Fatal("auto-migrate failed: " + err.Error())
	}
	return db
}

func buildPostgresDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.Db.Host,
		cfg.Db.Username,
		cfg.Db.Password,
		cfg.Db.Database,
		cfg.Db.Port,
	)
}
