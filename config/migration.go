package config

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

// RunMigrations runs database migrations
func RunMigrations(cfg *viper.Viper) error {
	// Get database configuration
	dbHost := cfg.GetString("DB_HOST")
	dbPort := cfg.GetInt("DB_PORT")
	dbUser := cfg.GetString("DB_USER")
	dbPassword := cfg.GetString("DB_PASSWORD")
	dbName := cfg.GetString("DB_NAME")
	dbType := cfg.GetString("DB_TYPE")

	var dbURL string
	switch dbType {
	case "mysql":
		dbURL = fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	case "postgres":
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Create migration instance
	m, err := migrate.New(
		"file://database/migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// RollbackMigrations rolls back the last migration
func RollbackMigrations(cfg *viper.Viper) error {
	// Get database configuration
	dbHost := cfg.GetString("DB_HOST")
	dbPort := cfg.GetInt("DB_PORT")
	dbUser := cfg.GetString("DB_USER")
	dbPassword := cfg.GetString("DB_PASSWORD")
	dbName := cfg.GetString("DB_NAME")
	dbType := cfg.GetString("DB_TYPE")

	var dbURL string
	switch dbType {
	case "mysql":
		dbURL = fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	case "postgres":
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Create migration instance
	m, err := migrate.New(
		"file://database/migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Rollback migrations
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Println("Migrations rolled back successfully")
	return nil
}
