package daraja

import "fmt"

// DarajaError represents an error returned by the Daraja API.
type DarajaError struct {
	HTTPStatus         int
	ResponseCode       string
	ResponseDesc       string
	Endpoint           string
	CorrelationID      string
	Retryable          bool
	RawResponseSnippet string
	Err                error
}

func (e *DarajaError) Error() string {
	msg := fmt.Sprintf("daraja: %s (HTTP %d, code: %s)", e.ResponseDesc, e.HTTPStatus, e.ResponseCode)
	if e.Endpoint != "" {
		msg += fmt.Sprintf(" endpoint: %s", e.Endpoint)
	}
	return msg
}

func (e *DarajaError) Unwrap() error {
	return e.Err
}
