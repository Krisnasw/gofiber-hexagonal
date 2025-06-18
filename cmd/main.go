package main

import (
	"fmt"

	"go.uber.org/zap"

	config "app-hexagonalconfig"
)

// @title SuperFood Panel Service API
// @version 1.0
// @description This is a server for SuperFood Panel Service API
// @termsOfService http://swagger.io/terms/
// @contact.name SuperFood
// @contact.email dev.superfood@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:4001
// @BasePath /
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, err := config.NewLogger(cfg)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	db, err := config.NewGormDB(cfg, log)
	if err != nil {
		panic(err)
	}

	validate := config.NewValidator(cfg)

	// Initialize RabbitMQ connection
	rabbitMq, err := config.NewRabbitMQConnection(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}

	app := config.NewFiberConfig(cfg)

	config.Boostrap(&config.BoostrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   cfg,
		RabbitMQ: rabbitMq,
	})

	appPort := cfg.GetInt("APP_PORT")
	// Start server
	log.Info("Starting server",
		zap.String("env", cfg.GetString("APP_ENV")),
		zap.Int("port", cfg.GetInt("APP_PORT")),
	)

	err = app.Listen(fmt.Sprintf(":%d", appPort))
	if err != nil {
		log.Fatal("Failed to start services", zap.Error(err))
	}
}
