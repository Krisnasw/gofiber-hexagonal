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
	v.SetDefault("APP_NAME", "Go Hexagonal Boilerplate")
	v.SetDefault("APP_VERSION", "1.0.0")
	v.SetDefault("APP_DEBUG", false)
	v.SetDefault("APP_PORT", 4001)
	v.SetDefault("APP_DEFAULT_LANG", "en")
	v.SetDefault("APP_TIMEZONE", "+07:00")
	v.SetDefault("APP_PREFORK", false)

	v.SetDefault("DATABASE_PORT", 3306)
	v.SetDefault("DATABASE_HOST", "localhost")
	v.SetDefault("DATABASE_USER", "root")
	v.SetDefault("DATABASE_PASSWORD", "")
	v.SetDefault("DATABASE_NAME", "superfood")
	v.SetDefault("DATABASE_SSL_MODE", "disable")
	v.SetDefault("DATABASE_LOG_ENABLED", true)
	v.SetDefault("DATABASE_LOG_LEVEL", 3)
	v.SetDefault("DATABASE_LOG_THRESHOLD", 200)

	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("REDIS_TTL", 3600*time.Second)

	v.SetDefault("SERVER_PORT", 4001)
	v.SetDefault("SERVER_READ_TIMEOUT", 5*time.Second)
	v.SetDefault("SERVER_WRITE_TIMEOUT", 10*time.Second)
	v.SetDefault("SERVER_IDLE_TIMEOUT", 60*time.Second)

	v.SetDefault("PG_PORT", 5432)
	v.SetDefault("PG_HOST", "localhost")
	v.SetDefault("PG_USER", "postgres")
	v.SetDefault("PG_PASSWORD", "")
	v.SetDefault("PG_NAME", "superfood")
	v.SetDefault("PG_LOG_ENABLED", true)
	v.SetDefault("PG_LOG_LEVEL", 3)
	v.SetDefault("PG_LOG_THRESHOLD", 200)

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

// LoadStructuredConfig loads configuration into a structured format
func LoadStructuredConfig() (*Config, error) {
	v, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	// Bind environment variables to struct fields
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Map individual values to structured config
	config.App.Environment = v.GetString("APP_ENV")
	config.App.Name = v.GetString("APP_NAME")
	config.App.Version = v.GetString("APP_VERSION")
	config.App.Debug = v.GetBool("APP_DEBUG")

	config.Database.Host = v.GetString("DATABASE_HOST")
	config.Database.Port = v.GetInt("DATABASE_PORT")
	config.Database.User = v.GetString("DATABASE_USER")
	config.Database.Password = v.GetString("DATABASE_PASSWORD")
	config.Database.Name = v.GetString("DATABASE_NAME")
	config.Database.SSLMode = v.GetString("DATABASE_SSL_MODE")

	config.Redis.Host = v.GetString("REDIS_HOST")
	config.Redis.Port = v.GetInt("REDIS_PORT")
	config.Redis.Password = v.GetString("REDIS_PASSWORD")
	config.Redis.DB = v.GetInt("REDIS_DB")
	config.Redis.TTL = v.GetDuration("REDIS_TTL")

	config.Server.Port = v.GetInt("SERVER_PORT")
	config.Server.ReadTimeout = v.GetDuration("SERVER_READ_TIMEOUT")
	config.Server.WriteTimeout = v.GetDuration("SERVER_WRITE_TIMEOUT")
	config.Server.IdleTimeout = v.GetDuration("SERVER_IDLE_TIMEOUT")

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}
