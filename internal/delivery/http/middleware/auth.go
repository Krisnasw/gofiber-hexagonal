package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AuthMiddleware provides basic authentication middleware
func AuthMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// In a real implementation, you would check for a valid JWT token or session
		// For this example, we'll just check for an Authorization header
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			logger.Warn("Unauthorized access attempt",
				zap.String("path", c.Path()),
				zap.String("ip", c.IP()),
				zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Missing authorization header",
			})
		}

		// In a real implementation, you would validate the token here
		// For now, we'll just check if it starts with "Bearer "
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			logger.Warn("Invalid authorization header",
				zap.String("path", c.Path()),
				zap.String("ip", c.IP()),
				zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid authorization header format",
			})
		}

		// Log successful authentication
		logger.Info("Authentication successful",
			zap.String("path", c.Path()),
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		)

		// Add user info to context if needed
		// c.Locals("user_id", userID)

		return c.Next()
	}
}
