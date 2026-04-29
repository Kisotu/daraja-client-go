package daraja

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// testRoundTripper routes all requests to a single test server.
type testRoundTripper struct {
	server *httptest.Server
}

func (t *testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// Rewrite the URL to point to the test server
	u := *r.URL
	u.Scheme = "http"
	u.Host = t.server.Listener.Addr().String()
	r2 := r.Clone(r.Context())
	r2.URL = &u
	return t.server.Client().Transport.RoundTrip(r2)
}

func TestConcurrentSTKPush(t *testing.T) {
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

	client := &http.Client{Transport: &testRoundTripper{server: srv}}

	c, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials("key", "secret"),
		WithShortCode("174379", "passkey"),
		WithHTTPClient(client),
	)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	errs := make(chan error, 50)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := c.STKPush(ctx, STKPushRequest{
				BusinessShortCode: "174379",
				Amount:            "1",
				PartyA:            "254712345678",
				PartyB:            "174379",
				PhoneNumber:       "254712345678",
				CallBackURL:       "https://example.com/callback",
				AccountReference:  "Test",
				TransactionDesc:   "Test payment",
			})
			if err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent STK push failed: %v", err)
	}
}

func TestConcurrentC2B(t *testing.T) {
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

	client := &http.Client{Transport: &testRoundTripper{server: srv}}

	c, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials("key", "secret"),
		WithHTTPClient(client),
	)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	errs := make(chan error, 50)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := c.C2BSimulate(ctx, C2BSimulationRequest{
				ShortCode: "174379",
				CommandID: "CustomerPayBillOnline",
				Amount:    "100",
				Msisdn:    "254712345678",
				BillRefNumber: "INV-001",
			})
			if err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent C2B failed: %v", err)
	}
}

func TestConcurrentMixedWorkload(t *testing.T) {
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
		case "/mpesa/c2b/v1/simulate":
			w.Write([]byte(`{
				"ConversationID":"conv-1",
				"OriginatorConversationID":"orig-1",
				"ResponseCode":"0",
				"ResponseDescription":"Success"
			}`))
		case "/mpesa/b2c/v1/paymentrequest":
			w.Write([]byte(`{
				"ConversationID":"conv-1",
				"OriginatorConversationID":"orig-1",
				"ResponseCode":"0",
				"ResponseDescription":"Success"
			}`))
		case "/mpesa/transactionstatus/v1/query":
			w.Write([]byte(`{
				"ConversationID":"conv-1",
				"OriginatorConversationID":"orig-1",
				"ResponseCode":"0",
				"ResponseDescription":"Success"
			}`))
		case "/mpesa/reversal/v1/request":
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
		WithShortCode("174379", "passkey"),
		WithHTTPClient(httpClient),
	)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	errs := make(chan error, 200)

	ops := []func(){
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := c.STKPush(ctx, STKPushRequest{
				BusinessShortCode: "174379", Amount: "1", PartyA: "254712345678",
				PartyB: "174379", PhoneNumber: "254712345678",
				CallBackURL: "https://example.com/callback", AccountReference: "Test",
				TransactionDesc: "Test",
			})
			if err != nil {
				errs <- err
			}
		},
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := c.C2BSimulate(ctx, C2BSimulationRequest{
				ShortCode: "174379", CommandID: "CustomerPayBillOnline",
				Amount: "100", Msisdn: "254712345678", BillRefNumber: "INV-001",
			})
			if err != nil {
				errs <- err
			}
		},
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := c.B2CPayment(ctx, B2CPaymentRequest{
				InitiatorName: "testapi", SecurityCredential: "cred",
				CommandID: "BusinessPayment", Amount: "100", PartyA: "174379",
				PartyB: "254712345678", Remarks: "Test",
				QueueTimeOutURL: "https://example.com/timeout",
				ResultURL:       "https://example.com/result",
			})
			if err != nil {
				errs <- err
			}
		},
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := c.TransactionStatus(ctx, TransactionStatusRequest{
				Initiator: "testapi", SecurityCredential: "cred",
				CommandID: "TransactionStatusQuery", TransactionID: "LHG123",
				PartyA: "174379", IdentifierType: "1",
				ResultURL:       "https://example.com/result",
				QueueTimeOutURL: "https://example.com/timeout", Remarks: "Test",
			})
			if err != nil {
				errs <- err
			}
		},
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := c.Reversal(ctx, ReversalRequest{
				Initiator: "testapi", SecurityCredential: "cred",
				CommandID: "TransactionReversal", TransactionID: "LHG123",
				PartyA: "174379", IdentifierType: "11",
				QueueTimeOutURL: "https://example.com/timeout",
				ResultURL: "https://example.com/result", Remarks: "Test",
			})
			if err != nil {
				errs <- err
			}
		},
	}

	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ops[idx%len(ops)]()
		}(i)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("mixed workload failed: %v", err)
	}
}