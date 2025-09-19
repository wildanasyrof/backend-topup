package db

import (
	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	logger "github.com/wildanasyrof/backend-topup/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config, logger logger.Logger) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.Db.DbUrl), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database error, " + err.Error())
	}
	if err := db.AutoMigrate(&entity.User{}, &entity.Menu{}, &entity.Settings{}); err != nil {
		logger.Fatal("auto-migrate failed: " + err.Error())
	}
	return db
}
