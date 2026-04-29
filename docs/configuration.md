# Configuration

The Daraja client uses a functional options pattern for configuration. All options are passed to `daraja.NewClient()`.

## Basic Options

### `WithEnvironment(env)`

Sets the API environment.

```go
daraja.WithEnvironment(daraja.Sandbox)
daraja.WithEnvironment(daraja.Production)
```

Default: `Sandbox`

### `WithCredentials(consumerKey, consumerSecret)`

Sets your Daraja API credentials.

```go
daraja.WithCredentials(os.Getenv("DARAJA_CONSUMER_KEY"), os.Getenv("DARAJA_CONSUMER_SECRET"))
```

### `WithShortCode(shortCode, passkey)`

Sets the business short code and STK Push passkey.

```go
daraja.WithShortCode("174379", os.Getenv("DARAJA_PASSKEY"))
```

## Advanced Options

### `WithHTTPClient(httpClient)`

Provides a custom `*http.Client`. The default is a secure client with TLS 1.2+, 30s timeouts, and connection pooling.

```go
custom := &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 50,
    },
}
daraja.WithHTTPClient(custom)
```

### `WithTracer(tracer)`

Attaches an OpenTelemetry-compatible tracer for distributed tracing.

```go
import "github.com/Kisotu/daraja-client-go/internal/trace"

daraja.WithTracer(trace.NewNoopTracer())
```

## Environment Variables Reference

| Variable | Required For | Description |
|---|---|---|
| `DARAJA_CONSUMER_KEY` | All operations | API consumer key |
| `DARAJA_CONSUMER_SECRET` | All operations | API consumer secret |
| `DARAJA_SHORTCODE` | STK Push, B2C | Business short code |
| `DARAJA_PASSKEY` | STK Push | Lipa Na M-Pesa Online passkey |
| `DARAJA_INITIATOR` | Transaction Status, Reversal | API initiator username |
| `DARAJA_SECURITY_CREDENTIAL` | Transaction Status, Reversal | Base64-encoded security credential |
| `DARAJA_ENVIRONMENT` | All operations | `sandbox` (default) or `production` |

## Production Checklist

Before switching to production:

1. Verify TLS 1.2+ is enforced (default client enforces this)
2. Set `WithEnvironment(daraja.Production)` 
3. Use real short codes and credentials
4. Implement callback URL verification with IP allowlist
5. Review [Security](security.md) for go-live requirements