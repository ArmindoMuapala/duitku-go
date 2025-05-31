package duitku

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// CallbackData represents the data received in a Duitku callback
type CallbackData struct {
	MerchantCode    string `json:"merchantCode"`
	Amount          int    `json:"amount"`
	MerchantOrderID string `json:"merchantOrderId"`
	ProductDetail   string `json:"productDetail"`
	AdditionalParam string `json:"additionalParam"`
	PaymentCode     string `json:"paymentCode"`
	ResultCode      string `json:"resultCode"`
	Reference       string `json:"reference"`
	Signature       string `json:"signature"`
}

// ParseCallback parses the callback data from an HTTP request
func (c *Client) ParseCallback(r *http.Request) (*CallbackData, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("error parsing form: %w", err)
	}

	// Extract callback data from form
	merchantCode := r.FormValue("merchantCode")
	amountStr := r.FormValue("amount")
	merchantOrderID := r.FormValue("merchantOrderId")
	productDetail := r.FormValue("productDetail")
	additionalParam := r.FormValue("additionalParam")
	paymentCode := r.FormValue("paymentCode")
	resultCode := r.FormValue("resultCode")
	reference := r.FormValue("reference")
	signature := r.FormValue("signature")

	// Validate required fields
	if merchantCode == "" || amountStr == "" || merchantOrderID == "" || resultCode == "" || signature == "" {
		return nil, errors.New("missing required callback parameters")
	}

	// Parse amount
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	// Create callback data
	callbackData := &CallbackData{
		MerchantCode:    merchantCode,
		Amount:          amount,
		MerchantOrderID: merchantOrderID,
		ProductDetail:   productDetail,
		AdditionalParam: additionalParam,
		PaymentCode:     paymentCode,
		ResultCode:      resultCode,
		Reference:       reference,
		Signature:       signature,
	}

	// Verify signature
	if !c.VerifyCallbackSignature(callbackData) {
		return nil, errors.New("invalid callback signature")
	}

	return callbackData, nil
}

// VerifyCallbackSignature verifies the signature of a callback
func (c *Client) VerifyCallbackSignature(data *CallbackData) bool {
	// Create signature string
	signatureStr := fmt.Sprintf("%s%d%s%s", 
		data.MerchantCode, 
		data.Amount, 
		data.MerchantOrderID, 
		c.config.APIKey,
	)

	// Calculate MD5 hash
	hash := md5.Sum([]byte(signatureStr))
	expectedSignature := hex.EncodeToString(hash[:])

	// Compare signatures (case-insensitive)
	return strings.EqualFold(expectedSignature, data.Signature)
}

// IsSuccessful returns true if the callback indicates a successful payment
func (data *CallbackData) IsSuccessful() bool {
	return data.ResultCode == "00"
}

// HandleCallback is a helper function to handle Duitku callbacks
func (c *Client) HandleCallback(w http.ResponseWriter, r *http.Request, handler func(*CallbackData) error) {
	// Parse callback data
	callbackData, err := c.ParseCallback(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call handler
	if err := handler(callbackData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
