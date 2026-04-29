package daraja

import (
	"testing"
)

func TestValidateSTKPushRequest_MissingFields(t *testing.T) {
	tests := []struct {
		name    string
		req     STKPushRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: STKPushRequest{
				BusinessShortCode: "174379",
				Amount:            "1",
				PartyA:            "254712345678",
				PartyB:            "174379",
				PhoneNumber:       "254712345678",
				CallBackURL:       "https://example.com/callback",
				AccountReference:  "Test",
				TransactionDesc:   "Test payment",
			},
			wantErr: false,
		},
		{
			name:    "missing BusinessShortCode",
			req:     buildReq(func(r *STKPushRequest) { r.BusinessShortCode = "" }),
			wantErr: true,
		},
		{
			name:    "missing Amount",
			req:     buildReq(func(r *STKPushRequest) { r.Amount = "" }),
			wantErr: true,
		},
		{
			name:    "missing PartyA",
			req:     buildReq(func(r *STKPushRequest) { r.PartyA = "" }),
			wantErr: true,
		},
		{
			name:    "missing PartyB",
			req:     buildReq(func(r *STKPushRequest) { r.PartyB = "" }),
			wantErr: true,
		},
		{
			name:    "missing PhoneNumber",
			req:     buildReq(func(r *STKPushRequest) { r.PhoneNumber = "" }),
			wantErr: true,
		},
		{
			name:    "missing CallBackURL",
			req:     buildReq(func(r *STKPushRequest) { r.CallBackURL = "" }),
			wantErr: true,
		},
		{
			name:    "missing AccountReference",
			req:     buildReq(func(r *STKPushRequest) { r.AccountReference = "" }),
			wantErr: true,
		},
		{
			name:    "missing TransactionDesc",
			req:     buildReq(func(r *STKPushRequest) { r.TransactionDesc = "" }),
			wantErr: true,
		},
		{
			name:    "unsupported phone prefix",
			req:     buildReq(func(r *STKPushRequest) { r.PhoneNumber = "0712345678" }),
			wantErr: true,
		},
		{
			name:    "invalid TransactionType",
			req:     buildReq(func(r *STKPushRequest) { r.TransactionType = "InvalidType" }),
			wantErr: true,
		},
		{
			name: "valid CustomerPayBillOnline",
			req: STKPushRequest{
				BusinessShortCode: "174379",
				TransactionType:   "CustomerPayBillOnline",
				Amount:            "1",
				PartyA:            "254712345678",
				PartyB:            "174379",
				PhoneNumber:       "254712345678",
				CallBackURL:       "https://example.com/callback",
				AccountReference:  "Test",
				TransactionDesc:   "Test payment",
			},
			wantErr: false,
		},
		{
			name: "valid CustomerBuyGoodsOnline",
			req: STKPushRequest{
				BusinessShortCode: "174379",
				TransactionType:   "CustomerBuyGoodsOnline",
				Amount:            "1",
				PartyA:            "254712345678",
				PartyB:            "174379",
				PhoneNumber:       "254712345678",
				CallBackURL:       "https://example.com/callback",
				AccountReference:  "Test",
				TransactionDesc:   "Test payment",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSTKPushRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSTKPushRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		code      int
		retryable bool
	}{
		{200, false},
		{400, false},
		{401, false},
		{403, false},
		{404, false},
		{429, true},
		{500, true},
		{502, true},
		{503, true},
		{504, true},
	}
	for _, tt := range tests {
		if got := isRetryable(tt.code); got != tt.retryable {
			t.Errorf("isRetryable(%d) = %v, want %v", tt.code, got, tt.retryable)
		}
	}
}

func TestHasSupportedPrefix(t *testing.T) {
	tests := []struct {
		phone    string
		expected bool
	}{
		{"254712345678", true},
		{"254700000000", true},
		{"0712345678", false},
		{"+254712345678", false},
		{"", false},
		{"123456789", false},
	}
	for _, tt := range tests {
		if got := hasSupportedPrefix(tt.phone); got != tt.expected {
			t.Errorf("hasSupportedPrefix(%q) = %v, want %v", tt.phone, got, tt.expected)
		}
	}
}

func buildReq(mod func(*STKPushRequest)) STKPushRequest {
	r := STKPushRequest{
		BusinessShortCode: "174379",
		Amount:            "1",
		PartyA:            "254712345678",
		PartyB:            "174379",
		PhoneNumber:       "254712345678",
		CallBackURL:       "https://example.com/callback",
		AccountReference:  "Test",
		TransactionDesc:   "Test payment",
	}
	mod(&r)
	return r
}
