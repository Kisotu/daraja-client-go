package daraja

import (
	"testing"
)

func TestValidateC2BSimulationRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     C2BSimulationRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: C2BSimulationRequest{
				ShortCode:     "174379",
				CommandID:     "CustomerPayBillOnline",
				Amount:        "100",
				Msisdn:        "254712345678",
				BillRefNumber: "INV-001",
			},
			wantErr: false,
		},
		{
			name:    "missing ShortCode",
			req:     buildC2BSimReq(func(r *C2BSimulationRequest) { r.ShortCode = "" }),
			wantErr: true,
		},
		{
			name:    "missing CommandID",
			req:     buildC2BSimReq(func(r *C2BSimulationRequest) { r.CommandID = "" }),
			wantErr: true,
		},
		{
			name:    "missing Amount",
			req:     buildC2BSimReq(func(r *C2BSimulationRequest) { r.Amount = "" }),
			wantErr: true,
		},
		{
			name:    "missing Msisdn",
			req:     buildC2BSimReq(func(r *C2BSimulationRequest) { r.Msisdn = "" }),
			wantErr: true,
		},
		{
			name:    "missing BillRefNumber",
			req:     buildC2BSimReq(func(r *C2BSimulationRequest) { r.BillRefNumber = "" }),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateC2BSimulationRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateC2BSimulationRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRegisterURLRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     RegisterURLRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: RegisterURLRequest{
				ShortCode:       "174379",
				ResponseType:    "Completed",
				ConfirmationURL: "https://example.com/confirm",
				ValidationURL:   "https://example.com/validate",
			},
			wantErr: false,
		},
		{
			name:    "missing ShortCode",
			req:     buildRegURLReq(func(r *RegisterURLRequest) { r.ShortCode = "" }),
			wantErr: true,
		},
		{
			name:    "missing ResponseType",
			req:     buildRegURLReq(func(r *RegisterURLRequest) { r.ResponseType = "" }),
			wantErr: true,
		},
		{
			name:    "missing ConfirmationURL",
			req:     buildRegURLReq(func(r *RegisterURLRequest) { r.ConfirmationURL = "" }),
			wantErr: true,
		},
		{
			name:    "missing ValidationURL",
			req:     buildRegURLReq(func(r *RegisterURLRequest) { r.ValidationURL = "" }),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegisterURLRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRegisterURLRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func buildC2BSimReq(mod func(*C2BSimulationRequest)) C2BSimulationRequest {
	r := C2BSimulationRequest{
		ShortCode:     "174379",
		CommandID:     "CustomerPayBillOnline",
		Amount:        "100",
		Msisdn:        "254712345678",
		BillRefNumber: "INV-001",
	}
	mod(&r)
	return r
}

func buildRegURLReq(mod func(*RegisterURLRequest)) RegisterURLRequest {
	r := RegisterURLRequest{
		ShortCode:       "174379",
		ResponseType:    "Completed",
		ConfirmationURL: "https://example.com/confirm",
		ValidationURL:   "https://example.com/validate",
	}
	mod(&r)
	return r
}
