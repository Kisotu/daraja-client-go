package daraja

import (
	"testing"
)

func TestValidateReversalRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     ReversalRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: ReversalRequest{
				Initiator:              "TestInit",
				SecurityCredential:     "cred123",
				CommandID:              "TransactionReversal",
				TransactionID:          "LHG7C0J0G",
				OriginalConversationID: "conv-123",
				PartyA:                 "600000",
				IdentifierType:         "1",
				ResultURL:              "https://example.com/result",
				QueueTimeOutURL:        "https://example.com/timeout",
				Remarks:                "Test reversal",
			},
			wantErr: false,
		},
		{
			name:    "missing Initiator",
			req:     buildReversalReq(func(r *ReversalRequest) { r.Initiator = "" }),
			wantErr: true,
		},
		{
			name:    "missing SecurityCredential",
			req:     buildReversalReq(func(r *ReversalRequest) { r.SecurityCredential = "" }),
			wantErr: true,
		},
		{
			name:    "missing CommandID",
			req:     buildReversalReq(func(r *ReversalRequest) { r.CommandID = "" }),
			wantErr: true,
		},
		{
			name:    "missing TransactionID",
			req:     buildReversalReq(func(r *ReversalRequest) { r.TransactionID = "" }),
			wantErr: true,
		},
		{
			name:    "missing PartyA",
			req:     buildReversalReq(func(r *ReversalRequest) { r.PartyA = "" }),
			wantErr: true,
		},
		{
			name:    "missing IdentifierType",
			req:     buildReversalReq(func(r *ReversalRequest) { r.IdentifierType = "" }),
			wantErr: true,
		},
		{
			name:    "missing ResultURL",
			req:     buildReversalReq(func(r *ReversalRequest) { r.ResultURL = "" }),
			wantErr: true,
		},
		{
			name:    "missing QueueTimeOutURL",
			req:     buildReversalReq(func(r *ReversalRequest) { r.QueueTimeOutURL = "" }),
			wantErr: true,
		},
		{
			name:    "missing Remarks",
			req:     buildReversalReq(func(r *ReversalRequest) { r.Remarks = "" }),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReversalRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateReversalRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func buildReversalReq(mod func(*ReversalRequest)) ReversalRequest {
	r := ReversalRequest{
		Initiator:              "TestInit",
		SecurityCredential:     "cred123",
		CommandID:              "TransactionReversal",
		TransactionID:          "LHG7C0J0G",
		OriginalConversationID: "conv-123",
		PartyA:                 "600000",
		IdentifierType:         "1",
		ResultURL:              "https://example.com/result",
		QueueTimeOutURL:        "https://example.com/timeout",
		Remarks:                "Test reversal",
	}
	mod(&r)
	return r
}