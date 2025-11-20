package resilience

import (
	"errors"
	"sync"
	"time"
)

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	state            CircuitBreakerState
	failureThreshold int
	successThreshold int
	timeout          time.Duration
	lastFailure      time.Time
	mutex            sync.Mutex

	failures  int
	successes int
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            Closed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	}
}

// Execute executes the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	cb.mutex.Lock()

	// Check if we can execute based on current state
	switch cb.state {
	case Open:
		// Check if timeout has elapsed
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = HalfOpen
		} else {
			cb.mutex.Unlock()
			return nil, errors.New("circuit breaker is open")
		}
	case HalfOpen:
		// Allow only one request through
	default:
		// Closed state - proceed normally
	}

	cb.mutex.Unlock()

	// Execute the function
	result, err := fn()

	// Update circuit breaker state based on result
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}

	return result, err
}

func (cb *CircuitBreaker) onFailure() {
	cb.failures++
	cb.lastFailure = time.Now()

	if cb.state == HalfOpen || cb.failures >= cb.failureThreshold {
		cb.state = Open
	}
}

func (cb *CircuitBreaker) onSuccess() {
	cb.successes++

	switch cb.state {
	case HalfOpen:
		if cb.successes >= cb.successThreshold {
			cb.state = Closed
			cb.failures = 0
			cb.successes = 0
		}
	case Closed:
		cb.failures = 0
	}
}
