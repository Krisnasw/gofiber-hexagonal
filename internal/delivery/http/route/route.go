package route

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	"app-hexagonal/internal/delivery/http"
	"app-hexagonal/internal/delivery/http/middleware"
)

type RouteConfig struct {
	App         *fiber.App
	UserHandler *http.UserHandler
	Logger      *zap.Logger
}

func (c *RouteConfig) Setup() {
	// Apply global middleware
	c.App.Use(middleware.CORSMiddleware())
	c.App.Use(middleware.LoggingMiddleware(c.Logger))

	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World")
	})

	c.App.Get("/swagger/*", swagger.HandlerDefault) // default swagger

	// Health check endpoints
	c.App.Get("/health", c.HealthCheck)
	c.App.Get("/ready", c.ReadinessCheck)
}

func (c *RouteConfig) SetupAuthRoute() {
	// Protected routes would go here
	// For example:
	// c.App.Use(middleware.AuthMiddleware(c.Logger))
	// c.UserHandler.RegisterRoutes(c.App)
}

// HealthCheck returns the health status of the application
func (c *RouteConfig) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// ReadinessCheck returns the readiness status of the application
// In a real implementation, this would check database connections, cache, etc.
func (c *RouteConfig) ReadinessCheck(ctx *fiber.Ctx) error {
	// TODO: Add actual checks for database, cache, etc.
	// For now, we'll just return ready
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ready",
		"timestamp": time.Now().Unix(),
		"checks": map[string]string{
			"database": "unknown", // Would be "ok" or "error" in real implementation
			"cache":    "unknown", // Would be "ok" or "error" in real implementation
		},
	})
}
