package resilience

import (
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	maxTokens  int64
	refillRate time.Duration
	tokens     int64
	lastRefill time.Time
	mutex      sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxTokens int64, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		maxTokens:  maxTokens,
		refillRate: refillRate,
		tokens:     maxTokens,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed based on the rate limit
func (rl *RateLimiter) Allow() bool {
	return rl.AllowN(1)
}

// AllowN checks if n requests are allowed based on the rate limit
func (rl *RateLimiter) AllowN(n int64) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Refill tokens
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int64(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		rl.tokens = min(rl.maxTokens, rl.tokens+tokensToAdd)
		rl.lastRefill = now
	}

	// Check if we have enough tokens
	if rl.tokens >= n {
		rl.tokens -= n
		return true
	}

	return false
}

// Wait blocks until the rate limiter allows the request
func (rl *RateLimiter) Wait() error {
	return rl.WaitN(1)
}

// WaitN blocks until the rate limiter allows n requests
func (rl *RateLimiter) WaitN(n int64) error {
	for {
		if rl.AllowN(n) {
			return nil
		}
		// Sleep for a short duration before checking again
		time.Sleep(rl.refillRate / 10)
	}
}

// GetAvailableTokens returns the number of available tokens
func (rl *RateLimiter) GetAvailableTokens() int64 {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	return rl.tokens
}

// min returns the smaller of two int64 values
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
