package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// LoggingMiddleware logs incoming requests
func LoggingMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Log the incoming request
		logger.Info("Incoming request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", string(c.Request().Header.UserAgent())),
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		)

		// Process request
		err := c.Next()

		// Log request details
		duration := time.Since(start)
		statusCode := c.Response().StatusCode()

		logger.Info("Request completed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		)

		return err
	}
}
