package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app-hexagonal/config"
	"app-hexagonal/internal/repository"
	"app-hexagonal/internal/resilience"
	"app-hexagonal/internal/usecase"

	"go.uber.org/zap"
)

func runWorker() {
	fmt.Println("Starting worker service...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, err := config.NewLogger(cfg)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	db, err := config.NewGormDB(cfg, log)
	if err != nil {
		panic(err)
	}

	// Initialize RabbitMQ connection
	_, err = config.NewRabbitMQConnection(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}

	// Create repository and usecase
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)

	// Set up resilience handler
	resilienceConfig := resilience.DefaultResilienceConfig()
	resilienceHandler := resilience.NewResilienceHandler(resilienceConfig)

	// Start worker loop
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Worker service started. Press Ctrl+C to stop.")

	for {
		select {
		case <-sigChan:
			fmt.Println("Shutting down worker...")
			return
		case <-ticker.C:
			// Process worker tasks with resilience
			_, err := resilienceHandler.Execute("worker_task", func() (interface{}, error) {
				log.Info("Processing worker task")

				// Example task - replace with actual worker logic
				// For now, we'll just get a user by ID as an example
				user, err := userUsecase.GetUserByID("example-id")
				if err != nil {
					log.Info("No user found with example ID (expected)", zap.Error(err))
					return nil, nil // Not an error for the worker
				}

				log.Info("Processed user", zap.String("user_id", user.ID))
				return user, nil
			})

			if err != nil {
				log.Error("Worker task failed", zap.Error(err))
			}
		}
	}
}
