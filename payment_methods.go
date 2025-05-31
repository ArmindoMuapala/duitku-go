package duitku

import (
	"fmt"
	"time"
)

// PaymentMethod represents a payment method available in Duitku
type PaymentMethod struct {
	PaymentMethod string `json:"paymentMethod"`
	PaymentName   string `json:"paymentName"`
	PaymentImage  string `json:"paymentImage"`
	TotalFee      string `json:"totalFee"`
}

// PaymentMethodResponse represents the response from the get payment methods endpoint
type PaymentMethodResponse struct {
	PaymentFee      []PaymentMethod `json:"paymentFee"`
	ResponseCode    string          `json:"responseCode"`
	ResponseMessage string          `json:"responseMessage"`
}

// GetPaymentMethodsRequest represents the request to get payment methods
type GetPaymentMethodsRequest struct {
	MerchantCode string `json:"merchantcode"`
	Amount       int    `json:"amount"`
	DateTime     string `json:"datetime"`
	Signature    string `json:"signature"`
}

// GetPaymentMethods retrieves the available payment methods for the specified amount
func (c *Client) GetPaymentMethods(amount int) ([]PaymentMethod, error) {
	datetime := time.Now().Format("2006-01-02 15:04:05")
	signature := c.createSignatureSHA256(c.config.MerchantCode, fmt.Sprintf("%d", amount), datetime)

	request := GetPaymentMethodsRequest{
		MerchantCode: c.config.MerchantCode,
		Amount:       amount,
		DateTime:     datetime,
		Signature:    signature,
	}

	var response PaymentMethodResponse
	err := c.doRequest("POST", "merchant/paymentmethod/getpaymentmethod", request, &response)
	if err != nil {
		return nil, err
	}

	if response.ResponseCode != "00" {
		return nil, fmt.Errorf("error getting payment methods: %s", response.ResponseMessage)
	}

	return response.PaymentFee, nil
}
