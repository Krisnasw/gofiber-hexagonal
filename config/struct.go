package config

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Config represents the structured application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Server   ServerConfig   `mapstructure:"server"`
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Environment string `mapstructure:"env" validate:"required,oneof=development production staging"`
	Name        string `mapstructure:"name" validate:"required"`
	Version     string `mapstructure:"version"`
	Debug       bool   `mapstructure:"debug"`
	GRPCPort    int    `mapstructure:"grpc_port" validate:"required,min=1,max=65535"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name" validate:"required"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string        `mapstructure:"host" validate:"required"`
	Port     int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	TTL      time.Duration `mapstructure:"ttl"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
