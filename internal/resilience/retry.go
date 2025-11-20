package resilience

import (
	"math/rand"
	"time"
)

// RetryConfig holds the configuration for retry logic
type RetryConfig struct {
	MaxRetries  int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
	Jitter      bool
	ShouldRetry func(error) bool
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   10 * time.Second,
		Multiplier: 2.0,
		Jitter:     true,
		ShouldRetry: func(err error) bool {
			// By default, retry on any error
			return err != nil
		},
	}
}

// Retry executes the given function with retry logic
func Retry(fn func() (interface{}, error), config *RetryConfig) (interface{}, error) {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error

	for i := 0; i <= config.MaxRetries; i++ {
		result, err := fn()

		// If successful or shouldn't retry, return immediately
		if err == nil || !config.ShouldRetry(err) {
			return result, err
		}

		lastErr = err

		// Don't delay after the last attempt
		if i == config.MaxRetries {
			break
		}

		// Calculate delay
		delay := calculateDelay(config, i)

		// Wait for the delay
		time.Sleep(delay)
	}

	return nil, lastErr
}

// calculateDelay calculates the delay for the given retry attempt
func calculateDelay(config *RetryConfig, attempt int) time.Duration {
	delay := time.Duration(float64(config.BaseDelay) * pow(config.Multiplier, float64(attempt)))

	// Cap the delay
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}

	// Add jitter if enabled
	if config.Jitter {
		delay = addJitter(delay)
	}

	return delay
}

// pow calculates base^exp
func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

// addJitter adds random jitter to the delay
func addJitter(delay time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(delay) / 2))
	return delay + jitter
}
