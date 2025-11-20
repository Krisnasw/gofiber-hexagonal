package resilience

import (
	"errors"
	"sync"
)

// Bulkhead implements the bulkhead pattern to limit concurrent requests
type Bulkhead struct {
	maxConcurrent int
	semaphore     chan struct{}
	mutex         sync.Mutex
	current       int
}

// NewBulkhead creates a new bulkhead with the specified maximum concurrent requests
func NewBulkhead(maxConcurrent int) *Bulkhead {
	return &Bulkhead{
		maxConcurrent: maxConcurrent,
		semaphore:     make(chan struct{}, maxConcurrent),
	}
}

// Execute executes the given function with bulkhead protection
func (b *Bulkhead) Execute(fn func() (interface{}, error)) (interface{}, error) {
	// Try to acquire a slot
	select {
	case b.semaphore <- struct{}{}:
		// Successfully acquired a slot
		defer func() {
			<-b.semaphore // Release the slot
		}()

		// Update current count
		b.mutex.Lock()
		b.current++
		current := b.current
		b.mutex.Unlock()

		defer func() {
			b.mutex.Lock()
			b.current--
			b.mutex.Unlock()
		}()

		// Check if we're still within limits
		if current > b.maxConcurrent {
			return nil, errors.New("bulkhead limit exceeded")
		}

		// Execute the function
		return fn()
	default:
		// No slots available
		return nil, errors.New("bulkhead capacity exceeded")
	}
}

// GetCurrent returns the current number of concurrent requests
func (b *Bulkhead) GetCurrent() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.current
}

// GetMaxConcurrent returns the maximum number of concurrent requests allowed
func (b *Bulkhead) GetMaxConcurrent() int {
	return b.maxConcurrent
}
