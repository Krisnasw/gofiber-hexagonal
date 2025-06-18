package config

import (
	"app-hexagonal/internal/delivery/http"
	"app-hexagonal/internal/delivery/http/route"
	"app-hexagonal/internal/repository"
	"app-hexagonal/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BoostrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *zap.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	RabbitMQ *amqp.Connection
}

func Boostrap(config *BoostrapConfig) {
	// Repository
	userRepository := repository.NewUserRepository(config.DB)

	// UseCase
	userUseCase := usecase.NewUserUsecase(userRepository)

	// Handler
	userHandler := http.NewUserHandler(userUseCase, config.Log)

	routeConfig := route.RouteConfig{
		App:         config.App,
		UserHandler: userHandler,
	}
	routeConfig.Setup()
}
