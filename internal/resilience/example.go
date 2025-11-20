package resilience

import (
	"errors"
	"fmt"
	"time"
)

// Example usage of resilience patterns
func Example() {
	// Create a resilience handler with default configuration
	config := DefaultResilienceConfig()
	handler := NewResilienceHandler(config)

	// Example function that sometimes fails
	failingFunction := func() (interface{}, error) {
		// Simulate a random failure
		if time.Now().Unix()%3 == 0 {
			return nil, errors.New("simulated failure")
		}
		return "success", nil
	}

	// Execute with resilience patterns
	result, err := handler.Execute("example-key", failingFunction)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %v\n", result)
	}

	// Show resilience handler stats
	fmt.Printf("Circuit Breaker State: %v\n", handler.GetCircuitBreakerState())
	fmt.Printf("Bulkhead Current: %d\n", handler.GetBulkheadCurrent())
	fmt.Printf("Rate Limiter Tokens: %d\n", handler.GetRateLimiterAvailableTokens())
	fmt.Printf("Dedupe Cache Size: %d\n", handler.GetDedupeCacheSize())
}
