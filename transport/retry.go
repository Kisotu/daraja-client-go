// Package transport provides secure HTTP client configuration and middleware.
package transport

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// Default retry configuration
const (
	DefaultMaxRetries = 3
	DefaultBaseDelay  = 100 * time.Millisecond
	DefaultMaxDelay   = 5 * time.Second
)

// RetryableFunc determines if an error or response should trigger a retry.
type RetryableFunc func(resp *http.Response, err error) bool

// RetryTransport wraps an http.RoundTripper with retry logic.
type RetryTransport struct {
	Base       http.RoundTripper
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	Retryable  RetryableFunc
}

// RetryOption configures a RetryTransport.
type RetryOption func(*RetryTransport)

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(n int) RetryOption {
	return func(rt *RetryTransport) {
		rt.MaxRetries = n
	}
}

// WithBaseDelay sets the initial delay between retries.
func WithBaseDelay(d time.Duration) RetryOption {
	return func(rt *RetryTransport) {
		rt.BaseDelay = d
	}
}

// WithMaxDelay sets the maximum delay between retries.
func WithMaxDelay(d time.Duration) RetryOption {
	return func(rt *RetryTransport) {
		rt.MaxDelay = d
	}
}

// WithRetryableFn sets a custom retryable function.
func WithRetryableFn(fn RetryableFunc) RetryOption {
	return func(rt *RetryTransport) {
		rt.Retryable = fn
	}
}

// NewRetryTransport wraps a transport with retry logic.
func NewRetryTransport(base http.RoundTripper, opts ...RetryOption) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	rt := &RetryTransport{
		Base:       base,
		MaxRetries: DefaultMaxRetries,
		BaseDelay:  DefaultBaseDelay,
		MaxDelay:   DefaultMaxDelay,
		Retryable:  DefaultRetryableFunc,
	}
	for _, opt := range opts {
		opt(rt)
	}
	return rt
}

// RoundTrip implements http.RoundTripper with retry logic.
func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request body if it exists so we can retry
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("transport: failed to read request body: %w", err)
		}
		req.Body.Close()
	}

	var lastResp *http.Response
	var lastErr error

	for attempt := 0; attempt <= rt.MaxRetries; attempt++ {
		// Check context cancellation before each attempt
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
		}

		// Reset the body for retries
		if bodyBytes != nil {
			req.Body = io.NopCloser(&byteReader{data: bodyBytes})
		}

		// Clone request for each attempt to avoid side effects
		reqClone := req.Clone(req.Context())
		resp, err := rt.Base.RoundTrip(reqClone)

		// If successful and not retryable, return immediately
		if !rt.Retryable(resp, err) {
			return resp, err
		}

		// Store last result for potential return
		lastResp = resp
		lastErr = err

		// Don't retry on the last attempt
		if attempt == rt.MaxRetries {
			break
		}

		// Close response body if we're going to retry
		if resp != nil && resp.Body != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}

		// Calculate backoff with jitter
		delay := calculateBackoff(attempt, rt.BaseDelay, rt.MaxDelay)

		// Wait with context awareness
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(delay):
		}
	}

	// All retries exhausted, return last result
	if lastErr != nil {
		return nil, fmt.Errorf("transport: retries exhausted, last error: %w", lastErr)
	}
	return lastResp, nil
}

// DefaultRetryableFunc determines if a request should be retried.
func DefaultRetryableFunc(resp *http.Response, err error) bool {
	// Network errors are retryable
	if err != nil {
		return true
	}

	// Without a response, can't determine retry-ability
	if resp == nil {
		return false
	}

	// Retryable status codes: 5xx server errors and 429 rate limiting
	switch resp.StatusCode {
	case http.StatusTooManyRequests, // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	default:
		return false
	}
}

// calculateBackoff computes the delay with exponential backoff and jitter.
func calculateBackoff(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	// Exponential backoff: baseDelay * 2^attempt
	delay := baseDelay * (1 << attempt)

	// Add jitter: +/-25% randomization
	jitter := float64(delay) * 0.25
	offset := (rand.Float64() * 2 * jitter) - jitter
	delay = time.Duration(float64(delay) + offset)

	// Cap at maxDelay
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

// byteReader is a simple byte slice reader that can be reused.
type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
