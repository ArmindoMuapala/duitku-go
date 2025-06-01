package duitku

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// CallbackData represents the data received in a Duitku callback
//
// Parameter descriptions:
// - merchantCode: Merchant code, sent by the Duitku server to inform which project code is on used. Example: DXXXX
// - amount: Transaction amount. Example: 150000
// - merchantOrderId: Transaction number from merchant. Example: abcde12345
// - productDetail: Description about product/service on transaction. Example: Payment example for example merchant
// - additionalParam: Additional parameters that you send at the beginning of the transaction request.
// - paymentCode: Payment method code. Example: VC
// - resultCode: Result code callback notification. Example: 00 - Success, 01 - Failed
// - merchantUserId: Customer's username or email on your site. Example: your_customer@example.com
// - reference: Transaction reference number from Duitku. Please keep it for the purposes of recording or tracking transactions. Example: DXXXXCX80TXXX5Q70QCI
// - signature: Transaction identification code. Contains transaction parameters which are hashed using the MD5 hashing method. Security parameters as a reference that the request received comes from the Duitku server. Formula: MD5(merchantcode + amount + merchantOrderId + apiKey). Example: 506f88f1000dfb4a6541ff94d9b8d1e6
// - publisherOrderId: Unique transaction payment number from Duitku. Please keep it for the purposes of recording or tracking transactions. Example: MGUHWKJX3M1KMSQN5
// - spUserHash: Will be sent to your callback if the payment method using ShopeePay(QRIS, App, and Account Link). If this string parameter contains alphabet and numeric, then it might been paid by Shopee itself. Example: xxxyyyzzz
// - settlementDate: Settlement date estimation information. Format: YYYY-MM-DD. Example: 2023-07-25
// - issuerCode: QRIS issuer code information. Example: 93600523
type CallbackData struct {
	MerchantCode     string `json:"merchantCode"`
	Amount           string `json:"amount"`
	MerchantOrderID  string `json:"merchantOrderId"`
	ProductDetail    string `json:"productDetail"`
	AdditionalParam  string `json:"additionalParam"`
	PaymentCode      string `json:"paymentCode"`
	ResultCode       string `json:"resultCode"`
	MerchantUserID   string `json:"merchantUserId"`
	Reference        string `json:"reference"`
	Signature        string `json:"signature"`
	PublisherOrderId string `json:"publisherOrderId"`
	SpUserHash       string `json:"spUserHash"`
	SettlementDate   string `json:"settlementDate"`
	IssuerCode       string `json:"issuerCode"`
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
	merchantUserID := r.FormValue("merchantUserId")
	reference := r.FormValue("reference")
	signature := r.FormValue("signature")
	publisherOrderId := r.FormValue("publisherOrderId")
	spUserHash := r.FormValue("spUserHash")
	settlementDate := r.FormValue("settlementDate")
	issuerCode := r.FormValue("issuerCode")

	// Validate required fields
	if merchantCode == "" || amountStr == "" || merchantOrderID == "" || resultCode == "" || signature == "" {
		return nil, errors.New("missing required callback parameters")
	}

	// Create callback data
	callbackData := &CallbackData{
		MerchantCode:     merchantCode,
		Amount:           amountStr,
		MerchantOrderID:  merchantOrderID,
		ProductDetail:    productDetail,
		AdditionalParam:  additionalParam,
		PaymentCode:      paymentCode,
		ResultCode:       resultCode,
		MerchantUserID:   merchantUserID,
		Reference:        reference,
		Signature:        signature,
		PublisherOrderId: publisherOrderId,
		SpUserHash:       spUserHash,
		SettlementDate:   settlementDate,
		IssuerCode:       issuerCode,
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
	signatureStr := fmt.Sprintf("%s%s%s%s",
		data.MerchantCode,
		data.Amount,
		data.MerchantOrderID,
		c.config.APIKey,
	)

	// Calculate MD5 hash
	hash := md5.Sum([]byte(signatureStr))
	expectedSignature := hex.EncodeToString(hash[:])

	// Compare signatures
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
