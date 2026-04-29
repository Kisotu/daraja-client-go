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
	c2bSimulatePath = "/mpesa/c2b/v1/simulate"
	c2bRegURLPath   = "/mpesa/c2b/v1/registerurl"
)

// C2BSimulationRequest represents a C2B simulation request.
type C2BSimulationRequest struct {
	ShortCode     string `json:"ShortCode"`
	CommandID     string `json:"CommandID"`
	Amount        string `json:"Amount"`
	Msisdn        string `json:"Msisdn"`
	BillRefNumber string `json:"BillRefNumber"`
}

// C2BResponse represents the generic response from C2B endpoints.
type C2BResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

// RegisterURLRequest represents the request to register C2B confirmation and
// validation URLs.
type RegisterURLRequest struct {
	ShortCode       string `json:"ShortCode"`
	ResponseType    string `json:"ResponseType"`
	ConfirmationURL string `json:"ConfirmationURL"`
	ValidationURL   string `json:"ValidationURL"`
}

// RegisterURLResponse carries the result of a URL registration.
type RegisterURLResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

func (c *client) C2BSimulate(ctx context.Context, req C2BSimulationRequest) (*C2BResponse, error) {
	if err := validateC2BSimulationRequest(&req); err != nil {
		return nil, err
	}

	token, err := c.authMgr.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to get token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to marshal C2B simulation request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.Environment.authURL()+c2bSimulatePath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to create C2B simulation request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("daraja: C2B simulation request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to read C2B simulation response: %w", err)
	}

	var c2bResp C2BResponse
	if err := json.Unmarshal(respBody, &c2bResp); err != nil {
		return nil, fmt.Errorf("daraja: failed to unmarshal C2B simulation response: %w", err)
	}

	if c2bResp.ResponseCode != "0" {
		return &c2bResp, &DarajaError{
			HTTPStatus:   resp.StatusCode,
			ResponseCode: c2bResp.ResponseCode,
			ResponseDesc: c2bResp.ResponseDescription,
			Endpoint:     c2bSimulatePath,
			Retryable:    isRetryable(resp.StatusCode),
		}
	}

	return &c2bResp, nil
}

func (c *client) C2BRegisterURL(ctx context.Context, req RegisterURLRequest) (*RegisterURLResponse, error) {
	if err := validateRegisterURLRequest(&req); err != nil {
		return nil, err
	}

	token, err := c.authMgr.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to get token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to marshal register URL request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.Environment.authURL()+c2bRegURLPath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to create register URL request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("daraja: register URL request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to read register URL response: %w", err)
	}

	var regResp RegisterURLResponse
	if err := json.Unmarshal(respBody, &regResp); err != nil {
		return nil, fmt.Errorf("daraja: failed to unmarshal register URL response: %w", err)
	}

	if regResp.ResponseCode != "0" {
		return &regResp, &DarajaError{
			HTTPStatus:   resp.StatusCode,
			ResponseCode: regResp.ResponseCode,
			ResponseDesc: regResp.ResponseDescription,
			Endpoint:     c2bRegURLPath,
			Retryable:    isRetryable(resp.StatusCode),
		}
	}

	return &regResp, nil
}

func validateC2BSimulationRequest(req *C2BSimulationRequest) error {
	if req.ShortCode == "" {
		return fmt.Errorf("daraja: C2B ShortCode is required")
	}
	if req.CommandID == "" {
		return fmt.Errorf("daraja: C2B CommandID is required")
	}
	if req.Amount == "" {
		return fmt.Errorf("daraja: C2B Amount is required")
	}
	if req.Msisdn == "" {
		return fmt.Errorf("daraja: C2B Msisdn is required")
	}
	if req.BillRefNumber == "" {
		return fmt.Errorf("daraja: C2B BillRefNumber is required")
	}
	return nil
}

func validateRegisterURLRequest(req *RegisterURLRequest) error {
	if req.ShortCode == "" {
		return fmt.Errorf("daraja: register ShortCode is required")
	}
	if req.ResponseType == "" {
		return fmt.Errorf("daraja: register ResponseType is required")
	}
	if req.ConfirmationURL == "" {
		return fmt.Errorf("daraja: register ConfirmationURL is required")
	}
	if req.ValidationURL == "" {
		return fmt.Errorf("daraja: register ValidationURL is required")
	}
	return nil
}
