// Package circuitbreaker provides circuit breaker pattern implementation.
// The circuit breaker prevents cascade failures by temporarily rejecting
// requests when a service is experiencing high error rates.
package circuitbreaker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	// StateClosed - normal operation, requests pass through
	StateClosed State = iota
	// StateOpen - failing fast, requests are rejected immediately
	StateOpen
	// StateHalfOpen - testing if service recovered
	StateHalfOpen
)

// String returns the string representation of the state.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Errors returned by the circuit breaker.
var (
	ErrCircuitOpen = errors.New("circuit breaker: circuit is open")
)

// Config holds circuit breaker configuration.
type Config struct {
	// FailureThreshold is the number of consecutive failures before opening
	FailureThreshold uint32
	// SuccessThreshold is the number of consecutive successes in half-open before closing
	SuccessThreshold uint32
	// Timeout is the duration to wait before transitioning to half-open
	Timeout time.Duration
}

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	return Config{
		FailureThreshold: 5,
		SuccessThreshold: 3,
		Timeout:          30 * time.Second,
	}
}

// Metrics contains circuit breaker telemetry data.
type Metrics struct {
	State           State
	FailureCount    uint32
	SuccessCount    uint32
	ConsecutiveSucc uint32
	ConsecutiveFail uint32
	LastFailureTime time.Time
	LastSuccessTime time.Time
	TotalRequests   uint64
	TotalSuccesses  uint64
	TotalFailures   uint64
}

// Breaker implements the circuit breaker pattern.
type Breaker struct {
	config Config
	mu     sync.RWMutex
	state  State
	// Consecutive failure count in closed state
	consecutiveFailures uint32
	// Consecutive success count in half-open state
	consecutiveSuccesses uint32
	// Last time the circuit transitioned to open
	lastFailureTime time.Time
	// Metrics tracking
	totalRequests   uint64
	totalSuccesses  uint64
	totalFailures   uint64
	lastSuccessTime time.Time
}

// New creates a new circuit breaker with the given configuration.
func New(cfg Config) *Breaker {
	if cfg.FailureThreshold == 0 {
		cfg.FailureThreshold = DefaultConfig().FailureThreshold
	}
	if cfg.SuccessThreshold == 0 {
		cfg.SuccessThreshold = DefaultConfig().SuccessThreshold
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultConfig().Timeout
	}

	return &Breaker{
		config: cfg,
		state:  StateClosed,
	}
}

// Execute runs the given function if the circuit allows it.
// Returns ErrCircuitOpen if the circuit is open.
func (b *Breaker) Execute(fn func() error) error {
	if err := b.allow(); err != nil {
		return err
	}

	return b.doExecute(fn)
}

// ExecuteContext runs the given function with context if the circuit allows it.
// Returns ErrCircuitOpen if the circuit is open.
func (b *Breaker) ExecuteContext(ctx context.Context, fn func() error) error {
	if err := b.allow(); err != nil {
		return err
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return b.doExecute(fn)
}

// allow checks if the request should be allowed and handles state transitions.
func (b *Breaker) allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.totalRequests++

	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if time.Since(b.lastFailureTime) > b.config.Timeout {
			// Transition to half-open
			b.state = StateHalfOpen
			b.consecutiveSuccesses = 0
			return nil
		}
		return ErrCircuitOpen
	case StateHalfOpen:
		return nil
	}

	return nil
}

// doExecute runs the function and records the result.
func (b *Breaker) doExecute(fn func() error) error {
	err := fn()

	b.mu.Lock()
	defer b.mu.Unlock()

	if err == nil {
		b.onSuccess()
	} else {
		b.onFailure()
	}

	return err
}

// onSuccess handles successful execution.
func (b *Breaker) onSuccess() {
	b.totalSuccesses++
	b.lastSuccessTime = time.Now()

	switch b.state {
	case StateClosed:
		b.consecutiveFailures = 0
	case StateHalfOpen:
		b.consecutiveSuccesses++
		if b.consecutiveSuccesses >= b.config.SuccessThreshold {
			// Transition to closed
			b.state = StateClosed
			b.consecutiveFailures = 0
			b.consecutiveSuccesses = 0
		}
	}
}

// onFailure handles failed execution.
// updateState checks for and executes Open-to-HalfOpen transitions.
// This must be called while holding the lock.
func (b *Breaker) updateState() {
	if b.state == StateOpen && time.Since(b.lastFailureTime) > b.config.Timeout {
		b.state = StateHalfOpen
		b.consecutiveSuccesses = 0
	}
}

func (b *Breaker) onFailure() {
	b.totalFailures++
	b.lastFailureTime = time.Now()

	switch b.state {
	case StateClosed:
		b.consecutiveFailures++
		if b.consecutiveFailures >= b.config.FailureThreshold {
			// Transition to open
			b.state = StateOpen
		}
	case StateHalfOpen:
		// Transition back to open
		b.state = StateOpen
		b.consecutiveSuccesses = 0
	}
}

// State returns the current state of the circuit breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.updateState()
	return b.state
}

// Metrics returns current metrics for the circuit breaker.
func (b *Breaker) Metrics() Metrics {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.updateState()

	return Metrics{
		State:           b.state,
		FailureCount:    b.consecutiveFailures,
		SuccessCount:    b.consecutiveSuccesses,
		ConsecutiveSucc: b.consecutiveSuccesses,
		ConsecutiveFail: b.consecutiveFailures,
		LastFailureTime: b.lastFailureTime,
		LastSuccessTime: b.lastSuccessTime,
		TotalRequests:   b.totalRequests,
		TotalSuccesses:  b.totalSuccesses,
		TotalFailures:   b.totalFailures,
	}
}

// Reset forces the circuit breaker to closed state.
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.state = StateClosed
	b.consecutiveFailures = 0
	b.consecutiveSuccesses = 0
}

// IsCircuitOpenError reports whether err is a circuit open error.
func IsCircuitOpenError(err error) bool {
	return errors.Is(err, ErrCircuitOpen)
}

func (b *Breaker) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return fmt.Sprintf("Breaker{state:%s, failures:%d, successes:%d}",
		b.state, b.consecutiveFailures, b.consecutiveSuccesses)
}
