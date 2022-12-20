package decentro

const (
	StagingAPIURL = "https://in.staging.decentro.tech"
	ProdAPIURL    = "https://in.decentro.tech"
)

const (
	GenerateUPIPaymentURLEndpoint = "/v2/payments/upi/link"
	CheckPaymentStatusEndpoint    = "/v2/payments/transaction/%s/status" //%s is the transaction id
)
