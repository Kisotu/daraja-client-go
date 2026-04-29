package daraja

import (
	"testing"
)

func TestValidateTransactionStatusRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     TransactionStatusRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: TransactionStatusRequest{
				Initiator:          "TestInit",
				SecurityCredential: "cred123",
				CommandID:          "TransactionStatusQuery",
				TransactionID:      "LHG7C0J0G",
				PartyA:             "600000",
				IdentifierType:     "1",
				ResultURL:          "https://example.com/result",
				QueueTimeOutURL:    "https://example.com/timeout",
				Remarks:            "Test status check",
				Occasion:           "Test",
			},
			wantErr: false,
		},
		{
			name:    "missing Initiator",
			req:     buildTxStatusReq(func(r *TransactionStatusRequest) { r.Initiator = "" }),
			wantErr: true,
		},
		{
			name:    "missing SecurityCredential",
			req:     buildTxStatusReq(func(r *TransactionStatusRequest) { r.SecurityCredential = "" }),
			wantErr: true,
		},
		{
			name:    "missing CommandID",
			req:     buildTxStatusReq(func(r *TransactionStatusRequest) { r.CommandID = "" }),
			wantErr: true,
		},
		{
			name:    "missing TransactionID",
			req:     buildTxStatusReq(func(r *TransactionStatusRequest) { r.TransactionID = "" }),
			wantErr: true,
		},
		{
			name:    "missing PartyA",
			req:     buildTxStatusReq(func(r *TransactionStatusRequest) { r.PartyA = "" }),
			wantErr: true,
		},
		{
			name:    "missing IdentifierType",
			req:     buildTxStatusReq(func(r *TransactionStatusRequest) { r.IdentifierType = "" }),
			wantErr: true,
		},
		{
			name:    "missing ResultURL",
			req:     buildTxStatusReq(func(r *TransactionStatusRequest) { r.ResultURL = "" }),
			wantErr: true,
		},
		{
			name:    "missing QueueTimeOutURL",
			req:     buildTxStatusReq(func(r *TransactionStatusRequest) { r.QueueTimeOutURL = "" }),
			wantErr: true,
		},
		{
			name:    "missing Remarks",
			req:     buildTxStatusReq(func(r *TransactionStatusRequest) { r.Remarks = "" }),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransactionStatusRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTransactionStatusRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func buildTxStatusReq(mod func(*TransactionStatusRequest)) TransactionStatusRequest {
	r := TransactionStatusRequest{
		Initiator:          "TestInit",
		SecurityCredential: "cred123",
		CommandID:          "TransactionStatusQuery",
		TransactionID:      "LHG7C0J0G",
		PartyA:             "600000",
		IdentifierType:     "1",
		ResultURL:          "https://example.com/result",
		QueueTimeOutURL:    "https://example.com/timeout",
		Remarks:            "Test status check",
		Occasion:           "Test",
	}
	mod(&r)
	return r
}
