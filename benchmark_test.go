package daraja

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkSTKPush(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/v1/generate":
			w.Write([]byte(`{"access_token":"test-token","expires_in":"3599"}`))
		case "/mpesa/stkpush/v1/processrequest":
			w.Write([]byte(`{
				"MerchantRequestID":"mreq-1",
				"CheckoutRequestID":"chk-1",
				"ResponseCode":"0",
				"ResponseDescription":"Success",
				"CustomerMessage":"Success"
			}`))
		}
	}))
	defer srv.Close()

	httpClient := &http.Client{Transport: &testRoundTripper{server: srv}}

	c, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials("key", "secret"),
		WithShortCode("174379", "passkey"),
		WithHTTPClient(httpClient),
	)
	if err != nil {
		b.Fatal(err)
	}

	req := STKPushRequest{
		BusinessShortCode: "174379",
		Amount:            "1",
		PartyA:            "254712345678",
		PartyB:            "174379",
		PhoneNumber:       "254712345678",
		CallBackURL:       "https://example.com/callback",
		AccountReference:  "Test",
		TransactionDesc:   "Test payment",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		c.STKPush(ctx, req)
		cancel()
	}
}

func BenchmarkSTKPushQuery(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/v1/generate":
			w.Write([]byte(`{"access_token":"test-token","expires_in":"3599"}`))
		case "/mpesa/stkpushquery/v1/query":
			w.Write([]byte(`{
				"ResponseCode":"0",
				"ResponseDescription":"Success",
				"MerchantRequestID":"mreq-1",
				"CheckoutRequestID":"chk-1",
				"ResultCode":"0",
				"ResultDesc":"Success"
			}`))
		}
	}))
	defer srv.Close()

	httpClient := &http.Client{Transport: &testRoundTripper{server: srv}}

	c, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials("key", "secret"),
		WithShortCode("174379", "passkey"),
		WithHTTPClient(httpClient),
	)
	if err != nil {
		b.Fatal(err)
	}

	req := STKPushQueryRequest{
		BusinessShortCode: "174379",
		CheckoutRequestID: "chk-1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		c.STKPushQuery(ctx, req)
		cancel()
	}
}

func BenchmarkC2BSimulate(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/v1/generate":
			w.Write([]byte(`{"access_token":"test-token","expires_in":"3599"}`))
		case "/mpesa/c2b/v1/simulate":
			w.Write([]byte(`{
				"ConversationID":"conv-1",
				"OriginatorConversationID":"orig-1",
				"ResponseCode":"0",
				"ResponseDescription":"Success"
			}`))
		}
	}))
	defer srv.Close()

	httpClient := &http.Client{Transport: &testRoundTripper{server: srv}}

	c, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials("key", "secret"),
		WithHTTPClient(httpClient),
	)
	if err != nil {
		b.Fatal(err)
	}

	req := C2BSimulationRequest{
		ShortCode: "174379",
		CommandID: "CustomerPayBillOnline",
		Amount:    "100",
		Msisdn:    "254712345678",
		BillRefNumber: "INV-001",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		c.C2BSimulate(ctx, req)
		cancel()
	}
}

func BenchmarkTokenRefresh(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"access_token":"test-token","expires_in":"3599"}`))
	}))
	defer srv.Close()

	httpClient := &http.Client{Transport: &testRoundTripper{server: srv}}

	c, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials("key", "secret"),
		WithHTTPClient(httpClient),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		c.Token(ctx)
		cancel()
	}
}

func BenchmarkMarshalSTKPush(b *testing.B) {
	req := STKPushRequest{
		BusinessShortCode: "174379",
		Password:          "dGhpcyBpcyBhIHRlc3QgcGFzc3dvcmQ=",
		Timestamp:         "20260429150405",
		TransactionType:   "CustomerPayBillOnline",
		Amount:            "1000",
		PartyA:            "254712345678",
		PartyB:            "174379",
		PhoneNumber:       "254712345678",
		CallBackURL:       "https://example.com/callback",
		AccountReference:  "INV-001",
		TransactionDesc:   "Payment for invoice INV-001",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(req)
	}
}

func BenchmarkUnmarshalSTKPushResponse(b *testing.B) {
	data := []byte(`{
		"MerchantRequestID":"mreq-1",
		"CheckoutRequestID":"chk-1",
		"ResponseCode":"0",
		"ResponseDescription":"Success",
		"CustomerMessage":"Success"
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var resp STKPushResponse
		_ = json.Unmarshal(data, &resp)
	}
}

func BenchmarkGeneratePassword(b *testing.B) {
	shortCode := "174379"
	passkey := "passkey"
	now := time.Now().UTC().Format("20060102150405")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%s%s%s", shortCode, passkey, now)
	}
}

func BenchmarkRandomPhoneNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(99999999))
		_ = fmt.Sprintf("2547%08d", n.Int64())
	}
}

func BenchmarkParallelSTKPush(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/v1/generate":
			w.Write([]byte(`{"access_token":"test-token","expires_in":"3599"}`))
		case "/mpesa/stkpush/v1/processrequest":
			w.Write([]byte(`{
				"MerchantRequestID":"mreq-1",
				"CheckoutRequestID":"chk-1",
				"ResponseCode":"0",
				"ResponseDescription":"Success",
				"CustomerMessage":"Success"
			}`))
		}
	}))
	defer srv.Close()

	httpClient := &http.Client{Transport: &testRoundTripper{server: srv}}

	c, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials("key", "secret"),
		WithShortCode("174379", "passkey"),
		WithHTTPClient(httpClient),
	)
	if err != nil {
		b.Fatal(err)
	}

	req := STKPushRequest{
		BusinessShortCode: "174379",
		Amount:            "1",
		PartyA:            "254712345678",
		PartyB:            "174379",
		PhoneNumber:       "254712345678",
		CallBackURL:       "https://example.com/callback",
		AccountReference:  "Test",
		TransactionDesc:   "Test payment",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			c.STKPush(ctx, req)
			cancel()
		}
	})
}

func BenchmarkDiscardLogger(b *testing.B) {
	// Benchmark logging overhead with io.Discard
	logger := &benchLogger{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Log("test message", "key", "value")
	}
}

type benchLogger struct{}

func (l *benchLogger) Log(msg string, keysAndValues ...any) {
	_ = io.Discard
	_ = msg
	_ = keysAndValues
}