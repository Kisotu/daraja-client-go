# API Compatibility Checklist

This document captures the public API surface of `daraja-client-go` at v1.0.0.
Any changes to these items must be considered breaking and require a MAJOR
version bump per Semantic Versioning.

## Public Interface: `Client`

File: `daraja.go`

| Method | Signature | Status |
|---|---|---|
| `Token` | `Token(ctx context.Context) (string, error)` | Frozen |
| `STKPush` | `STKPush(ctx context.Context, req STKPushRequest) (*STKPushResponse, error)` | Frozen |
| `STKPushQuery` | `STKPushQuery(ctx context.Context, req STKPushQueryRequest) (*STKPushQueryResponse, error)` | Frozen |
| `C2BSimulate` | `C2BSimulate(ctx context.Context, req C2BSimulationRequest) (*C2BResponse, error)` | Frozen |
| `C2BRegisterURL` | `C2BRegisterURL(ctx context.Context, req RegisterURLRequest) (*RegisterURLResponse, error)` | Frozen |
| `B2CPayment` | `B2CPayment(ctx context.Context, req B2CPaymentRequest) (*B2CPaymentResponse, error)` | Frozen |
| `TransactionStatus` | `TransactionStatus(ctx context.Context, req TransactionStatusRequest) (*TransactionStatusResponse, error)` | Frozen |
| `Reversal` | `Reversal(ctx context.Context, req ReversalRequest) (*ReversalResponse, error)` | Frozen |

## Public Types

### Request Structs

| Type | File | Status |
|---|---|---|
| `STKPushRequest` | `stkpush.go` | Frozen |
| `STKPushQueryRequest` | `stkpush.go` | Frozen |
| `C2BSimulationRequest` | `c2b.go` | Frozen |
| `RegisterURLRequest` | `c2b.go` | Frozen |
| `B2CPaymentRequest` | `b2c.go` | Frozen |
| `TransactionStatusRequest` | `transaction_status.go` | Frozen |
| `ReversalRequest` | `reversal.go` | Frozen |

### Response Structs

| Type | File | Status |
|---|---|---|
| `STKPushResponse` | `stkpush.go` | Frozen |
| `STKPushQueryResponse` | `stkpush.go` | Frozen |
| `C2BResponse` | `c2b.go` | Frozen |
| `RegisterURLResponse` | `c2b.go` | Frozen |
| `B2CPaymentResponse` | `b2c.go` | Frozen |
| `TransactionStatusResponse` | `transaction_status.go` | Frozen |
| `ReversalResponse` | `reversal.go` | Frozen |

### Configuration Types

| Type | File | Status |
|---|---|---|
| `Environment` | `config.go` | Frozen |
| `Config` | `config.go` | Frozen |
| `Option` | `config.go` | Frozen |

### Error Types

| Type | File | Status |
|---|---|---|
| `DarajaError` | `errors.go` | Frozen |

## Public Functions

| Function | File | Status |
|---|---|---|
| `NewClient` | `daraja.go` | Frozen |
| `WithEnvironment` | `config.go` | Frozen |
| `WithCredentials` | `config.go` | Frozen |
| `WithShortCode` | `config.go` | Frozen |
| `WithHTTPClient` | `config.go` | Frozen |
| `WithTracer` | `config.go` | Frozen |

## Public Constants/Variables

| Name | File | Status |
|---|---|---|
| `Sandbox` | `config.go` | Frozen |
| `Production` | `config.go` | Frozen |

## Public Sub-packages

### `auth`
| Symbol | Status |
|---|---|
| `AuthConfig` | Frozen |
| `AuthManager` | Frozen |
| `NewAuthManager` | Frozen |
| `Token` | Frozen |
| `TokenResponse` | Frozen |

### `callback`
| Symbol | Status |
|---|---|
| `Verifier` | Frozen |
| `NewVerifier` | Frozen |
| `WithIPAllowlist` | Frozen |
| `WithIPCheck` | Frozen |
| `VerifySTKPush` | Frozen |
| `ValidateCallbackPayload` | Frozen |
| `ConstantTimeEqual` | Frozen |
| `WriteOK` | Frozen |

### `transport`
| Symbol | Status |
|---|---|
| `NewSecureClient` | Frozen |
| `RetryTransport` | Frozen |
| `NewRetryTransport` | Frozen |
| `DefaultRetryableFunc` | Frozen |
| `RetryableFunc` | Frozen |
| `RetryOption` | Frozen |

## Deprecation Policy

1. Deprecated symbols will be marked with a Go doc comment containing `Deprecated:`
2. Deprecated symbols will be maintained for at least 2 MINOR versions
3. Removal only happens in a MAJOR version bump

## Exceptions

The following may change in MINOR releases:
- Internal packages under `internal/` (no compatibility guarantees)
- New fields added to request/response structs (backward-compatible additions)
- New methods added to the `Client` interface (new methods are additive)