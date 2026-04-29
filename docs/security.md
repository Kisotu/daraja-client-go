# Security

This document covers security practices for integrating with Safaricom Daraja APIs using this client.

## Credential Management

- **Never hardcode secrets** in source code. Use environment variables or a secret manager.
- The client accepts credentials via `WithCredentials()` — feed these from secure sources only.
- All sensitive fields (passwords, tokens, secrets) are automatically **redacted** from logs and error messages.

## Transport Security

The default HTTP client enforces:

- **TLS 1.2 minimum** — connections to Daraja reject older protocols
- **30-second timeouts** for connect, TLS handshake, and response headers
- **Connection pooling** with sane defaults (100 idle, 10 per host)

## Callback Verification

Use the `callback` package to secure incoming Daraja callbacks:

### IP Allowlist

Restrict callbacks to known Daraja IP ranges:

```go
v := callback.NewVerifier()
v.WithIPAllowlist([]string{"196.201.214.200/32", "196.201.214.206/32"})
```

### Payload Validation

Verify callback structure, method, and content type:

```go
payload, err := v.VerifySTKPush(r)
if err != nil {
    http.Error(w, "invalid callback", http.StatusBadRequest)
    return
}
```

### Correlation ID Matching

Ensure callbacks match your pending transactions:

```go
if err := v.ValidateCallbackPayload(payload, expectedID); err != nil {
    // Mismatch — possible replay or misrouted callback
}
```

### Constant-Time Comparison

All sensitive comparisons use `crypto/subtle.ConstantTimeCompare` to prevent timing attacks.

## Secrets in Logs

The logger automatically redacts these fields from all output:

- `password`, `passkey`
- `token`, `authorization`
- `secret`, `consumer_secret`
- `api_key`, `credential`
- URL credentials

## Go-Live Security Checklist

- [ ] Consumer key/secret loaded from secure storage (env vars, vault, KMS)
- [ ] Short code and passkey configured correctly
- [ ] Callback endpoint validates source IP (Daraja IP ranges)
- [ ] Callback correlation ID checking implemented
- [ ] HTTPS enforced for all endpoints
- [ ] Logging configured with redaction verified
- [ ] Dependencies scanned for vulnerabilities (`govulncheck`)
- [ ] Circuit breaker and retry policies tuned for production