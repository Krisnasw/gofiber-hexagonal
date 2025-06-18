package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

func NewFiberConfig(cfg *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      cfg.GetString("APP_NAME"),
		Prefork:      cfg.GetBool("APP_PREFORK"),
		ErrorHandler: NewErrorHandler(),
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(healthcheck.New())

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"error":   true,
			"message": "Internal Server Error",
			"errors":  err.Error(),
		})
	}
}
