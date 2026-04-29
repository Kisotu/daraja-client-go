// Package circuitbreaker provides circuit breaker pattern implementation.
package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	b := New(DefaultConfig())

	if b.State() != StateClosed {
		t.Errorf("Initial state = %v, want %v", b.State(), StateClosed)
	}
}

func TestNew_DefaultConfig(t *testing.T) {
	cfg := Config{
		FailureThreshold: 0,
		SuccessThreshold: 0,
		Timeout:          0,
	}
	b := New(cfg)

	// Should use defaults
	if b.config.FailureThreshold != 5 {
		t.Errorf("FailureThreshold = %d, want 5", b.config.FailureThreshold)
	}
	if b.config.SuccessThreshold != 3 {
		t.Errorf("SuccessThreshold = %d, want 3", b.config.SuccessThreshold)
	}
	if b.config.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", b.config.Timeout)
	}
}

func TestBreaker_NormalOperation(t *testing.T) {
	b := New(DefaultConfig())

	err := b.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if b.State() != StateClosed {
		t.Errorf("State = %v, want %v", b.State(), StateClosed)
	}
}

func TestBreaker_OpensAfterThreshold(t *testing.T) {
	cfg := Config{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		Timeout:          30 * time.Second,
	}
	b := New(cfg)

	// Trigger 3 failures
	for i := 0; i < 3; i++ {
		if err := b.Execute(func() error {
			return errors.New("failure")
		}); err == nil {
			t.Errorf("expected failure while opening circuit")
		}
	}

	if b.State() != StateOpen {
		t.Errorf("State = %v, want %v", b.State(), StateOpen)
	}
}

func TestBreaker_ReturnsCircuitOpenError(t *testing.T) {
	cfg := Config{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          30 * time.Second,
	}
	b := New(cfg)

	// Trigger failure to open circuit
	if err := b.Execute(func() error {
		return errors.New("failure")
	}); err == nil {
		t.Fatalf("expected failure while opening circuit")
	}

	if b.State() != StateOpen {
		t.Fatalf("Expected open state")
	}

	// Next request should fail fast
	err := b.Execute(func() error {
		return nil
	})

	if !errors.Is(err, ErrCircuitOpen) {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}
}

func TestBreaker_TransitionsToHalfOpen(t *testing.T) {
	cfg := Config{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          50 * time.Millisecond,
	}
	b := New(cfg)

	// Open the circuit
	if err := b.Execute(func() error {
		return errors.New("failure")
	}); err == nil {
		t.Fatalf("expected failure while opening circuit")
	}

	if b.State() != StateOpen {
		t.Fatalf("Expected open state")
	}

	// Wait for timeout (use longer sleep to ensure transition)
	time.Sleep(100 * time.Millisecond)

	// Should now be half-open
	if b.State() != StateHalfOpen {
		t.Errorf("State = %v, want %v", b.State(), StateHalfOpen)
	}
}

func TestBreaker_ClosesAfterSuccessThreshold(t *testing.T) {
	cfg := Config{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		Timeout:          50 * time.Millisecond,
	}
	b := New(cfg)

	// Open the circuit
	if err := b.Execute(func() error {
		return errors.New("failure")
	}); err == nil {
		t.Fatalf("expected failure while opening circuit")
	}

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)

	// Should be half-open now
	if b.State() != StateHalfOpen {
		t.Fatalf("Expected half-open state, got %v", b.State())
	}

	// Success 1
	if err := b.Execute(func() error {
		return nil
	}); err != nil {
		t.Fatalf("expected success in half-open state: %v", err)
	}

	if b.State() != StateHalfOpen {
		t.Errorf("State = %v, want %v after first success", b.State(), StateHalfOpen)
	}

	// Success 2 - should close
	if err := b.Execute(func() error {
		return nil
	}); err != nil {
		t.Fatalf("expected success closing circuit: %v", err)
	}

	if b.State() != StateClosed {
		t.Errorf("State = %v, want %v after second success", b.State(), StateClosed)
	}
}

func TestBreaker_ReturnsToOpenOnHalfOpenFailure(t *testing.T) {
	cfg := Config{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		Timeout:          50 * time.Millisecond,
	}
	b := New(cfg)

	// Open the circuit
	if err := b.Execute(func() error {
		return errors.New("failure")
	}); err == nil {
		t.Fatalf("expected failure while opening circuit")
	}

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)

	// Should be half-open now
	if b.State() != StateHalfOpen {
		t.Fatalf("Expected half-open state")
	}

	// Failure in half-open returns to open
	if err := b.Execute(func() error {
		return errors.New("failure")
	}); err == nil {
		t.Fatalf("expected failure in half-open state")
	}

	if b.State() != StateOpen {
		t.Errorf("State = %v, want %v", b.State(), StateOpen)
	}
}

func TestBreaker_Reset(t *testing.T) {
	cfg := Config{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          30 * time.Second,
	}
	b := New(cfg)

	// Open the circuit
	if err := b.Execute(func() error {
		return errors.New("failure")
	}); err == nil {
		t.Fatalf("expected failure while opening circuit")
	}

	if b.State() != StateOpen {
		t.Fatalf("Expected open state")
	}

	// Reset
	b.Reset()

	if b.State() != StateClosed {
		t.Errorf("State = %v, want %v after reset", b.State(), StateClosed)
	}
}

func TestBreaker_ExecuteContext(t *testing.T) {
	b := New(DefaultConfig())

	// Success path
	err := b.ExecuteContext(context.Background(), func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = b.ExecuteContext(ctx, func() error {
		return nil
	})

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestBreaker_Metrics(t *testing.T) {
	cfg := Config{
		FailureThreshold: 5,
		SuccessThreshold: 3,
		Timeout:          30 * time.Second,
	}
	b := New(cfg)

	// Execute some requests
	if err := b.Execute(func() error { return nil }); err != nil {
		t.Fatalf("expected success: %v", err)
	}
	if err := b.Execute(func() error { return errors.New("fail") }); err == nil {
		t.Fatalf("expected failure")
	}
	if err := b.Execute(func() error { return nil }); err != nil {
		t.Fatalf("expected success: %v", err)
	}

	metrics := b.Metrics()

	if metrics.State != StateClosed {
		t.Errorf("Metrics.State = %v, want %v", metrics.State, StateClosed)
	}
	if metrics.TotalRequests != 3 {
		t.Errorf("TotalRequests = %d, want 3", metrics.TotalRequests)
	}
	if metrics.TotalSuccesses != 2 {
		t.Errorf("TotalSuccesses = %d, want 2", metrics.TotalSuccesses)
	}
	if metrics.TotalFailures != 1 {
		t.Errorf("TotalFailures = %d, want 1", metrics.TotalFailures)
	}
}

func TestBreaker_ConcurrentAccess(t *testing.T) {
	cfg := Config{
		FailureThreshold: 100,
		SuccessThreshold: 50,
		Timeout:          30 * time.Second,
	}
	b := New(cfg)

	done := make(chan bool, 100)

	// Concurrent requests
	for i := 0; i < 100; i++ {
		go func(fail bool) {
			if fail {
				if err := b.Execute(func() error { return errors.New("fail") }); err == nil {
					t.Errorf("expected failure")
				}
			} else {
				if err := b.Execute(func() error { return nil }); err != nil {
					t.Errorf("expected success: %v", err)
				}
			}
			done <- true
		}(i%2 == 0)
	}

	// Wait for all
	for i := 0; i < 100; i++ {
		<-done
	}

	metrics := b.Metrics()
	if metrics.TotalRequests != 100 {
		t.Errorf("TotalRequests = %d, want 100", metrics.TotalRequests)
	}
}

func TestIsCircuitOpenError(t *testing.T) {
	if !IsCircuitOpenError(ErrCircuitOpen) {
		t.Error("IsCircuitOpenError(ErrCircuitOpen) should be true")
	}
	if IsCircuitOpenError(errors.New("other error")) {
		t.Error("IsCircuitOpenError(other error) should be false")
	}
}

func TestState_String(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{State(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("State(%d).String() = %v, want %v", tt.state, got, tt.want)
		}
	}
}
