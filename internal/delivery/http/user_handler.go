package http

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"app-hexagonalinternal/usecase"
)

type UserHandler struct {
	uc     *usecase.UserUsecase
	logger *zap.Logger
}

func NewUserHandler(uc *usecase.UserUsecase, logger *zap.Logger) *UserHandler {
	return &UserHandler{uc: uc, logger: logger}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.uc.GetUserByID(id)
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/users/:id", h.GetUser)
}
