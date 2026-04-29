# Examples

## STK Push (Lipa Na M-Pesa Online)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Kisotu/daraja-client-go"
)

func main() {
	client, err := daraja.NewClient(
		daraja.WithEnvironment(daraja.Sandbox),
		daraja.WithCredentials(os.Getenv("DARAJA_CONSUMER_KEY"), os.Getenv("DARAJA_CONSUMER_SECRET")),
		daraja.WithShortCode(os.Getenv("DARAJA_SHORTCODE"), os.Getenv("DARAJA_PASSKEY")),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.STKPush(context.Background(), daraja.STKPushRequest{
		BusinessShortCode: os.Getenv("DARAJA_SHORTCODE"),
		Amount:            "1",
		PartyA:            "254712345678",
		PartyB:            os.Getenv("DARAJA_SHORTCODE"),
		PhoneNumber:       "254712345678",
		CallBackURL:       "https://example.com/callback",
		AccountReference:  "INV-001",
		TransactionDesc:   "Payment for invoice INV-001",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CheckoutRequestID: %s\n", resp.CheckoutRequestID)
}
```

## STK Push Query

```go
resp, err := client.STKPushQuery(ctx, daraja.STKPushQueryRequest{
    BusinessShortCode: os.Getenv("DARAJA_SHORTCODE"),
    CheckoutRequestID: "ws_CO_040420251234567890",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %s - %s\n", resp.ResultCode, resp.ResultDesc)
```

## C2B Simulation

```go
resp, err := client.C2BSimulate(ctx, daraja.C2BSimulationRequest{
    ShortCode:  os.Getenv("DARAJA_SHORTCODE"),
    CommandID:  "CustomerPayBillOnline",
    Amount:     "100",
    Msisdn:     "254712345678",
    BillRef:    "INV-001",
})
```

## C2B Register URL

```go
resp, err := client.C2BRegisterURL(ctx, daraja.RegisterURLRequest{
    ShortCode:       os.Getenv("DARAJA_SHORTCODE"),
    ResponseType:    "Completed",
    ConfirmationURL: "https://example.com/confirmation",
    ValidationURL:   "https://example.com/validation",
})
```

## B2C Payment

```go
resp, err := client.B2CPayment(ctx, daraja.B2CPaymentRequest{
    InitiatorName:      "testapi",
    SecurityCredential: os.Getenv("DARAJA_SECURITY_CREDENTIAL"),
    CommandID:          "BusinessPayment",
    Amount:             "100",
    PartyA:             os.Getenv("DARAJA_SHORTCODE"),
    PartyB:             "254712345678",
    Remarks:            "Salary payment",
    QueueTimeOutURL:    "https://example.com/timeout",
    ResultURL:          "https://example.com/result",
    Occasion:           "",
})
```

## Transaction Status

```go
resp, err := client.TransactionStatus(ctx, daraja.TransactionStatusRequest{
    Initiator:          "testapi",
    SecurityCredential: os.Getenv("DARAJA_SECURITY_CREDENTIAL"),
    CommandID:          "TransactionStatusQuery",
    TransactionID:      "LHG123ABC",
    PartyA:             os.Getenv("DARAJA_SHORTCODE"),
    IdentifierType:     "1",
    ResultURL:          "https://example.com/result",
    QueueTimeOutURL:    "https://example.com/timeout",
    Remarks:            "Status check",
})
```

## Reversal

```go
resp, err := client.Reversal(ctx, daraja.ReversalRequest{
    Initiator:              "testapi",
    SecurityCredential:     os.Getenv("DARAJA_SECURITY_CREDENTIAL"),
    CommandID:              "TransactionReversal",
    TransactionID:          "LHG123ABC",
    Amount:                 "100",
    ReceiverParty:          os.Getenv("DARAJA_SHORTCODE"),
    RecieverIdentifierType: "11",
    QueueTimeOutURL:        "https://example.com/timeout",
    ResultURL:              "https://example.com/result",
    Remarks:                "Customer requested reversal",
    Occasion:               "",
})
```

## Callback Verification

```go
package main

import (
	"log"
	"net/http"

	"github.com/Kisotu/daraja-client-go/callback"
)

func main() {
	v := callback.NewVerifier()
	v.WithIPAllowlist([]string{"196.201.214.200/32"})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		payload, err := v.VerifySTKPush(r)
		if err != nil {
			http.Error(w, "invalid callback", http.StatusBadRequest)
			return
		}
		log.Printf("STK callback: %+v", payload)
		callback.WriteOK(w)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```