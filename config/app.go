package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"app-hexagonalinternal/delivery/http"
	"app-hexagonalinternal/delivery/http/route"
	"app-hexagonalinternal/repository"
	"app-hexagonalinternal/usecase"
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
