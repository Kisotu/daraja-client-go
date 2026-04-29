# Getting Started

This guide walks through setting up the Daraja client and making your first API call to the Safaricom sandbox.

## Prerequisites

- Go 1.23 or later
- Safaricom Daraja API credentials ([register on the Daraja portal](https://developer.safaricom.co.ke/))

## Obtaining Credentials

1. Sign up at the [Safaricom Developer Portal](https://developer.safaricom.co.ke/)
2. Create an app to get your **Consumer Key** and **Consumer Secret**
3. For STK Push, request a **short code** and **passkey** from the test credentials page
4. Set up a callback URL (use a service like ngrok for local development)

## Installing the Library

```bash
go get github.com/Kisotu/daraja-client-go
```

## Minimal Example

Create a file `main.go`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Kisotu/daraja-client-go"
)

func main() {
	client, err := daraja.NewClient(
		daraja.WithEnvironment(daraja.Sandbox),
		daraja.WithCredentials(
			os.Getenv("DARAJA_CONSUMER_KEY"),
			os.Getenv("DARAJA_CONSUMER_SECRET"),
		),
		daraja.WithShortCode(
			os.Getenv("DARAJA_SHORTCODE"),
			os.Getenv("DARAJA_PASSKEY"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.STKPush(context.Background(), daraja.STKPushRequest{
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
		log.Fatal(err)
	}
	fmt.Printf("CheckoutRequestID: %s\n", resp.CheckoutRequestID)
}
```

Run it:

```bash
export DARAJA_CONSUMER_KEY=your_key
export DARAJA_CONSUMER_SECRET=your_secret
export DARAJA_SHORTCODE=174379
export DARAJA_PASSKEY=your_passkey
go run main.go
```

## Next Steps

- See [Configuration](configuration.md) for all available options
- Read [Security](security.md) for callback verification and secret handling
- Browse [Examples](examples.md) for real-world integration patterns
- Run the CLI: `go run ./cmd/daraja-cli auth`