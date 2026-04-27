//go:build integration

package daraja

import (
	"context"
	"os"
	"testing"
)

func TestSTKPush_Integration(t *testing.T) {
	consumerKey := os.Getenv("DARAJA_CONSUMER_KEY")
	consumerSecret := os.Getenv("DARAJA_CONSUMER_SECRET")
	shortCode := os.Getenv("DARAJA_SHORTCODE")
	passkey := os.Getenv("DARAJA_PASSKEY")

	if consumerKey == "" || consumerSecret == "" || shortCode == "" || passkey == "" {
		t.Skip("DARAJA_CONSUMER_KEY, DARAJA_CONSUMER_SECRET, DARAJA_SHORTCODE, and DARAJA_PASSKEY must be set")
	}

	client, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials(consumerKey, consumerSecret),
		WithShortCode(shortCode, passkey),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.STKPush(context.Background(), STKPushRequest{
		BusinessShortCode: shortCode,
		Amount:            "1",
		PartyA:            "254708374149",
		PartyB:            shortCode,
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