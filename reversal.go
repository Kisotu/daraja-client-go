package daraja

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const (
	reversalPath = "/mpesa/reversal/v1/request"
)

// ReversalRequest represents a request to reverse a transaction.
type ReversalRequest struct {
	Initiator              string `json:"Initiator"`
	SecurityCredential     string `json:"SecurityCredential"`
	CommandID              string `json:"CommandID"`
	TransactionID          string `json:"TransactionID"`
	OriginalConversationID string `json:"OriginalConversationID"`
	PartyA                 string `json:"PartyA"`
	IdentifierType         string `json:"IdentifierType"`
	ResultURL              string `json:"ResultURL"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL"`
	Remarks                string `json:"Remarks"`
	Occasion               string `json:"Occasion"`
}

// ReversalResponse represents the response from a reversal request.
type ReversalResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

func (c *client) Reversal(ctx context.Context, req ReversalRequest) (*ReversalResponse, error) {
	if req.OriginalConversationID == "" {
		req.OriginalConversationID = uuid.NewString()
	}

	if err := validateReversalRequest(&req); err != nil {
		return nil, err
	}

	token, err := c.authMgr.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to get token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to marshal reversal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.Environment.authURL()+reversalPath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to create reversal request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("daraja: reversal request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to read reversal response: %w", err)
	}

	var revResp ReversalResponse
	if err := json.Unmarshal(respBody, &revResp); err != nil {
		return nil, fmt.Errorf("daraja: failed to unmarshal reversal response: %w", err)
	}

	if revResp.ResponseCode != "0" {
		return &revResp, &DarajaError{
			HTTPStatus:   resp.StatusCode,
			ResponseCode: revResp.ResponseCode,
			ResponseDesc: revResp.ResponseDescription,
			Endpoint:     reversalPath,
			Retryable:    isRetryable(resp.StatusCode),
		}
	}

	return &revResp, nil
}

func validateReversalRequest(req *ReversalRequest) error {
	if req.Initiator == "" {
		return fmt.Errorf("daraja: reversal Initiator is required")
	}
	if req.SecurityCredential == "" {
		return fmt.Errorf("daraja: reversal SecurityCredential is required")
	}
	if req.CommandID == "" {
		return fmt.Errorf("daraja: reversal CommandID is required")
	}
	if req.TransactionID == "" {
		return fmt.Errorf("daraja: reversal TransactionID is required")
	}
	if req.PartyA == "" {
		return fmt.Errorf("daraja: reversal PartyA is required")
	}
	if req.IdentifierType == "" {
		return fmt.Errorf("daraja: reversal IdentifierType is required")
	}
	if req.ResultURL == "" {
		return fmt.Errorf("daraja: reversal ResultURL is required")
	}
	if req.QueueTimeOutURL == "" {
		return fmt.Errorf("daraja: reversal QueueTimeOutURL is required")
	}
	if req.Remarks == "" {
		return fmt.Errorf("daraja: reversal Remarks is required")
	}
	return nil
}
