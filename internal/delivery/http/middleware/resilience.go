package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"app-hexagonal/internal/resilience"
)

// ResilienceMiddleware provides resilience patterns for HTTP handlers
type ResilienceMiddleware struct {
	handler *resilience.ResilienceHandler
	logger  *zap.Logger
}

// NewResilienceMiddleware creates a new resilience middleware
func NewResilienceMiddleware(handler *resilience.ResilienceHandler, logger *zap.Logger) *ResilienceMiddleware {
	return &ResilienceMiddleware{
		handler: handler,
		logger:  logger,
	}
}

// WithCircuitBreaker applies circuit breaker pattern to the handler
func (rm *ResilienceMiddleware) WithCircuitBreaker() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a key based on the request path
		key := c.Path()

		_, err := rm.handler.Execute(key, func() (interface{}, error) {
			// Process the request
			err := c.Next()
			return nil, err
		})

		if err != nil {
			rm.logger.Error("Request failed with resilience pattern",
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.Error(err),
			)

			// Check if it's a circuit breaker error
			if err.Error() == "circuit breaker is open" {
				return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
					"error":   "Service temporarily unavailable",
					"message": "Circuit breaker is open, service is temporarily unavailable",
				})
			}

			return err
		}

		return nil
	}
}

// WithRateLimiting applies rate limiting pattern to the handler
func (rm *ResilienceMiddleware) WithRateLimiting() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request is allowed
		availableTokens := rm.handler.GetRateLimiterAvailableTokens()
		if availableTokens <= 0 {
			rm.logger.Warn("Rate limit exceeded",
				zap.String("ip", c.IP()),
				zap.String("path", c.Path()),
			)

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
		}

		return c.Next()
	}
}

// WithTimeout applies timeout pattern to the handler
func (rm *ResilienceMiddleware) WithTimeout(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Apply timeout to the request
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		// Create a copy of the context with timeout
		c.SetUserContext(ctx)

		return c.Next()
	}
}
