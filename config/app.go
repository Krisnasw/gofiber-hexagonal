package config

import (
	"app-hexagonal/internal/delivery/http"
	"app-hexagonal/internal/delivery/http/route"
	"app-hexagonal/internal/repository"
	"app-hexagonal/internal/resilience"
	"app-hexagonal/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BoostrapConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Log         *zap.Logger
	Validate    *validator.Validate
	Config      *viper.Viper
	RabbitMQ    *amqp.Connection
	UserUsecase usecase.UserUsecaseInterface
}

func Boostrap(config *BoostrapConfig) {
	// Repository
	userRepository := repository.NewUserRepository(config.DB)
	// No need for authRepository since we're using JWT

	// UseCase
	var userUseCase usecase.UserUsecaseInterface
	if config.UserUsecase != nil {
		userUseCase = config.UserUsecase
	} else {
		userUseCase = usecase.NewUserUsecase(userRepository)
	}

	// Create auth usecase (no need for auth repository with JWT)
	authUseCase := usecase.NewAuthUsecase(userRepository)

	// Resilience Handler
	resilienceConfig := resilience.DefaultResilienceConfig()
	resilienceHandler := resilience.NewResilienceHandler(resilienceConfig)

	// Handler
	userHandler := http.NewUserHandler(userUseCase, config.Log, resilienceHandler)
	authHandler := http.NewAuthHandler(authUseCase, config.Log, resilienceHandler)

	routeConfig := route.RouteConfig{
		App:         config.App,
		UserHandler: userHandler,
		AuthHandler: authHandler,
		Logger:      config.Log,
	}
	routeConfig.Setup()
}
