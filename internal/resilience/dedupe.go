package resilience

import (
	"sync"
	"time"
)

// Dedupe implements request deduplication to prevent processing duplicate requests
type Dedupe struct {
	cache   map[string]*dedupeEntry
	mutex   sync.Mutex
	timeout time.Duration
}

// dedupeEntry represents a cached request
type dedupeEntry struct {
	result    interface{}
	err       error
	timestamp time.Time
}

// NewDedupe creates a new deduplication handler
func NewDedupe(timeout time.Duration) *Dedupe {
	return &Dedupe{
		cache:   make(map[string]*dedupeEntry),
		timeout: timeout,
	}
}

// Execute executes the given function with deduplication based on the key
func (d *Dedupe) Execute(key string, fn func() (interface{}, error)) (interface{}, error) {
	d.mutex.Lock()

	// Check if we have a cached result
	if entry, exists := d.cache[key]; exists {
		// Check if the cached result is still valid
		if time.Since(entry.timestamp) < d.timeout {
			d.mutex.Unlock()
			return entry.result, entry.err
		}
		// Remove expired entry
		delete(d.cache, key)
	}

	d.mutex.Unlock()

	// Execute the function
	result, err := fn()

	// Cache the result
	d.mutex.Lock()
	d.cache[key] = &dedupeEntry{
		result:    result,
		err:       err,
		timestamp: time.Now(),
	}
	d.mutex.Unlock()

	// Clean up expired entries periodically
	go d.cleanup()

	return result, err
}

// cleanup removes expired entries from the cache
func (d *Dedupe) cleanup() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	now := time.Now()
	for key, entry := range d.cache {
		if now.Sub(entry.timestamp) >= d.timeout {
			delete(d.cache, key)
		}
	}
}

// GetCacheSize returns the number of cached entries
func (d *Dedupe) GetCacheSize() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return len(d.cache)
}
