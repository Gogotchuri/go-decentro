package decentro

import (
	"fmt"
	"github.com/gogotchuri/go-decentro/decentroModels"
	"net/http"
)

type CollectionsAPI interface {
	CheckPaymentStatus(transactionID string) (*decentroModels.PaymentStatusResponseData, error)
	GenerateUPIPaymentLink(request decentroModels.PaymentLinkRequest) (*decentroModels.PaymentLinkResponseData, error)
}

func (c client) CheckPaymentStatus(transactionID string) (*decentroModels.PaymentStatusResponseData, error) {
	resp := decentroModels.PaymentStatusResponse{}
	if err := c.send(http.MethodGet, fmt.Sprintf(CheckPaymentStatusEndpoint, transactionID), nil, &resp); err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, &decentroModels.ErrorResponse{BaseResponse: resp.BaseResponse}
	}
	return &resp.Data, nil
}

func (c client) GenerateUPIPaymentLink(request decentroModels.PaymentLinkRequest) (*decentroModels.PaymentLinkResponseData, error) {
	if request.PayeeAccount == "" {
		request.PayeeAccount = c.defaultPayeeNumber
	}
	resp := decentroModels.PaymentLinkResponse{}
	if err := c.send(http.MethodPost, GenerateUPIPaymentURLEndpoint, request, &resp); err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, &decentroModels.ErrorResponse{BaseResponse: resp.BaseResponse}
	}
	return &resp.Data, nil
}
