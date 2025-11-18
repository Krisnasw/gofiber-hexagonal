package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	helper "app-hexagonal/internal/helper"
	"app-hexagonal/internal/usecase"
)

// UserRequest represents the user creation/update request structure with validation tags
type UserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Email string `json:"email" validate:"required,email"`
}

type UserHandler struct {
	uc       *usecase.UserUsecase
	logger   *zap.Logger
	validate *validator.Validate
}

func NewUserHandler(uc *usecase.UserUsecase, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		uc:       uc,
		logger:   logger,
		validate: validator.New(),
	}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Log the incoming request
	h.logger.Info("Getting user by ID",
		zap.String("user_id", id),
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
	)

	user, err := h.uc.GetUserByID(id)
	if err != nil {
		h.logger.Error("Failed to get user",
			zap.String("user_id", id),
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.Error(err),
		)
		return c.Status(fiber.StatusNotFound).JSON(helper.SuccessResponse(nil,
			fiber.StatusNotFound,
			"User Not Found"))
	}

	h.logger.Info("Successfully retrieved user",
		zap.String("user_id", id),
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		zap.String("user_name", user.Name),
	)

	return c.JSON(helper.SuccessResponse(user, fiber.StatusOK, "User retrieved successfully"))
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req UserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.Error(err),
		)
		return c.Status(fiber.StatusBadRequest).JSON(helper.SuccessResponse(nil,
			fiber.StatusBadRequest,
			"Invalid request body"))
	}

	// Log the incoming request
	h.logger.Info("Creating new user",
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		zap.String("user_name", req.Name),
		zap.String("user_email", req.Email),
	)

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Validation failed for user creation",
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.String("user_name", req.Name),
			zap.String("user_email", req.Email),
			zap.Error(err),
		)
		return c.Status(fiber.StatusBadRequest).JSON(helper.SuccessResponse(nil,
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error()))
	}

	// In a real implementation, you would create a domain user and save it
	// For now, we'll just return a success response
	h.logger.Info("User created successfully",
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		zap.String("user_name", req.Name),
		zap.String("user_email", req.Email),
	)

	return c.Status(fiber.StatusCreated).JSON(helper.SuccessResponse(nil,
		fiber.StatusCreated,
		"User created successfully"))
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/users/:id", h.GetUser)
	app.Post("/users", h.CreateUser)
}
