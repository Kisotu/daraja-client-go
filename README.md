# daraja-client-go

A production-ready, open-source Go client for Safaricom Daraja (M-Pesa) APIs.

## Features

- OAuth token management with caching and concurrent refresh protection
- Secure-by-default HTTP transport (TLS 1.2+, timeouts, context-aware)
- STK Push, C2B, B2C, Transaction Status, Reversal, and more
- Strong error model with typed DarajaError
- Functional options pattern for configuration
- Sandbox and production environment support

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/anomalyco/daraja-client-go"
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

## Installation

```bash
go get github.com/anomalyco/daraja-client-go
```

## Requirements

- Go 1.23 or later
- Safaricom Daraja API credentials

## License

MIT - see [LICENSE](LICENSE)