package daraja

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	transactionStatusPath = "/mpesa/transactionstatus/v1/query"
)

// TransactionStatusRequest represents a transaction status query request.
type TransactionStatusRequest struct {
	Initiator              string `json:"Initiator"`
	SecurityCredential     string `json:"SecurityCredential"`
	CommandID              string `json:"CommandID"`
	TransactionID          string `json:"TransactionID"`
	OriginalConversationID string `json:"OriginalConversationID,omitempty"`
	PartyA                 string `json:"PartyA"`
	IdentifierType         string `json:"IdentifierType"`
	ResultURL              string `json:"ResultURL"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL"`
	Remarks                string `json:"Remarks"`
	Occasion               string `json:"Occasion"`
}

// TransactionStatusResponse represents the response from a transaction status
// query.
type TransactionStatusResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

func (c *client) TransactionStatus(ctx context.Context, req TransactionStatusRequest) (*TransactionStatusResponse, error) {
	if err := validateTransactionStatusRequest(&req); err != nil {
		return nil, err
	}

	token, err := c.authMgr.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to get token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to marshal transaction status request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.Environment.authURL()+transactionStatusPath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to create transaction status request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("daraja: transaction status request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to read transaction status response: %w", err)
	}

	var tsResp TransactionStatusResponse
	if err := json.Unmarshal(respBody, &tsResp); err != nil {
		return nil, fmt.Errorf("daraja: failed to unmarshal transaction status response: %w", err)
	}

	if tsResp.ResponseCode != "0" {
		return &tsResp, &DarajaError{
			HTTPStatus:   resp.StatusCode,
			ResponseCode: tsResp.ResponseCode,
			ResponseDesc: tsResp.ResponseDescription,
			Endpoint:     transactionStatusPath,
			Retryable:    isRetryable(resp.StatusCode),
		}
	}

	return &tsResp, nil
}

func validateTransactionStatusRequest(req *TransactionStatusRequest) error {
	if req.Initiator == "" {
		return fmt.Errorf("daraja: transaction status Initiator is required")
	}
	if req.SecurityCredential == "" {
		return fmt.Errorf("daraja: transaction status SecurityCredential is required")
	}
	if req.CommandID == "" {
		return fmt.Errorf("daraja: transaction status CommandID is required")
	}
	if req.TransactionID == "" {
		return fmt.Errorf("daraja: transaction status TransactionID is required")
	}
	if req.PartyA == "" {
		return fmt.Errorf("daraja: transaction status PartyA is required")
	}
	if req.IdentifierType == "" {
		return fmt.Errorf("daraja: transaction status IdentifierType is required")
	}
	if req.ResultURL == "" {
		return fmt.Errorf("daraja: transaction status ResultURL is required")
	}
	if req.QueueTimeOutURL == "" {
		return fmt.Errorf("daraja: transaction status QueueTimeOutURL is required")
	}
	if req.Remarks == "" {
		return fmt.Errorf("daraja: transaction status Remarks is required")
	}
	return nil
}
