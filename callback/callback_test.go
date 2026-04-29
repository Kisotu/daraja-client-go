package callback

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifySTKPush_Valid(t *testing.T) {
	payload := STKPushResult{}
	payload.Body.StkCallback.MerchantRequestID = "mr-001"
	payload.Body.StkCallback.CheckoutRequestID = "chk-001"
	payload.Body.StkCallback.ResultCode = 0
	payload.Body.StkCallback.ResultDesc = "Success"

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	v := NewVerifier()
	result, err := v.VerifySTKPush(req)
	if err != nil {
		t.Fatalf("VerifySTKPush() error = %v", err)
	}
	if result.Body.StkCallback.CheckoutRequestID != "chk-001" {
		t.Errorf("CheckoutRequestID = %q, want %q", result.Body.StkCallback.CheckoutRequestID, "chk-001")
	}
}

func TestVerifySTKPush_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/callback", nil)
	v := NewVerifier()
	_, err := v.VerifySTKPush(req)
	if err == nil {
		t.Fatal("expected error for non-POST method")
	}
}

func TestVerifySTKPush_InvalidContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "text/plain")
	v := NewVerifier()
	_, err := v.VerifySTKPush(req)
	if err == nil {
		t.Fatal("expected error for non-JSON content type")
	}
}

func TestVerifySTKPush_MissingCheckoutID(t *testing.T) {
	payload := STKPushResult{}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	v := NewVerifier()
	_, err := v.VerifySTKPush(req)
	if err == nil {
		t.Fatal("expected error for missing CheckoutRequestID")
	}
}

func TestVerifySTKPush_MalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader([]byte(`not json`)))
	req.Header.Set("Content-Type", "application/json")

	v := NewVerifier()
	_, err := v.VerifySTKPush(req)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestVerifySTKPush_IPAllowlist(t *testing.T) {
	payload := STKPushResult{}
	payload.Body.StkCallback.CheckoutRequestID = "chk-001"
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.168.1.1:1234"

	v := NewVerifier(WithIPAllowlist([]string{"10.0.0.0/8"}))
	_, err := v.VerifySTKPush(req)
	if err == nil {
		t.Fatal("expected error for non-matching IP")
	}
}

func TestValidateCallbackPayload(t *testing.T) {
	result := &STKPushResult{}
	result.Body.StkCallback.CheckoutRequestID = "chk-001"
	result.Body.StkCallback.MerchantRequestID = "mr-001"

	t.Run("matching IDs", func(t *testing.T) {
		if err := ValidateCallbackPayload(result, "chk-001", "mr-001"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("mismatched CheckoutID", func(t *testing.T) {
		if err := ValidateCallbackPayload(result, "chk-002", "mr-001"); err == nil {
			t.Fatal("expected error for mismatched CheckoutRequestID")
		}
	})
	t.Run("mismatched MerchantID", func(t *testing.T) {
		if err := ValidateCallbackPayload(result, "chk-001", "mr-002"); err == nil {
			t.Fatal("expected error for mismatched MerchantRequestID")
		}
	})
}

func TestConstantTimeEqual(t *testing.T) {
	if !ConstantTimeEqual([]byte("abc"), []byte("abc")) {
		t.Error("equal slices should match")
	}
	if ConstantTimeEqual([]byte("abc"), []byte("def")) {
		t.Error("different slices should not match")
	}
	if ConstantTimeEqual([]byte("abc"), []byte("abcd")) {
		t.Error("different-length slices should not match")
	}
}

func TestWriteOK(t *testing.T) {
	w := httptest.NewRecorder()
	WriteOK(w)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", w.Header().Get("Content-Type"))
	}
}
