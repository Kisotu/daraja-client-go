package daraja

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// STKPushRequest represents a Lipa Na M-Pesa Online STK Push request.
type STKPushRequest struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            string `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

// STKPushResponse represents the response from an STK Push request.
type STKPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

func (c *client) STKPush(ctx context.Context, req STKPushRequest) (*STKPushResponse, error) {
	if err := validateSTKPushRequest(&req); err != nil {
		return nil, err
	}

	token, err := c.authMgr.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to get token: %w", err)
	}

	now := time.Now().UTC()
	req.Timestamp = now.Format("20060102150405")
	req.Password = base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s%s%s", req.BusinessShortCode, c.config.Passkey, req.Timestamp)),
	)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to marshal STK push request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.Environment.authURL()+"/mpesa/stkpush/v1/processrequest", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to create STK push request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("daraja: STK push request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("daraja: failed to read STK push response: %w", err)
	}

	var stkResp STKPushResponse
	if err := json.Unmarshal(respBody, &stkResp); err != nil {
		return nil, fmt.Errorf("daraja: failed to unmarshal STK push response: %w", err)
	}

	if stkResp.ResponseCode != "0" {
		return &stkResp, &DarajaError{
			HTTPStatus:    resp.StatusCode,
			ResponseCode:  stkResp.ResponseCode,
			ResponseDesc:  stkResp.ResponseDescription,
			Endpoint:      "/mpesa/stkpush/v1/processrequest",
			Retryable:     isRetryable(resp.StatusCode),
		}
	}

	return &stkResp, nil
}

func validateSTKPushRequest(req *STKPushRequest) error {
	if req.BusinessShortCode == "" {
		return fmt.Errorf("daraja: STK push BusinessShortCode is required")
	}
	if req.Amount == "" {
		return fmt.Errorf("daraja: STK push Amount is required")
	}
	if req.PartyA == "" {
		return fmt.Errorf("daraja: STK push PartyA is required")
	}
	if req.PartyB == "" {
		return fmt.Errorf("daraja: STK push PartyB is required")
	}
	if req.PhoneNumber == "" {
		return fmt.Errorf("daraja: STK push PhoneNumber is required")
	}
	if req.CallBackURL == "" {
		return fmt.Errorf("daraja: STK push CallBackURL is required")
	}
	if req.AccountReference == "" {
		return fmt.Errorf("daraja: STK push AccountReference is required")
	}
	if req.TransactionDesc == "" {
		return fmt.Errorf("daraja: STK push TransactionDesc is required")
	}
	return nil
}

func isRetryable(statusCode int) bool {
	return statusCode == 429 || statusCode >= 500
}