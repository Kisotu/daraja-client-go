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
	b2cPaymentPath = "/mpesa/b2c/v1/paymentrequest"
)

// B2CPaymentRequest represents a B2C payment request.
type B2CPaymentRequest struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	InitiatorName            string `json:"InitiatorName"`
	SecurityCredential       string `json:"SecurityCredential"`
	CommandID                string `json:"CommandID"`
	Amount                   string `json:"Amount"`
	PartyA                   string `json:"PartyA"`
	PartyB                   string `json:"PartyB"`
	Remarks                  string `json:"Remarks"`
	QueueTimeOutURL          string `json:"QueueTimeOutURL"`
	ResultURL                string `json:"ResultURL"`
	Occasion                 string `json:"Occasion"`
}

// B2CPaymentResponse represents the response from a B2C payment request.
type B2CPaymentResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

func (c *client) B2CPayment(ctx context.Context, req B2CPaymentRequest) (*B2CPaymentResponse, error) {
	if req.OriginatorConversationID == "" {
		req.OriginatorConversationID = uuid.NewString()
	}

	if err := validateB2CPaymentRequest(&req); err != nil {
		return nil, err
	}

	token, err := c.authMgr.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to get token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to marshal B2C payment request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.Environment.authURL()+b2cPaymentPath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to create B2C payment request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("daraja: B2C payment request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to read B2C payment response: %w", err)
	}

	var b2cResp B2CPaymentResponse
	if err := json.Unmarshal(respBody, &b2cResp); err != nil {
		return nil, fmt.Errorf("daraja: failed to unmarshal B2C payment response: %w", err)
	}

	if b2cResp.ResponseCode != "0" {
		return &b2cResp, &DarajaError{
			HTTPStatus:    resp.StatusCode,
			ResponseCode:  b2cResp.ResponseCode,
			ResponseDesc:  b2cResp.ResponseDescription,
			Endpoint:      b2cPaymentPath,
			Retryable:     isRetryable(resp.StatusCode),
		}
	}

	return &b2cResp, nil
}

func validateB2CPaymentRequest(req *B2CPaymentRequest) error {
	if req.InitiatorName == "" {
		return fmt.Errorf("daraja: B2C InitiatorName is required")
	}
	if req.SecurityCredential == "" {
		return fmt.Errorf("daraja: B2C SecurityCredential is required")
	}
	if req.CommandID == "" {
		return fmt.Errorf("daraja: B2C CommandID is required")
	}
	if req.Amount == "" {
		return fmt.Errorf("daraja: B2C Amount is required")
	}
	if req.PartyA == "" {
		return fmt.Errorf("daraja: B2C PartyA is required")
	}
	if req.PartyB == "" {
		return fmt.Errorf("daraja: B2C PartyB is required")
	}
	if req.Remarks == "" {
		return fmt.Errorf("daraja: B2C Remarks is required")
	}
	if req.QueueTimeOutURL == "" {
		return fmt.Errorf("daraja: B2C QueueTimeOutURL is required")
	}
	if req.ResultURL == "" {
		return fmt.Errorf("daraja: B2C ResultURL is required")
	}
	return nil
}