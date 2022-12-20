package decentroModels

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ResponseStatus string
type TransactionStatus string

const (
	Success ResponseStatus = "SUCCESS"
	Failure ResponseStatus = "FAILURE"
)
const (
	TransactionStatusSuccess TransactionStatus = "SUCCESS"
	TransactionStatusFailure TransactionStatus = "FAILED"
)

type BaseResponse struct {
	TransactionID string         `json:"decentroTxnId"`
	Status        ResponseStatus `json:"status"`
	ResponseCode  string         `json:"responseCode"`
	Message       string         `json:"message"`
}

func (b BaseResponse) IsSuccess() bool {
	return b.Status == Success
}

type ErrorResponse struct {
	BaseResponse
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("txID: %s, error: %s, message: %s", e.TransactionID, e.ResponseCode, e.Message)
}

// https://docs.decentro.tech/reference/payments_api-check-payment-status
// Check Payment Status
type PaymentStatusResponseData struct {
	TransactionStatus            TransactionStatus `json:"transactionStatus"`
	TransactionStatusDescription string            `json:"transactionStatusDescription"`
	BankReferenceNumber          string            `json:"bankReferenceNumber"`
	NPCITransactionID            string            `json:"npciTxnId"`
	ProviderMessage              string            `json:"providerMessage"`
}

type PaymentStatusResponse struct {
	BaseResponse
	Data PaymentStatusResponseData `json:"data"`
}

// https://docs.decentro.tech/reference/payments_api-generate-upi-payment-link
// Generate UPI Payment Link
type PaymentLinkResponseData struct {
	UPIsURI              string `json:"upiUri"`
	EncodedDynamicQrCode string `json:"encodedDynamicQrCode"`
	GeneratedLink        string `json:"generatedLink"`
	TransactionID        string `json:"transactionId"`
	TransactionStatus    string `json:"transactionStatus"`
}

type PaymentLinkResponse struct {
	BaseResponse
	Data PaymentLinkResponseData `json:"data"`
}

type PaymentLinkRequest struct {
	ReferenceId          string  `json:"reference_id"`
	PayeeAccount         string  `json:"payee_account"`
	Amount               float64 `json:"amount"`
	PurposeMessage       string  `json:"purpose_message"`
	ExpiryTimeMinutes    int32   `json:"expiry_time"`
	GenerateQr           bool    `json:"generate_qr"`
	CustomizedQrWithLogo bool    `json:"customized_qr_with_logo"`
	GenerateUPI          bool    `json:"generate_upi"`
}

func (p PaymentLinkRequest) MarshalJSON() ([]byte, error) {
	type Alias PaymentLinkRequest
	if p.ExpiryTimeMinutes < 0 || p.ExpiryTimeMinutes > 64800 {
		return nil, errors.New("expiry time must be between 0 and 64800 minutes")
	}
	if p.Amount < 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if len(p.PurposeMessage) < 5 || len(p.PurposeMessage) > 50 {
		return nil, errors.New("purpose message must be between 5 and 50 characters")
	}
	toEncode := &struct {
		*Alias
		GenerateQr           int `json:"generate_qr"`
		CustomizedQrWithLogo int `json:"customized_qr_with_logo"`
		GenerateUPI          int `json:"generate_upi"`
	}{
		Alias: (*Alias)(&p),
	}
	if p.GenerateQr {
		toEncode.GenerateQr = 1
	}
	if p.CustomizedQrWithLogo {
		toEncode.CustomizedQrWithLogo = 1
	}
	if p.GenerateUPI {
		toEncode.GenerateUPI = 1
	}
	return json.Marshal(toEncode)
}
