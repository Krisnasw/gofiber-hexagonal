package gorm

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"
)

// TransactionManager handles database transactions with safety features
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// TransactionFunc represents a function that operates within a transaction
type TransactionFunc func(tx *gorm.DB) error

// WithTransaction executes a function within a database transaction
// It handles retries and provides safety mechanisms for race conditions
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	return tm.WithTransactionOptions(ctx, fn, &TransactionOptions{})
}

// TransactionOptions configures transaction behavior
type TransactionOptions struct {
	MaxRetries     int
	RetryDelay     time.Duration
	IsolationLevel string
	Timeout        time.Duration
}

// DefaultTransactionOptions returns default transaction options
func DefaultTransactionOptions() *TransactionOptions {
	return &TransactionOptions{
		MaxRetries:     3,
		RetryDelay:     100 * time.Millisecond,
		IsolationLevel: "READ_COMMITTED",
		Timeout:        30 * time.Second,
	}
}

// WithTransactionOptions executes a function within a database transaction with custom options
func (tm *TransactionManager) WithTransactionOptions(ctx context.Context, fn TransactionFunc, opts *TransactionOptions) error {
	if opts == nil {
		opts = DefaultTransactionOptions()
	}

	// Apply timeout to context if specified
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	var lastErr error

	// Retry mechanism for handling transient errors and race conditions
	for attempt := 0; attempt <= opts.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff with jitter
			delay := time.Duration(float64(opts.RetryDelay) * pow(2, float64(attempt-1)))
			jitter := time.Duration(float64(delay) * 0.1 * float64(attempt%3)) // Simple jitter
			time.Sleep(delay + jitter)
		}

		// Create a new context for each attempt to avoid timeout carryover
		attemptCtx := ctx
		if opts.Timeout > 0 && attempt > 0 {
			var cancel context.CancelFunc
			attemptCtx, cancel = context.WithTimeout(context.Background(), opts.Timeout)
			defer cancel()
		}

		err := tm.executeTransaction(attemptCtx, fn, opts)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on non-retryable errors
		if !isRetryableError(err) {
			return err
		}

		// Log retry attempt
		if attempt < opts.MaxRetries {
			// In a real implementation, you would use a logger
		}
	}

	return lastErr
}

// executeTransaction performs the actual transaction execution
func (tm *TransactionManager) executeTransaction(ctx context.Context, fn TransactionFunc, opts *TransactionOptions) error {
	// Begin transaction with isolation level
	tx := tm.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Ensure rollback in case of panic or early return
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-panic
		}
	}()

	// Set isolation level if specified
	if opts.IsolationLevel != "" {
		if err := tx.Exec("SET TRANSACTION ISOLATION LEVEL " + opts.IsolationLevel).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Execute the transaction function
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// ConcurrentTransactionManager handles concurrent transactions with additional safety
type ConcurrentTransactionManager struct {
	TransactionManager
	lock sync.Mutex
}

// NewConcurrentTransactionManager creates a new concurrent transaction manager
func NewConcurrentTransactionManager(db *gorm.DB) *ConcurrentTransactionManager {
	return &ConcurrentTransactionManager{
		TransactionManager: *NewTransactionManager(db),
	}
}

// WithTransaction executes a function within a database transaction with concurrency protection
func (ctm *ConcurrentTransactionManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	return ctm.WithTransactionOptions(ctx, fn, &TransactionOptions{})
}

// WithTransactionOptions executes a function within a database transaction with concurrency protection and custom options
func (ctm *ConcurrentTransactionManager) WithTransactionOptions(ctx context.Context, fn TransactionFunc, opts *TransactionOptions) error {
	// Acquire lock to prevent race conditions in critical sections
	ctm.lock.Lock()
	defer ctm.lock.Unlock()

	return ctm.TransactionManager.WithTransactionOptions(ctx, fn, opts)
}

// isRetryableError determines if an error should trigger a retry
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Common retryable errors (MySQL/PostgreSQL)
	retryableErrors := []string{
		"Deadlock found when trying to get lock",
		"try restarting transaction",
		"database is locked",
		"Lock wait timeout exceeded",
		"Too many connections",
		"connection refused",
	}

	errStr := err.Error()
	for _, retryable := range retryableErrors {
		if contains(errStr, retryable) {
			return true
		}
	}

	return false
}

// pow calculates the power of a number
func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
