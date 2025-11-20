package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"app-hexagonal/internal/domain"
	helper "app-hexagonal/internal/helper"
	"app-hexagonal/internal/resilience"
	"app-hexagonal/internal/usecase"
)

// LoginRequest represents the login request structure
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RefreshRequest represents the refresh token request structure
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
	logger      *zap.Logger
	validate    *validator.Validate
	resilience  *resilience.ResilienceHandler
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase *usecase.AuthUsecase, logger *zap.Logger, resilience *resilience.ResilienceHandler) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		logger:      logger,
		validate:    validator.New(),
		resilience:  resilience,
	}
}

// Login handles user login requests
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse login request body",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.Error(err),
		)
		return c.Status(fiber.StatusBadRequest).JSON(helper.ErrorResponse(nil,
			fiber.StatusBadRequest,
			"Invalid request body"))
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Validation failed for login request",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return c.Status(fiber.StatusBadRequest).JSON(helper.ErrorResponse(nil,
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error()))
	}

	// Log the login attempt
	h.logger.Info("Login attempt",
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		zap.String("email", req.Email),
	)

	// Execute with resilience patterns
	result, err := h.resilience.Execute("auth_login", func() (interface{}, error) {
		credentials := &domain.Credentials{
			Email:    req.Email,
			Password: req.Password,
		}

		return h.authUsecase.Login(credentials)
	})

	if err != nil {
		h.logger.Error("Login failed",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse(nil,
			fiber.StatusUnauthorized,
			"Invalid credentials"))
	}

	tokenResponse := result.(*domain.TokenResponse)

	h.logger.Info("Login successful",
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		zap.String("email", req.Email),
	)

	return c.JSON(helper.SuccessResponseWithMetadata(tokenResponse,
		fiber.StatusOK,
		"Login successful",
		helper.Metadata{}))
}

// Refresh handles token refresh requests
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse refresh request body",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.Error(err),
		)
		return c.Status(fiber.StatusBadRequest).JSON(helper.ErrorResponse(nil,
			fiber.StatusBadRequest,
			"Invalid request body"))
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Validation failed for refresh request",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.Error(err),
		)
		return c.Status(fiber.StatusBadRequest).JSON(helper.ErrorResponse(nil,
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error()))
	}

	// Log the refresh attempt
	h.logger.Info("Token refresh attempt",
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
	)

	// Execute with resilience patterns
	result, err := h.resilience.Execute("auth_refresh", func() (interface{}, error) {
		return h.authUsecase.RefreshToken(req.RefreshToken)
	})

	if err != nil {
		h.logger.Error("Token refresh failed",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.Error(err),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse(nil,
			fiber.StatusUnauthorized,
			"Invalid refresh token"))
	}

	tokenResponse := result.(*domain.TokenResponse)

	h.logger.Info("Token refresh successful",
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
	)

	return c.JSON(helper.SuccessResponseWithMetadata(tokenResponse,
		fiber.StatusOK,
		"Token refreshed successfully",
		helper.Metadata{}))
}

// Logout handles user logout requests
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get the access token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		h.logger.Error("Missing Authorization header",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse(nil,
			fiber.StatusUnauthorized,
			"Missing Authorization header"))
	}

	// Extract the token (Bearer token)
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		h.logger.Error("Invalid Authorization header format",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.String("auth_header", authHeader),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse(nil,
			fiber.StatusUnauthorized,
			"Invalid Authorization header format"))
	}

	// Log the logout attempt
	h.logger.Info("Logout attempt",
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
	)

	// Execute with resilience patterns
	_, err := h.resilience.Execute("auth_logout", func() (interface{}, error) {
		return nil, h.authUsecase.Logout(token)
	})

	if err != nil {
		h.logger.Error("Logout failed",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.Error(err),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse(nil,
			fiber.StatusUnauthorized,
			"Invalid token"))
	}

	h.logger.Info("Logout successful",
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
	)

	return c.JSON(helper.SuccessResponse(nil,
		fiber.StatusOK,
		"Logout successful"))
}

// RegisterRoutes registers the authentication routes
func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	// Public routes
	app.Post("/auth/login", h.Login)
	app.Post("/auth/refresh", h.Refresh)

	// Protected routes (would typically require auth middleware)
	app.Post("/auth/logout", h.Logout)
}
