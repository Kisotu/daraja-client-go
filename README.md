# daraja-client-go

A production-ready, open-source Go client for [Safaricom Daraja (M-Pesa) APIs](https://developer.safaricom.co.ke/).

[![CI](https://github.com/Kisotu/daraja-client-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Kisotu/daraja-client-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Kisotu/daraja-client-go.svg)](https://pkg.go.dev/github.com/Kisotu/daraja-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kisotu/daraja-client-go)](https://goreportcard.com/report/github.com/Kisotu/daraja-client-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Features

- **Complete API coverage**: STK Push, C2B, B2C, Transaction Status, Reversal
- **Secure by default**: TLS 1.2+, credential redaction, callback verification
- **Resilient**: Exponential backoff with jitter, circuit breaker per endpoint group
- **Observable**: Structured logging, Prometheus metrics, OpenTelemetry tracing
- **Idiomatic Go**: Functional options, clean interfaces, no magic

## Installation

```bash
go get github.com/Kisotu/daraja-client-go
```

## Quick Start

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

## CLI

A CLI tool is included for testing and operations:

```bash
# Install
go install ./cmd/daraja-cli

# Test authentication
export DARAJA_CONSUMER_KEY=your_key
export DARAJA_CONSUMER_SECRET=your_secret
daraja-cli auth

# Send STK Push
daraja-cli stkpush -phone 254712345678 -amount 1

# Query transaction status
daraja-cli status --transaction-id LHG123ABC
```

## Documentation

- [Getting Started](docs/getting-started.md) — sandbox setup and first request
- [Configuration](docs/configuration.md) — all options and defaults
- [Security](docs/security.md) — callback verification, secret handling, go-live checklist
- [Error Handling](docs/error-handling.md) — error taxonomy and retry guidance
- [Examples](docs/examples.md) — practical integration patterns
- [Release Process](docs/release-process.md) — versioning and publishing
- [Go-Live Checklist](docs/go-live-checklist.md) — production readiness
- [Architecture Decisions](docs/adr/) — design rationale

## API Coverage

| Endpoint | Status |
|---|---|
| STK Push (Lipa Na M-Pesa Online) | ✅ |
| STK Push Query | ✅ |
| C2B Simulate | ✅ |
| C2B Register URL | ✅ |
| B2C Payment | ✅ |
| Transaction Status | ✅ |
| Reversal | ✅ |

## Requirements

- Go 1.23 or later
- Safaricom Daraja API credentials

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) and [our PR template](.github/pull_request_template.md). All contributions are welcome!

## Security

Found a vulnerability? See [SECURITY.md](SECURITY.md) for our disclosure process.

## License

MIT — see [LICENSE](LICENSE).

## Support

- [Documentation](docs/getting-started.md)
- [Issue tracker](https://github.com/Kisotu/daraja-client-go/issues)
- [SUPPORT.md](SUPPORT.md)