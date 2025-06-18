package config

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"app-hexagonalpkg/mysql"
	"app-hexagonalpkg/postgres"
)

func NewGormDB(cfg *viper.Viper, log *zap.Logger) (*gorm.DB, error) {
	db, err := mysql.Connect(
		cfg.GetString("DATABASE_HOST"),
		cfg.GetInt("DATABASE_PORT"),
		cfg.GetString("DATABASE_USER"),
		cfg.GetString("DATABASE_PASSWORD"),
		cfg.GetString("DATABASE_NAME"),
		mysql.SetPrintLog(
			cfg.GetBool("DATABASE_LOG_ENABLED"),
			logger.LogLevel(cfg.GetInt("DATABASE_LOG_LEVEL")),
			time.Duration(cfg.GetInt("DATABASE_LOG_THRESHOLD"))*time.Millisecond,
		),
	)
	if err != nil {
		log.Fatal("Failed to initialize mysql database connection", zap.Error(err))
	}

	_, err = postgres.Connect(
		cfg.GetString("PG_HOST"),
		cfg.GetInt("PG_PORT"),
		cfg.GetString("PG_USER"),
		cfg.GetString("PG_PASSWORD"),
		cfg.GetString("PG_NAME"),
	)
	if err != nil {
		log.Fatal("Failed to initialize postgres database connection", zap.Error(err))
	}

	return db, nil
}
