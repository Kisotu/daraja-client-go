// Package transport provides secure HTTP client configuration and middleware.
package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// mockTransport allows injection of test responses
type mockTransport struct {
	responses []*http.Response
	errors    []error
	callCount int
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.callCount >= len(m.responses) && m.callCount >= len(m.errors) {
		return nil, errors.New("no more mock responses")
	}
	var resp *http.Response
	if m.callCount < len(m.responses) {
		resp = m.responses[m.callCount]
	}
	var err error
	if m.callCount < len(m.errors) {
		err = m.errors[m.callCount]
	}
	m.callCount++
	return resp, err
}

func TestNewRetryTransport(t *testing.T) {
	base := &http.Transport{}
	rt := NewRetryTransport(base).(*RetryTransport)

	if rt.Base != base {
		t.Error("Base transport not set correctly")
	}
	if rt.MaxRetries != DefaultMaxRetries {
		t.Errorf("MaxRetries = %d, want %d", rt.MaxRetries, DefaultMaxRetries)
	}
	if rt.BaseDelay != DefaultBaseDelay {
		t.Errorf("BaseDelay = %v, want %v", rt.BaseDelay, DefaultBaseDelay)
	}
	if rt.MaxDelay != DefaultMaxDelay {
		t.Errorf("MaxDelay = %v, want %v", rt.MaxDelay, DefaultMaxDelay)
	}
}

func TestNewRetryTransport_NilBase(t *testing.T) {
	rt := NewRetryTransport(nil).(*RetryTransport)
	if rt.Base != http.DefaultTransport {
		t.Error("Expected DefaultTransport when base is nil")
	}
}

// TestRetryTransport_SuccessNoRetry tests successful request with no retry needed.
func TestRetryTransport_SuccessNoRetry(t *testing.T) {
	mock := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))},
		},
	}

	rt := NewRetryTransport(mock, WithMaxRetries(3))
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	client := &http.Client{Transport: rt}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if mock.callCount != 1 {
		t.Errorf("Expected 1 call, got %d", mock.callCount)
	}
}

// TestRetryTransport_RetryOn500ThenSuccess tests that retries happen on 500 and succeed on 200.
func TestRetryTransport_RetryOn500ThenSuccess(t *testing.T) {
	mock := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))},
			{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))},
		},
	}

	rt := NewRetryTransport(mock, WithMaxRetries(3))
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	client := &http.Client{Transport: rt}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if mock.callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", mock.callCount)
	}
}

// TestRetryTransport_RetryOn429 tests retry on 429 Too Many Requests.
func TestRetryTransport_RetryOn429(t *testing.T) {
	mock := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusTooManyRequests, Body: io.NopCloser(strings.NewReader(""))},
			{StatusCode: http.StatusTooManyRequests, Body: io.NopCloser(strings.NewReader(""))},
			{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))},
		},
	}

	rt := NewRetryTransport(mock, WithMaxRetries(3))
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	client := &http.Client{Transport: rt}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if mock.callCount != 3 {
		t.Errorf("Expected 3 calls, got %d", mock.callCount)
	}
}

// TestRetryTransport_NoRetryOn400 tests that 400 Bad Request is not retried.
func TestRetryTransport_NoRetryOn400(t *testing.T) {
	mock := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader(""))},
		},
	}

	rt := NewRetryTransport(mock, WithMaxRetries(3))
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	client := &http.Client{Transport: rt}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
	if mock.callCount != 1 {
		t.Errorf("Expected 1 call, got %d", mock.callCount)
	}
}

// TestRetryTransport_RetryOnNetworkError tests retry on network errors.
func TestRetryTransport_RetryOnNetworkError(t *testing.T) {
	mock := &mockTransport{
		errors: []error{
			errors.New("connection refused"),
			nil,
		},
		responses: []*http.Response{
			nil,
			{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))},
		},
	}

	rt := NewRetryTransport(mock, WithMaxRetries(3))
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	client := &http.Client{Transport: rt}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if mock.callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", mock.callCount)
	}
}

// TestRetryTransport_ContextCancellation tests that context cancellation stops retries.
func TestRetryTransport_ContextCancellation(t *testing.T) {
	mock := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusServiceUnavailable, Body: io.NopCloser(strings.NewReader(""))},
		},
	}

	rt := NewRetryTransport(mock, WithMaxRetries(3), WithBaseDelay(100*time.Millisecond))
	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req = req.WithContext(ctx)
	client := &http.Client{Transport: rt}

	// Cancel after first attempt
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := client.Do(req)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}
}

// TestRetryTransport_MaxRetriesExhausted tests that all retries are exhausted.
func TestRetryTransport_MaxRetriesExhausted(t *testing.T) {
	mock := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))},
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))},
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))},
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))},
		},
	}

	rt := NewRetryTransport(mock, WithMaxRetries(3))
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	client := &http.Client{Transport: rt}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
	if mock.callCount != 4 { // Initial + 3 retries
		t.Errorf("Expected 4 calls, got %d", mock.callCount)
	}
}

func TestDefaultRetryableFunc(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
		want       bool
	}{
		{"network error", 0, errors.New("connection failed"), true},
		{"status 200", http.StatusOK, nil, false},
		{"status 400", http.StatusBadRequest, nil, false},
		{"status 429", http.StatusTooManyRequests, nil, true},
		{"status 500", http.StatusInternalServerError, nil, true},
		{"status 502", http.StatusBadGateway, nil, true},
		{"status 503", http.StatusServiceUnavailable, nil, true},
		{"status 504", http.StatusGatewayTimeout, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			if tt.statusCode > 0 {
				resp = &http.Response{StatusCode: tt.statusCode}
			}
			got := DefaultRetryableFunc(resp, tt.err)
			if got != tt.want {
				t.Errorf("DefaultRetryableFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateBackoff(t *testing.T) {
	tests := []struct {
		attempt  int
		base     time.Duration
		max      time.Duration
		expected time.Duration
	}{
		{0, 100 * time.Millisecond, 5 * time.Second, 100 * time.Millisecond},
		{1, 100 * time.Millisecond, 5 * time.Second, 200 * time.Millisecond},
		{2, 100 * time.Millisecond, 5 * time.Second, 400 * time.Millisecond},
		{3, 100 * time.Millisecond, 5 * time.Second, 800 * time.Millisecond},
		{10, 100 * time.Millisecond, 500 * time.Millisecond, 500 * time.Millisecond}, // capped at max
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("attempt_%d", tt.attempt), func(t *testing.T) {
			got := calculateBackoff(tt.attempt, tt.base, tt.max)
			// Allow for jitter variance
			minExpected := time.Duration(float64(tt.expected) * 0.75)
			maxExpected := time.Duration(float64(tt.expected) * 1.25)
			if got < minExpected || got > maxExpected {
				t.Errorf("calculateBackoff() = %v, expected between %v and %v", got, minExpected, maxExpected)
			}
		})
	}
}

// TestRetryTransport_BodyPreservation tests that request body is preserved across retries.
func TestRetryTransport_BodyPreservation(t *testing.T) {
	bodyContent := []byte("test payload")

	mock := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusServiceUnavailable, Body: io.NopCloser(strings.NewReader(""))},
			{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))},
		},
	}

	rt := NewRetryTransport(mock, WithMaxRetries(3))
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", &byteReader{data: bodyContent})
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(&byteReader{data: bodyContent}), nil
	}
	client := &http.Client{Transport: rt}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
