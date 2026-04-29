package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kisotu/daraja-client-go"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	loadEnvVars()

	switch os.Args[1] {
	case "auth":
		runAuth()
	case "stkpush":
		runSTKPush()
	case "status":
		runStatus()
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`daraja-cli — Safaricom Daraja (M-Pesa) CLI

Usage:
  daraja-cli <command> [flags]

Commands:
  auth              Test authentication and print token
  stkpush           Send an STK Push request
  status            Query transaction status
  help              Show this help

Environment:
  DARAJA_CONSUMER_KEY      (required)
  DARAJA_CONSUMER_SECRET   (required)
  DARAJA_SHORTCODE         (required for stkpush)
  DARAJA_PASSKEY           (required for stkpush)
  DARAJA_ENVIRONMENT       (sandbox or production, default: sandbox)

Examples:
  daraja-cli auth
  daraja-cli stkpush -phone 254712345678 -amount 1 -account "Test123"
  daraja-cli status -transaction-id ABC123`)
}

func loadEnvVars() {
	flag.Parse()
}

func newClient() (daraja.Client, error) {
	consumerKey := os.Getenv("DARAJA_CONSUMER_KEY")
	consumerSecret := os.Getenv("DARAJA_CONSUMER_SECRET")
	if consumerKey == "" || consumerSecret == "" {
		return nil, fmt.Errorf("DARAJA_CONSUMER_KEY and DARAJA_CONSUMER_SECRET must be set")
	}

	env := daraja.Sandbox
	if os.Getenv("DARAJA_ENVIRONMENT") == "production" {
		env = daraja.Production
	}

	return daraja.NewClient(
		daraja.WithEnvironment(env),
		daraja.WithCredentials(consumerKey, consumerSecret),
	)
}

func runAuth() {
	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := client.Token(ctx)
	if err != nil {
		log.Fatalf("authentication failed: %v", err)
	}

	fmt.Println("Authentication successful")
	fmt.Printf("Token (first 16 chars): %s...\n", token[:min(len(token), 16)])
}

func runSTKPush() {
	fs := flag.NewFlagSet("stkpush", flag.ExitOnError)
	phone := fs.String("phone", "", "Customer phone number (e.g. 254712345678)")
	amount := fs.String("amount", "1", "Amount to charge")
	account := fs.String("account", "CLITest", "Account reference")
	desc := fs.String("desc", "CLI STK Push", "Transaction description")
	callback := fs.String("callback", "https://mydomain.com/callback", "Callback URL")
	shortCode := fs.String("shortcode", "", "Business short code (overrides DARAJA_SHORTCODE)")
	_ = fs.Parse(os.Args[2:])

	if *phone == "" {
		log.Fatal("--phone is required")
	}

	shortcode := *shortCode
	if shortcode == "" {
		shortcode = os.Getenv("DARAJA_SHORTCODE")
	}
	if shortcode == "" {
		log.Fatal("short code required: set DARAJA_SHORTCODE or pass --shortcode")
	}

	passkey := os.Getenv("DARAJA_PASSKEY")
	if passkey == "" {
		log.Fatal("DARAJA_PASSKEY environment variable is required")
	}

	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.STKPush(ctx, daraja.STKPushRequest{
		BusinessShortCode: shortcode,
		Amount:            *amount,
		PartyA:            *phone,
		PartyB:            shortcode,
		PhoneNumber:       *phone,
		CallBackURL:       *callback,
		AccountReference:  *account,
		TransactionDesc:   *desc,
	})
	if err != nil {
		log.Fatalf("STK push failed: %v", err)
	}

	printJSON(resp)
}

func runStatus() {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	transactionID := fs.String("transaction-id", "", "Transaction ID to query")
	_ = fs.Parse(os.Args[2:])

	if *transactionID == "" {
		log.Fatal("--transaction-id is required")
	}

	initiator := os.Getenv("DARAJA_INITIATOR")
	credential := os.Getenv("DARAJA_SECURITY_CREDENTIAL")
	shortCode := os.Getenv("DARAJA_SHORTCODE")
	if initiator == "" || credential == "" || shortCode == "" {
		log.Fatal("DARAJA_INITIATOR, DARAJA_SECURITY_CREDENTIAL, and DARAJA_SHORTCODE must be set")
	}

	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.TransactionStatus(ctx, daraja.TransactionStatusRequest{
		Initiator:          initiator,
		SecurityCredential: credential,
		CommandID:          "TransactionStatusQuery",
		TransactionID:      *transactionID,
		PartyA:             shortCode,
		IdentifierType:     "1",
		ResultURL:          "https://mydomain.com/result",
		QueueTimeOutURL:    "https://mydomain.com/timeout",
		Remarks:            "CLI status query",
	})
	if err != nil {
		log.Fatalf("status query failed: %v", err)
	}

	printJSON(resp)
}

func printJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Fatalf("failed to encode output: %v", err)
	}
}