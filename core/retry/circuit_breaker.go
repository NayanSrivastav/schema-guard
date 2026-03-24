package retry

import (
	"errors"
	"sync"
	"time"
)

// ErrCircuitOpen is thrown when execution is blocked due to threshold failures
var ErrCircuitOpen = errors.New("circuit breaker is open, LLM calls suspended")

// CircuitBreaker limits uncontrolled failures and costs across the app
type CircuitBreaker struct {
	mu        sync.Mutex
	failures  int
	threshold int
	timeout   time.Duration
	openedAt  time.Time
}

// NewCircuitBreaker creates a circuit breaker that trips after maxFailures and resets after timeout.
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: maxFailures,
		timeout:   resetTimeout,
	}
}

// CanExecute checks if the circuit is closed or if the timeout has elapsed
func (cb *CircuitBreaker) CanExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.failures >= cb.threshold {
		// If timeout has elapsed, transition to half-open (allow 1 execution to test)
		if time.Since(cb.openedAt) > cb.timeout {
			// For simplicity in MVP, we just reset it fully instead of full half-open concurrency mgmt
			cb.failures = 0
			cb.openedAt = time.Time{}
			return true
		}
		return false
	}
	return true
}

// RecordSuccess zeroes out error counts
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.openedAt = time.Time{}
}

// RecordFailure bumps error count and opens circuit if passing bounds
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.failures++
	if cb.failures >= cb.threshold {
		// Only set timestamp if this is the failure that tripped the circuit, or reset the window?
		// Restarting the window on continuous failures prevents recovery during cascades.
		cb.openedAt = time.Now() 
	}
}
