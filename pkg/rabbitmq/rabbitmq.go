package rabbitmq

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// Config holds RabbitMQ configuration
type Config struct {
	URL               string
	ReconnectInterval time.Duration
	ReconnectAttempts int
	HeartbeatInterval time.Duration
	ConnectionTimeout time.Duration
}

// Option is a function that configures RabbitMQ
type Option func(*Config)

// Connection represents a RabbitMQ connection wrapper
type Connection struct {
	conn     *amqp.Connection
	config   *Config
	channels map[string]*amqp.Channel
}

// New creates a new RabbitMQ connection
func New(url string, options ...Option) (*Connection, error) {
	config := &Config{
		URL:               url,
		ReconnectInterval: 5 * time.Second,
		ReconnectAttempts: 10,
		HeartbeatInterval: 10 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	}

	for _, option := range options {
		option(config)
	}

	conn, err := amqp.DialConfig(config.URL, amqp.Config{
		Heartbeat: config.HeartbeatInterval,
		Locale:    "en_US",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return &Connection{
		conn:     conn,
		config:   config,
		channels: make(map[string]*amqp.Channel),
	}, nil
}

// WithReconnectInterval sets the reconnect interval
func WithReconnectInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.ReconnectInterval = interval
	}
}

// WithReconnectAttempts sets the number of reconnect attempts
func WithReconnectAttempts(attempts int) Option {
	return func(c *Config) {
		c.ReconnectAttempts = attempts
	}
}

// WithHeartbeatInterval sets the heartbeat interval
func WithHeartbeatInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.HeartbeatInterval = interval
	}
}

// WithConnectionTimeout sets the connection timeout
func WithConnectionTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ConnectionTimeout = timeout
	}
}

// GetConnection returns the underlying AMQP connection
func (c *Connection) GetConnection() *amqp.Connection {
	return c.conn
}

// CreateChannel creates a new channel with the given name
func (c *Connection) CreateChannel(name string) (*amqp.Channel, error) {
	channel, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	c.channels[name] = channel
	return channel, nil
}

// GetChannel returns a channel by name
func (c *Connection) GetChannel(name string) (*amqp.Channel, bool) {
	channel, exists := c.channels[name]
	return channel, exists
}

// Close closes the RabbitMQ connection and all channels
func (c *Connection) Close() error {
	// Close all channels first
	for _, channel := range c.channels {
		if channel != nil {
			channel.Close()
		}
	}

	// Close the connection
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// IsClosed checks if the connection is closed
func (c *Connection) IsClosed() bool {
	if c.conn == nil {
		return true
	}
	return c.conn.IsClosed()
}

// Reconnect attempts to reconnect to RabbitMQ
func (c *Connection) Reconnect() error {
	if !c.IsClosed() {
		c.Close()
	}

	conn, err := amqp.DialConfig(c.config.URL, amqp.Config{
		Heartbeat: c.config.HeartbeatInterval,
		Locale:    "en_US",
	})
	if err != nil {
		return fmt.Errorf("failed to reconnect to RabbitMQ: %w", err)
	}

	c.conn = conn
	return nil
}
