# Incident Response

This document outlines the incident response process for production issues
with `daraja-client-go` or services using it.

## Severity Levels

| Level | Definition | Response Time |
|---|---|---|
| **SEV1** | Production outage — STK Push / payments failing | 30 minutes |
| **SEV2** | Partial degradation — one endpoint failing | 2 hours |
| **SEV3** | Non-critical — bug with workaround | 1 business day |

## Response Steps

1. **Identify**: Check metrics, logs, and error rates to confirm the issue
2. **Contain**: If the issue is in client config, roll back the deployment.
   If the issue is in the Daraja API, activate circuit breaker or fallback.
3. **Diagnose**: Review DarajaError codes, HTTP statuses, and retry counts.
   Check if the issue is retryable or requires credential rotation.
4. **Resolve**: Apply fix, verify in sandbox, deploy to production.
5. **Learn**: File a postmortem for SEV1/SEV2 incidents.

## Rollback

```bash
# Revert to previous version
go get github.com/Kisotu/daraja-client-go@v0.9.0
```

## Communication

- Post incident updates in the team's communication channel
- For SEV1, notify stakeholders within 15 minutes of confirmation
- After resolution, share a postmortem with root cause and preventive actions