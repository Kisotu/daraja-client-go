//go:build integration

package daraja

import (
	"context"
	"os"
	"testing"
)

func TestSTKPush_Integration(t *testing.T) {
	creds := loadCreds()
	if !creds.ok() {
		t.Skip("Daraja credentials not set")
	}

	client := newTestClient(t, creds)

	resp, err := client.STKPush(context.Background(), STKPushRequest{
		BusinessShortCode: creds.shortCode,
		Amount:            "1",
		PartyA:            "254708374149",
		PartyB:            creds.shortCode,
		PhoneNumber:       "254708374149",
		CallBackURL:       "https://example.com/callback",
		AccountReference:  "Integration Test",
		TransactionDesc:   "Smoke test",
	})
	if err != nil {
		t.Fatalf("STKPush() error = %v", err)
	}

	if resp.CheckoutRequestID == "" {
		t.Error("expected non-empty CheckoutRequestID")
	}
	if resp.MerchantRequestID == "" {
		t.Error("expected non-empty MerchantRequestID")
	}
}

func TestSTKPushQuery_Integration(t *testing.T) {
	creds := loadCreds()
	if !creds.ok() {
		t.Skip("Daraja credentials not set")
	}

	client := newTestClient(t, creds)

	pushResp, err := client.STKPush(context.Background(), STKPushRequest{
		BusinessShortCode: creds.shortCode,
		Amount:            "1",
		PartyA:            "254708374149",
		PartyB:            creds.shortCode,
		PhoneNumber:       "254708374149",
		CallBackURL:       "https://example.com/callback",
		AccountReference:  "Integration Test",
		TransactionDesc:   "Query smoke test",
	})
	if err != nil {
		t.Fatalf("STKPush() error = %v", err)
	}

	_, err = client.STKPushQuery(context.Background(), STKPushQueryRequest{
		BusinessShortCode: creds.shortCode,
		CheckoutRequestID: pushResp.CheckoutRequestID,
	})
	if err != nil {
		t.Logf("STKPushQuery returned expected sandbox error: %v", err)
	}
}

func TestC2BSimulate_Integration(t *testing.T) {
	creds := loadCreds()
	if !creds.ok() {
		t.Skip("Daraja credentials not set")
	}

	client := newTestClient(t, creds)

	resp, err := client.C2BSimulate(context.Background(), C2BSimulationRequest{
		ShortCode:     creds.shortCode,
		CommandID:     "CustomerPayBillOnline",
		Amount:        "10",
		Msisdn:        "254708374149",
		BillRefNumber: "INV-TEST-001",
	})
	if err != nil {
		t.Fatalf("C2BSimulate() error = %v", err)
	}

	if resp.ConversationID == "" {
		t.Error("expected non-empty ConversationID")
	}
	if resp.ResponseCode != "0" {
		t.Errorf("expected response code 0, got %s: %s", resp.ResponseCode, resp.ResponseDescription)
	}
}

func TestB2CPayment_Integration(t *testing.T) {
	creds := loadCreds()
	if !creds.ok() {
		t.Skip("Daraja credentials not set")
	}

	client := newTestClient(t, creds)

	resp, err := client.B2CPayment(context.Background(), B2CPaymentRequest{
		InitiatorName:      "testapi",
		SecurityCredential: os.Getenv("DARAJA_B2C_SECURITY_CREDENTIAL"),
		CommandID:          "BusinessPayment",
		Amount:             "10",
		PartyA:             creds.shortCode,
		PartyB:             "254708374149",
		Remarks:            "B2C Integration test",
		QueueTimeOutURL:    "https://example.com/timeout",
		ResultURL:          "https://example.com/result",
	})
	if err != nil {
		t.Fatalf("B2CPayment() error = %v", err)
	}

	if resp.ConversationID == "" {
		t.Error("expected non-empty ConversationID")
	}
	if resp.ResponseCode != "0" {
		t.Errorf("expected response code 0, got %s: %s", resp.ResponseCode, resp.ResponseDescription)
	}
}

func TestTransactionStatus_Integration(t *testing.T) {
	creds := loadCreds()
	if !creds.ok() {
		t.Skip("Daraja credentials not set")
	}

	client := newTestClient(t, creds)

	resp, err := client.TransactionStatus(context.Background(), TransactionStatusRequest{
		Initiator:          "testapi",
		SecurityCredential: os.Getenv("DARAJA_SECURITY_CREDENTIAL"),
		CommandID:          "TransactionStatusQuery",
		TransactionID:      "LHG0000000",
		PartyA:             creds.shortCode,
		IdentifierType:     "1",
		ResultURL:          "https://example.com/result",
		QueueTimeOutURL:    "https://example.com/timeout",
		Remarks:            "Status integration test",
	})
	if err != nil {
		t.Fatalf("TransactionStatus() error = %v", err)
	}

	if resp.ConversationID == "" {
		t.Error("expected non-empty ConversationID")
	}
}

func TestReversal_Integration(t *testing.T) {
	creds := loadCreds()
	if !creds.ok() {
		t.Skip("Daraja credentials not set")
	}

	client := newTestClient(t, creds)

	_, err := client.Reversal(context.Background(), ReversalRequest{
		Initiator:          "testapi",
		SecurityCredential: os.Getenv("DARAJA_SECURITY_CREDENTIAL"),
		CommandID:          "TransactionReversal",
		TransactionID:      "LHG0000000",
		PartyA:             creds.shortCode,
		IdentifierType:     "1",
		ResultURL:          "https://example.com/result",
		QueueTimeOutURL:    "https://example.com/timeout",
		Remarks:            "Reversal integration test",
	})
	if err != nil {
		t.Logf("Reversal returned expected sandbox error: %v", err)
	}
}

type testCreds struct {
	consumerKey    string
	consumerSecret string
	shortCode      string
	passkey        string
}

func (tc testCreds) ok() bool {
	return tc.consumerKey != "" && tc.consumerSecret != "" &&
		tc.shortCode != "" && tc.passkey != ""
}

func loadCreds() testCreds {
	return testCreds{
		consumerKey:    os.Getenv("DARAJA_CONSUMER_KEY"),
		consumerSecret: os.Getenv("DARAJA_CONSUMER_SECRET"),
		shortCode:      os.Getenv("DARAJA_SHORTCODE"),
		passkey:        os.Getenv("DARAJA_PASSKEY"),
	}
}

func newTestClient(t *testing.T, creds testCreds) Client {
	t.Helper()
	client, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials(creds.consumerKey, creds.consumerSecret),
		WithShortCode(creds.shortCode, creds.passkey),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return client
}