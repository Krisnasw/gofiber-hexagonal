package config

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func NewRabbitMQConnection(cfg *viper.Viper, log *zap.Logger) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.GetString("RABBITMQ_URL"))
	if err != nil {
		log.Error("Failed to connect to RabbitMQ", zap.Error(err))
		return nil, err
	}

	log.Info("Successfully connected to RabbitMQ")
	return conn, nil
}

func SetupRabbitMQChannels(conn *amqp.Connection, log *zap.Logger) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open RabbitMQ channel", zap.Error(err))
		return nil, err
	}

	return channel, nil
}
