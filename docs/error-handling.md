# Error Handling

All API methods return a typed `*daraja.DarajaError` on failure.

## Error Structure

```go
type DarajaError struct {
    HTTPStatus         int      // HTTP status code (e.g. 400, 401, 500)
    ResponseCode       string   // Daraja API response code
    ResponseDesc       string   // Daraja API response description
    Endpoint           string   // The API endpoint that was called
    CorrelationID      string   // Correlation ID for tracing
    Retryable          bool     // Whether the operation can be retried safely
    RawResponseSnippet string   // Truncated raw response (secrets redacted)
    Err                error    // Wrapped root cause
}
```

## Checking Error Type

```go
resp, err := client.STKPush(ctx, req)
if err != nil {
    var darajaErr *daraja.DarajaError
    if errors.As(err, &darajaErr) {
        fmt.Printf("HTTP %d | Code: %s | Desc: %s\n",
            darajaErr.HTTPStatus, darajaErr.ResponseCode, darajaErr.ResponseDesc)
    } else {
        // Network or client-level error
        fmt.Printf("Request failed: %v\n", err)
    }
}
```

## Retryable Errors

Use `errors.As` to check if an error is retryable:

```go
if darajaErr.Retryable {
    // Implement exponential backoff or re-enqueue
}
```

Errors are retryable when:
- HTTP status is `429 Too Many Requests`
- HTTP status is `5xx` (server errors)
- Network/timeout errors occur

Non-retryable errors include:
- `400 Bad Request` (invalid input)
- `401 Unauthorized` (bad credentials — token refresh won't help)
- `404 Not Found`

## Validation Errors

Input validation failures return `*DarajaError` with `HTTPStatus: 0` and a descriptive message. These indicate programming errors and should never be retried.

## Error Wrapping

The client uses `fmt.Errorf("...: %w", err)` to wrap errors, so you can use `errors.Is()` and `errors.As()` for inspection.