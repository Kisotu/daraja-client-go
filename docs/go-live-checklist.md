# Go-Live Checklist

Use this checklist when moving from sandbox testing to production.

## Configuration

- [ ] Environment set to `daraja.Production`
- [ ] Production consumer key and secret configured
- [ ] Real short code configured (not `174379`)
- [ ] Production passkey configured
- [ ] Security credentials configured for admin operations

## Callback Endpoints

- [ ] Callback URLs are HTTPS with valid TLS certificates
- [ ] IP allowlist configured for Daraja source IPs
- [ ] Correlation ID validation implemented
- [ ] Payload schema validation implemented
- [ ] Acknowledgment responses (`callback.WriteOK`) returned
- [ ] Idempotency/replay protection in place
- [ ] Callback endpoints are monitored and alert on errors

## Resilience

- [ ] Retry policy tuned for production tolerance
- [ ] Circuit breaker thresholds configured
- [ ] Timeouts set appropriate to SLA expectations
- [ ] Graceful shutdown implemented for background workers

## Observability

- [ ] Logger configured with production output (JSON)
- [ ] Metrics endpoint accessible and scraped
- [ ] Tracing instrumentation active (if using OTel)
- [ ] Correlation IDs propagated across service boundaries
- [ ] Logs verified for secret redaction

## Security

- [ ] Credentials loaded from secure storage (not hardcoded)
- [ ] TLS 1.2+ enforced
- [ ] `govulncheck` passes with no findings
- [ ] `gosec` passes with no high/critical findings

## Testing

- [ ] Unit tests pass (`go test ./... -race`)
- [ ] Integration tests pass against sandbox
- [ ] Smoke tests pass against production (small amounts)
- [ ] Rollback plan documented