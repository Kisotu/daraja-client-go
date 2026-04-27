package daraja

import (
	"testing"
)

func TestValidateB2CPaymentRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     B2CPaymentRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: B2CPaymentRequest{
				OriginatorConversationID: "occ-123",
				InitiatorName:            "TestInit",
				SecurityCredential:       "cred123",
				CommandID:                "BusinessPayment",
				Amount:                   "100",
				PartyA:                   "600000",
				PartyB:                   "254712345678",
				Remarks:                  "Test payment",
				QueueTimeOutURL:          "https://example.com/timeout",
				ResultURL:                "https://example.com/result",
			},
			wantErr: false,
		},
		{
			name:    "missing InitiatorName",
			req:     buildB2CReq(func(r *B2CPaymentRequest) { r.InitiatorName = "" }),
			wantErr: true,
		},
		{
			name:    "missing SecurityCredential",
			req:     buildB2CReq(func(r *B2CPaymentRequest) { r.SecurityCredential = "" }),
			wantErr: true,
		},
		{
			name:    "missing CommandID",
			req:     buildB2CReq(func(r *B2CPaymentRequest) { r.CommandID = "" }),
			wantErr: true,
		},
		{
			name:    "missing Amount",
			req:     buildB2CReq(func(r *B2CPaymentRequest) { r.Amount = "" }),
			wantErr: true,
		},
		{
			name:    "missing PartyA",
			req:     buildB2CReq(func(r *B2CPaymentRequest) { r.PartyA = "" }),
			wantErr: true,
		},
		{
			name:    "missing PartyB",
			req:     buildB2CReq(func(r *B2CPaymentRequest) { r.PartyB = "" }),
			wantErr: true,
		},
		{
			name:    "missing Remarks",
			req:     buildB2CReq(func(r *B2CPaymentRequest) { r.Remarks = "" }),
			wantErr: true,
		},
		{
			name:    "missing QueueTimeOutURL",
			req:     buildB2CReq(func(r *B2CPaymentRequest) { r.QueueTimeOutURL = "" }),
			wantErr: true,
		},
		{
			name:    "missing ResultURL",
			req:     buildB2CReq(func(r *B2CPaymentRequest) { r.ResultURL = "" }),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateB2CPaymentRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateB2CPaymentRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func buildB2CReq(mod func(*B2CPaymentRequest)) B2CPaymentRequest {
	r := B2CPaymentRequest{
		OriginatorConversationID: "occ-123",
		InitiatorName:            "TestInit",
		SecurityCredential:       "cred123",
		CommandID:                "BusinessPayment",
		Amount:                   "100",
		PartyA:                   "600000",
		PartyB:                   "254712345678",
		Remarks:                  "Test payment",
		QueueTimeOutURL:          "https://example.com/timeout",
		ResultURL:                "https://example.com/result",
	}
	mod(&r)
	return r
}