// Package callback provides helpers for secure callback handling.
package callback

import (
	"bytes"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// STKPushResult contains the full STK Push callback payload.
type STKPushResult struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `json:"MerchantRequestID"`
			CheckoutRequestID string `json:"CheckoutRequestID"`
			ResultCode        int    `json:"ResultCode"`
			ResultDesc        string `json:"ResultDesc"`
			CallbackMetadata  struct {
				Item []struct {
					Name  string      `json:"Name"`
					Value interface{} `json:"Value"`
				} `json:"Item"`
			} `json:"CallbackMetadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

// Verifier validates incoming Daraja callbacks.
type Verifier struct {
	ipAllowlist []*net.IPNet
	ipCheckFn   func(net.IP) bool
}

// VerifierOption configures a Verifier.
type VerifierOption func(*Verifier)

// WithIPAllowlist restricts callback acceptance to the given CIDR ranges.
func WithIPAllowlist(cidrs []string) VerifierOption {
	return func(v *Verifier) {
		for _, cidr := range cidrs {
			_, n, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			v.ipAllowlist = append(v.ipAllowlist, n)
		}
	}
}

// WithIPCheckFn sets a custom function to validate callback source IPs.
func WithIPCheckFn(fn func(net.IP) bool) VerifierOption {
	return func(v *Verifier) {
		v.ipCheckFn = fn
	}
}

// NewVerifier returns a new callback Verifier.
func NewVerifier(opts ...VerifierOption) *Verifier {
	v := &Verifier{}
	for _, o := range opts {
		o(v)
	}
	return v
}

// VerifySTKPush reads and validates an STK Push callback from r.
func (v *Verifier) VerifySTKPush(r *http.Request) (*STKPushResult, error) {
	if !v.ipAllowed(r) {
		return nil, fmt.Errorf("callback: source IP not allowed")
	}

	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("callback: expected POST, got %s", r.Method)
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		return nil, fmt.Errorf("callback: expected application/json, got %s", ct)
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("callback: failed to read body: %w", err)
	}

	var result STKPushResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("callback: failed to decode callback payload: %w", err)
	}

	if result.Body.StkCallback.CheckoutRequestID == "" {
		return nil, fmt.Errorf("callback: missing CheckoutRequestID")
	}

	return &result, nil
}

func (v *Verifier) ipAllowed(r *http.Request) bool {
	if len(v.ipAllowlist) == 0 && v.ipCheckFn == nil {
		return true
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	if v.ipCheckFn != nil && v.ipCheckFn(ip) {
		return true
	}

	if len(v.ipAllowlist) == 0 {
		return false
	}

	for _, n := range v.ipAllowlist {
		if n.Contains(ip) {
			return true
		}
	}

	return false
}

// ConstantTimeEqual performs a constant-time comparison of two byte slices.
func ConstantTimeEqual(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// ValidateCallbackPayload validates an STK Push callback payload against
// expected values.
func ValidateCallbackPayload(result *STKPushResult, expectedCheckoutID, expectedMerchantID string) error {
	if result.Body.StkCallback.CheckoutRequestID != expectedCheckoutID {
		return fmt.Errorf("callback: CheckoutRequestID mismatch: expected %q, got %q",
			expectedCheckoutID, result.Body.StkCallback.CheckoutRequestID)
	}
	if expectedMerchantID != "" && result.Body.StkCallback.MerchantRequestID != expectedMerchantID {
		return fmt.Errorf("callback: MerchantRequestID mismatch: expected %q, got %q",
			expectedMerchantID, result.Body.StkCallback.MerchantRequestID)
	}
	return nil
}

// WriteOK writes a 200 OK acknowledgement.
func WriteOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, bytes.NewReader([]byte(`{"ResultCode":0,"ResultDesc":"Success"}`)))
}
