package duitku

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestParseCallback(t *testing.T) {
	// Create a client
	client := NewClient(Config{
		MerchantCode: "DXXXX",
		APIKey:       "DXXXXCX80TZJ85Q70QCI",
		IsSandbox:    true,
	})

	// Calculate the actual expected signature
	signatureStr := fmt.Sprintf("%s%d%s%s", "DXXXX", 40000, "ORDER123", "DXXXXCX80TZJ85Q70QCI")
	hash := md5.Sum([]byte(signatureStr))
	expectedSignature := hex.EncodeToString(hash[:])

	// Create form data
	form := url.Values{}
	form.Add("merchantCode", "DXXXX")
	form.Add("amount", "40000")
	form.Add("merchantOrderId", "ORDER123")
	form.Add("productDetail", "Test Product")
	form.Add("additionalParam", "")
	form.Add("paymentCode", "VC")
	form.Add("resultCode", "00")
	form.Add("reference", "DEV123456789")
	form.Add("signature", expectedSignature)

	// Create a request with form data
	req, err := http.NewRequest("POST", "/callback", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Parse callback
	callbackData, err := client.ParseCallback(req)
	if err != nil {
		t.Errorf("ParseCallback() error = %v, want nil", err)
	}

	// Check callback data
	if callbackData.MerchantCode != "DXXXX" {
		t.Errorf("CallbackData MerchantCode = %s, want DXXXX", callbackData.MerchantCode)
	}
	if callbackData.Amount != 40000 {
		t.Errorf("CallbackData Amount = %d, want 40000", callbackData.Amount)
	}
	if callbackData.MerchantOrderID != "ORDER123" {
		t.Errorf("CallbackData MerchantOrderID = %s, want ORDER123", callbackData.MerchantOrderID)
	}
	if callbackData.ProductDetail != "Test Product" {
		t.Errorf("CallbackData ProductDetail = %s, want Test Product", callbackData.ProductDetail)
	}
	if callbackData.PaymentCode != "VC" {
		t.Errorf("CallbackData PaymentCode = %s, want VC", callbackData.PaymentCode)
	}
	if callbackData.ResultCode != "00" {
		t.Errorf("CallbackData ResultCode = %s, want 00", callbackData.ResultCode)
	}
	if callbackData.Reference != "DEV123456789" {
		t.Errorf("CallbackData Reference = %s, want DEV123456789", callbackData.Reference)
	}
	if callbackData.Signature != expectedSignature {
		t.Errorf("CallbackData Signature = %s, want %s", callbackData.Signature, expectedSignature)
	}
}

func TestParseCallbackMissingFields(t *testing.T) {
	// Create a client
	client := NewClient(Config{
		MerchantCode: "DXXXX",
		APIKey:       "DXXXXCX80TZJ85Q70QCI",
		IsSandbox:    true,
	})

	// Create form data with missing fields
	form := url.Values{}
	form.Add("merchantCode", "DXXXX")
	// Missing amount
	form.Add("merchantOrderId", "ORDER123")
	form.Add("signature", "d5df5a9d6807a8d7fae5b76e14c6bf4a")

	// Create a request with form data
	req, err := http.NewRequest("POST", "/callback", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Parse callback
	_, err = client.ParseCallback(req)
	if err == nil {
		t.Errorf("ParseCallback() error = nil, want error")
	}
}

func TestVerifyCallbackSignature(t *testing.T) {
	// Create a client
	client := NewClient(Config{
		MerchantCode: "DXXXX",
		APIKey:       "DXXXXCX80TZJ85Q70QCI",
		IsSandbox:    true,
	})

	// Calculate the actual expected signature
	signatureStr := fmt.Sprintf("%s%d%s%s", "DXXXX", 40000, "ORDER123", "DXXXXCX80TZJ85Q70QCI")
	hash := md5.Sum([]byte(signatureStr))
	expectedSignature := hex.EncodeToString(hash[:])

	// Create a valid callback data
	callbackData := &CallbackData{
		MerchantCode:    "DXXXX",
		Amount:          40000,
		MerchantOrderID: "ORDER123",
		ResultCode:      "00",
		Signature:       expectedSignature, // Valid signature
	}

	// Verify signature
	if !client.VerifyCallbackSignature(callbackData) {
		t.Errorf("VerifyCallbackSignature() = false, want true")
	}

	// Create an invalid callback data
	invalidCallbackData := &CallbackData{
		MerchantCode:    "DXXXX",
		Amount:          40000,
		MerchantOrderID: "ORDER123",
		ResultCode:      "00",
		Signature:       "invalid", // Invalid signature
	}

	// Verify signature
	if client.VerifyCallbackSignature(invalidCallbackData) {
		t.Errorf("VerifyCallbackSignature() = true, want false")
	}
}

func TestIsSuccessful(t *testing.T) {
	// Create a successful callback data
	successfulCallback := &CallbackData{
		ResultCode: "00",
	}

	// Check if successful
	if !successfulCallback.IsSuccessful() {
		t.Errorf("IsSuccessful() = false, want true")
	}

	// Create a failed callback data
	failedCallback := &CallbackData{
		ResultCode: "01",
	}

	// Check if successful
	if failedCallback.IsSuccessful() {
		t.Errorf("IsSuccessful() = true, want false")
	}
}

func TestHandleCallback(t *testing.T) {
	// Create a client
	client := NewClient(Config{
		MerchantCode: "DXXXX",
		APIKey:       "DXXXXCX80TZJ85Q70QCI",
		IsSandbox:    true,
	})

	// Calculate the actual expected signature
	signatureStr := fmt.Sprintf("%s%d%s%s", "DXXXX", 40000, "ORDER123", "DXXXXCX80TZJ85Q70QCI")
	hash := md5.Sum([]byte(signatureStr))
	expectedSignature := hex.EncodeToString(hash[:])

	// Create form data
	form := url.Values{}
	form.Add("merchantCode", "DXXXX")
	form.Add("amount", "40000")
	form.Add("merchantOrderId", "ORDER123")
	form.Add("productDetail", "Test Product")
	form.Add("additionalParam", "")
	form.Add("paymentCode", "VC")
	form.Add("resultCode", "00")
	form.Add("reference", "DEV123456789")
	form.Add("signature", expectedSignature)

	// Create a request with form data
	req, err := http.NewRequest("POST", "/callback", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler
	var handlerCalled bool
	handler := func(data *CallbackData) error {
		handlerCalled = true
		return nil
	}

	// Handle callback
	client.HandleCallback(rr, req, handler)

	// Check if handler was called
	if !handlerCalled {
		t.Errorf("HandleCallback() did not call handler")
	}

	// Check response status code
	if rr.Code != http.StatusOK {
		t.Errorf("HandleCallback() status code = %d, want %d", rr.Code, http.StatusOK)
	}

	// Check response body
	if rr.Body.String() != "OK" {
		t.Errorf("HandleCallback() response body = %s, want OK", rr.Body.String())
	}
}
