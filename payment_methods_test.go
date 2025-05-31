package duitku

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPaymentMethods(t *testing.T) {
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
		var request GetPaymentMethodsRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Error decoding request body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Check request fields
		if request.MerchantCode != "DXXXX" {
			t.Errorf("Expected MerchantCode: DXXXX, got %s", request.MerchantCode)
		}
		if request.Amount != 10000 {
			t.Errorf("Expected Amount: 10000, got %d", request.Amount)
		}
		if request.DateTime == "" {
			t.Errorf("Expected DateTime to be non-empty")
		}
		if request.Signature == "" {
			t.Errorf("Expected Signature to be non-empty")
		}

		// Return a successful response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"paymentFee": [
				{
					"paymentMethod": "VA",
					"paymentName": "MAYBANK VA",
					"paymentImage": "https://images.duitku.com/hotlink-ok/VA.PNG",
					"totalFee": "0"
				},
				{
					"paymentMethod": "BT",
					"paymentName": "PERMATA VA",
					"paymentImage": "https://images.duitku.com/hotlink-ok/PERMATA.PNG",
					"totalFee": "0"
				}
			],
			"responseCode": "00",
			"responseMessage": "SUCCESS"
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

	// Test getting payment methods
	methods, err := client.GetPaymentMethods(10000)
	if err != nil {
		t.Errorf("GetPaymentMethods() error = %v, want nil", err)
	}

	// Check the response
	if len(methods) != 2 {
		t.Errorf("GetPaymentMethods() returned %d methods, want 2", len(methods))
	}

	// Check the first method
	if methods[0].PaymentMethod != "VA" {
		t.Errorf("First method PaymentMethod = %s, want VA", methods[0].PaymentMethod)
	}
	if methods[0].PaymentName != "MAYBANK VA" {
		t.Errorf("First method PaymentName = %s, want MAYBANK VA", methods[0].PaymentName)
	}

	// Check the second method
	if methods[1].PaymentMethod != "BT" {
		t.Errorf("Second method PaymentMethod = %s, want BT", methods[1].PaymentMethod)
	}
	if methods[1].PaymentName != "PERMATA VA" {
		t.Errorf("Second method PaymentName = %s, want PERMATA VA", methods[1].PaymentName)
	}
}

func TestGetPaymentMethodsError(t *testing.T) {
	// Create a test server that returns an error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"paymentFee": [],
			"responseCode": "01",
			"responseMessage": "ERROR"
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

	// Test getting payment methods with error
	_, err := client.GetPaymentMethods(10000)
	if err == nil {
		t.Errorf("GetPaymentMethods() error = nil, want error")
	}
}
