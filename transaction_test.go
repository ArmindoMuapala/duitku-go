package duitku

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	// Create a test server that returns a successful response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Decode request body
		var request struct {
			MerchantCode    string `json:"merchantCode"`
			PaymentAmount   int    `json:"paymentAmount"`
			PaymentMethod   string `json:"paymentMethod"`
			MerchantOrderID string `json:"merchantOrderId"`
			ProductDetails  string `json:"productDetails"`
			Email           string `json:"email"`
			Signature       string `json:"signature"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Error decoding request body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Check request fields
		if request.MerchantCode != "DXXXX" {
			t.Errorf("Expected MerchantCode: DXXXX, got %s", request.MerchantCode)
		}
		if request.PaymentAmount != 40000 {
			t.Errorf("Expected PaymentAmount: 40000, got %d", request.PaymentAmount)
		}
		if request.PaymentMethod != "VC" {
			t.Errorf("Expected PaymentMethod: VC, got %s", request.PaymentMethod)
		}
		if request.MerchantOrderID != "ORDER123" {
			t.Errorf("Expected MerchantOrderID: ORDER123, got %s", request.MerchantOrderID)
		}
		if request.ProductDetails != "Test Product" {
			t.Errorf("Expected ProductDetails: Test Product, got %s", request.ProductDetails)
		}
		if request.Email != "customer@example.com" {
			t.Errorf("Expected Email: customer@example.com, got %s", request.Email)
		}
		if request.Signature == "" {
			t.Errorf("Expected Signature to be non-empty")
		}

		// Return a successful response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"merchantCode": "DXXXX",
			"reference": "DEV123456789",
			"paymentUrl": "https://sandbox.duitku.com/payment/DEV123456789",
			"vaNumber": "123456789",
			"amount": "40000",
			"statusCode": "00",
			"statusMessage": "SUCCESS"
		}`))
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := &Client{
		config: Config{
			MerchantCode: "DXXXX",
			APIKey:       "DXXXXCX80TZJ85Q70QCI",
		},
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	// Create a transaction request
	request := TransactionRequest{
		PaymentAmount:   40000,
		PaymentMethod:   "VC",
		MerchantOrderID: "ORDER123",
		ProductDetails:  "Test Product",
		CustomerVaName:  "John Doe",
		Email:           "customer@example.com",
		CallbackURL:     "https://example.com/callback",
		ReturnURL:       "https://example.com/return",
		ExpiryPeriod:    10,
	}

	// Test creating a transaction
	response, err := client.CreateTransaction(request)
	if err != nil {
		t.Errorf("CreateTransaction() error = %v, want nil", err)
	}

	// Check the response
	if response.MerchantCode != "DXXXX" {
		t.Errorf("Response MerchantCode = %s, want DXXXX", response.MerchantCode)
	}
	if response.Reference != "DEV123456789" {
		t.Errorf("Response Reference = %s, want DEV123456789", response.Reference)
	}
	if response.PaymentURL != "https://sandbox.duitku.com/payment/DEV123456789" {
		t.Errorf("Response PaymentURL = %s, want https://sandbox.duitku.com/payment/DEV123456789", response.PaymentURL)
	}
	if response.VANumber != "123456789" {
		t.Errorf("Response VANumber = %s, want 123456789", response.VANumber)
	}
	if response.Amount != "40000" {
		t.Errorf("Response Amount = %s, want 40000", response.Amount)
	}
	if response.StatusCode != "00" {
		t.Errorf("Response StatusCode = %s, want 00", response.StatusCode)
	}
	if response.StatusMessage != "SUCCESS" {
		t.Errorf("Response StatusMessage = %s, want SUCCESS", response.StatusMessage)
	}
}

func TestCreateTransactionError(t *testing.T) {
	// Create a test server that returns an error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"merchantCode": "DXXXX",
			"reference": "",
			"paymentUrl": "",
			"amount": 40000,
			"statusCode": "01",
			"statusMessage": "ERROR"
		}`))
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := &Client{
		config: Config{
			MerchantCode: "DXXXX",
			APIKey:       "DXXXXCX80TZJ85Q70QCI",
		},
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	// Create a transaction request
	request := TransactionRequest{
		PaymentAmount:   40000,
		PaymentMethod:   "VC",
		MerchantOrderID: "ORDER123",
		ProductDetails:  "Test Product",
		CustomerVaName:  "John Doe",
		Email:           "customer@example.com",
		CallbackURL:     "https://example.com/callback",
		ReturnURL:       "https://example.com/return",
		ExpiryPeriod:    10,
	}

	// Test creating a transaction with error
	_, err := client.CreateTransaction(request)
	if err == nil {
		t.Errorf("CreateTransaction() error = nil, want error")
	}
}

func TestCheckTransaction(t *testing.T) {
	// Create a test server that returns a successful response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Decode request body
		var request CheckTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Error decoding request body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Check request fields
		if request.MerchantCode != "DXXXX" {
			t.Errorf("Expected MerchantCode: DXXXX, got %s", request.MerchantCode)
		}
		if request.MerchantOrderID != "ORDER123" {
			t.Errorf("Expected MerchantOrderID: ORDER123, got %s", request.MerchantOrderID)
		}
		if request.Signature == "" {
			t.Errorf("Expected Signature to be non-empty")
		}

		// Return a successful response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"merchantOrderId": "ORDER123",
			"reference": "DEV123456789",
			"amount": "40000",
			"fee": "2000",
			"statusCode": "00",
			"statusMessage": "SUCCESS"
		}`))
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := &Client{
		config: Config{
			MerchantCode: "DXXXX",
			APIKey:       "DXXXXCX80TZJ85Q70QCI",
		},
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	// Test checking a transaction
	response, err := client.CheckTransaction("ORDER123")
	if err != nil {
		t.Errorf("CheckTransaction() error = %v, want nil", err)
	}

	// Check the response
	if response.MerchantOrderID != "ORDER123" {
		t.Errorf("Response MerchantOrderID = %s, want ORDER123", response.MerchantOrderID)
	}
	if response.Reference != "DEV123456789" {
		t.Errorf("Response Reference = %s, want DEV123456789", response.Reference)
	}
	if response.Amount != "40000" {
		t.Errorf("Response Amount = %s, want 40000", response.Amount)
	}
	if response.Fee != "2000" {
		t.Errorf("Response Fee = %s, want 2000", response.Fee)
	}
	if response.StatusCode != "00" {
		t.Errorf("Response StatusCode = %s, want 00", response.StatusCode)
	}
	if response.StatusMessage != "SUCCESS" {
		t.Errorf("Response StatusMessage = %s, want SUCCESS", response.StatusMessage)
	}
}
