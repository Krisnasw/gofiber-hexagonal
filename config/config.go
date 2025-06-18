package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", 4001)
	v.SetDefault("APP_DEFAULT_LANG", "en")
	v.SetDefault("APP_TIMEZONE", "+07:00")
	v.SetDefault("APP_PREFORK", false)

	v.SetDefault("DATABASE_PORT", 3306)
	v.SetDefault("DATABASE_HOST", "localhost")
	v.SetDefault("DATABASE_USER", "root")
	v.SetDefault("DATABASE_PASSWORD", "")
	v.SetDefault("DATABASE_NAME", "superfood")
	v.SetDefault("DATABASE_LOG_ENABLED", true)
	v.SetDefault("DATABASE_LOG_LEVEL", 3)
	v.SetDefault("DATABASE_LOG_THRESHOLD", 200)

	v.SetDefault("PG_PORT", 5432)
	v.SetDefault("PG_HOST", "localhost")
	v.SetDefault("PG_USER", "postgres")
	v.SetDefault("PG_PASSWORD", "")
	v.SetDefault("PG_NAME", "superfood")
	v.SetDefault("PG_LOG_ENABLED", true)
	v.SetDefault("PG_LOG_LEVEL", 3)
	v.SetDefault("PG_LOG_THRESHOLD", 200)

	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_TTL", 3600*time.Second)
	v.SetDefault("REDIS_ENABLE_SEARCH_PLACEHOLDER", true)

	v.SetDefault("RATE_LIMIT_ENABLED", false)
	v.SetDefault("RATE_LIMIT_WINDOWMS", 2*time.Second)
	v.SetDefault("RATE_LIMIT_MAX", 500)

	v.SetDefault("HTTP_TIMEOUT", 30000*time.Millisecond)
	v.SetDefault("HTTP_MAX_REDIRECTS", 5)

	v.SetDefault("BATCH_PROCESSING_SIZE", 500)

	// Set up to read from .env file
	v.SetConfigType("env")
	v.SetConfigFile(".env")

	// Read the main .env file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}

	// Automatically load environment variables
	v.AutomaticEnv()

	return v, nil
}
