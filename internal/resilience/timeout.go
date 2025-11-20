package resilience

import (
	"context"
	"time"
)

// TimeoutConfig holds the configuration for timeout logic
type TimeoutConfig struct {
	Timeout time.Duration
}

// DefaultTimeoutConfig returns a default timeout configuration
func DefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		Timeout: 30 * time.Second,
	}
}

// Timeout executes the given function with a timeout
func Timeout(fn func() (interface{}, error), config *TimeoutConfig) (interface{}, error) {
	if config == nil {
		config = DefaultTimeoutConfig()
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Create channels for result and error
	resultChan := make(chan interface{}, 1)
	errorChan := make(chan error, 1)

	// Execute the function in a goroutine
	go func() {
		result, err := fn()
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()

	// Wait for either the function to complete or the timeout
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
