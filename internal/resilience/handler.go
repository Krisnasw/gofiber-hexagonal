package resilience

import (
	"time"
)

// ResilienceHandler combines all resilience patterns into a single handler
type ResilienceHandler struct {
	circuitBreaker *CircuitBreaker
	retryConfig    *RetryConfig
	timeoutConfig  *TimeoutConfig
	bulkhead       *Bulkhead
	rateLimiter    *RateLimiter
	dedupe         *Dedupe
}

// ResilienceConfig holds the configuration for all resilience patterns
type ResilienceConfig struct {
	// Circuit Breaker
	CircuitFailureThreshold int
	CircuitSuccessThreshold int
	CircuitTimeout          time.Duration

	// Retry
	RetryMaxRetries int
	RetryBaseDelay  time.Duration
	RetryMaxDelay   time.Duration
	RetryMultiplier float64
	RetryJitter     bool

	// Timeout
	TimeoutDuration time.Duration

	// Bulkhead
	BulkheadMaxConcurrent int

	// Rate Limiter
	RateLimitMaxTokens  int64
	RateLimitRefillRate time.Duration

	// Dedupe
	DedupeTimeout time.Duration
}

// DefaultResilienceConfig returns a default resilience configuration
func DefaultResilienceConfig() *ResilienceConfig {
	return &ResilienceConfig{
		CircuitFailureThreshold: 5,
		CircuitSuccessThreshold: 3,
		CircuitTimeout:          60 * time.Second,

		RetryMaxRetries: 3,
		RetryBaseDelay:  100 * time.Millisecond,
		RetryMaxDelay:   10 * time.Second,
		RetryMultiplier: 2.0,
		RetryJitter:     true,

		TimeoutDuration: 30 * time.Second,

		BulkheadMaxConcurrent: 10,

		RateLimitMaxTokens:  100,
		RateLimitRefillRate: time.Second,

		DedupeTimeout: 5 * time.Second,
	}
}

// NewResilienceHandler creates a new resilience handler with the given configuration
func NewResilienceHandler(config *ResilienceConfig) *ResilienceHandler {
	if config == nil {
		config = DefaultResilienceConfig()
	}

	return &ResilienceHandler{
		circuitBreaker: NewCircuitBreaker(
			config.CircuitFailureThreshold,
			config.CircuitSuccessThreshold,
			config.CircuitTimeout,
		),
		retryConfig: &RetryConfig{
			MaxRetries:  config.RetryMaxRetries,
			BaseDelay:   config.RetryBaseDelay,
			MaxDelay:    config.RetryMaxDelay,
			Multiplier:  config.RetryMultiplier,
			Jitter:      config.RetryJitter,
			ShouldRetry: func(err error) bool { return err != nil },
		},
		timeoutConfig: &TimeoutConfig{
			Timeout: config.TimeoutDuration,
		},
		bulkhead: NewBulkhead(config.BulkheadMaxConcurrent),
		rateLimiter: NewRateLimiter(
			config.RateLimitMaxTokens,
			config.RateLimitRefillRate,
		),
		dedupe: NewDedupe(config.DedupeTimeout),
	}
}

// Execute executes the given function with all resilience patterns applied
func (rh *ResilienceHandler) Execute(key string, fn func() (interface{}, error)) (interface{}, error) {
	// Apply deduplication first
	return rh.dedupe.Execute(key, func() (interface{}, error) {
		// Apply rate limiting
		if !rh.rateLimiter.Allow() {
			// Fallback when rate limited
			return Fallback(
				func() (interface{}, error) {
					return nil, &RateLimitError{Msg: "rate limit exceeded"}
				},
				func(err error) (interface{}, error) {
					return nil, err
				},
			)
		}

		// Apply bulkhead
		return rh.bulkhead.Execute(func() (interface{}, error) {
			// Apply circuit breaker
			return rh.circuitBreaker.Execute(func() (interface{}, error) {
				// Apply retry with timeout
				return Retry(func() (interface{}, error) {
					return Timeout(fn, rh.timeoutConfig)
				}, rh.retryConfig)
			})
		})
	})
}

// RateLimitError represents a rate limit error
type RateLimitError struct {
	Msg string
}

func (e *RateLimitError) Error() string {
	return e.Msg
}

// GetCircuitBreakerState returns the current state of the circuit breaker
func (rh *ResilienceHandler) GetCircuitBreakerState() CircuitBreakerState {
	return rh.circuitBreaker.state
}

// GetBulkheadCurrent returns the current number of concurrent requests in the bulkhead
func (rh *ResilienceHandler) GetBulkheadCurrent() int {
	return rh.bulkhead.GetCurrent()
}

// GetRateLimiterAvailableTokens returns the number of available tokens in the rate limiter
func (rh *ResilienceHandler) GetRateLimiterAvailableTokens() int64 {
	return rh.rateLimiter.GetAvailableTokens()
}

// GetDedupeCacheSize returns the number of cached entries in the dedupe handler
func (rh *ResilienceHandler) GetDedupeCacheSize() int {
	return rh.dedupe.GetCacheSize()
}
