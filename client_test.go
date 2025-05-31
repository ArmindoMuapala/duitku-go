package duitku

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{
			name: "Sandbox",
			config: Config{
				MerchantCode: "DXXXX",
				APIKey:       "DXXXXCX80TZJ85Q70QCI",
				IsSandbox:    true,
			},
			want: SandboxBaseURL,
		},
		{
			name: "Production",
			config: Config{
				MerchantCode: "DXXXX",
				APIKey:       "DXXXXCX80TZJ85Q70QCI",
				IsSandbox:    false,
			},
			want: ProductionBaseURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)
			if client.baseURL != tt.want {
				t.Errorf("NewClient() baseURL = %v, want %v", client.baseURL, tt.want)
			}
			if client.config.MerchantCode != tt.config.MerchantCode {
				t.Errorf("NewClient() merchantCode = %v, want %v", client.config.MerchantCode, tt.config.MerchantCode)
			}
			if client.config.APIKey != tt.config.APIKey {
				t.Errorf("NewClient() apiKey = %v, want %v", client.config.APIKey, tt.config.APIKey)
			}
		})
	}
}

func TestCreateSignatureSHA256(t *testing.T) {
	client := NewClient(Config{
		MerchantCode: "DXXXX",
		APIKey:       "DXXXXCX80TZJ85Q70QCI",
		IsSandbox:    true,
	})

	signature := client.createSignatureSHA256("DXXXX", "10000", "2022-01-25 16:23:08")
	// The expected value should be calculated based on the actual algorithm
	// This is just a placeholder test
	if signature == "" {
		t.Errorf("createSignatureSHA256() = empty string, want non-empty")
	}
}

func TestCreateSignatureMD5(t *testing.T) {
	client := NewClient(Config{
		MerchantCode: "DXXXX",
		APIKey:       "DXXXXCX80TZJ85Q70QCI",
		IsSandbox:    true,
	})

	signature := client.createSignatureMD5("DXXXX", "ORDER123", "10000")
	// The expected value should be calculated based on the actual algorithm
	// This is just a placeholder test
	if signature == "" {
		t.Errorf("createSignatureMD5() = empty string, want non-empty")
	}
}

func TestDoRequest(t *testing.T) {
	// Create a test server that returns a successful response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"responseCode":"00","responseMessage":"SUCCESS"}`))
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

	// Test a successful request
	var response struct {
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	}
	err := client.doRequest("POST", "", nil, &response)
	if err != nil {
		t.Errorf("doRequest() error = %v, want nil", err)
	}
	if response.ResponseCode != "00" {
		t.Errorf("doRequest() responseCode = %v, want 00", response.ResponseCode)
	}
	if response.ResponseMessage != "SUCCESS" {
		t.Errorf("doRequest() responseMessage = %v, want SUCCESS", response.ResponseMessage)
	}
}

func TestErrorResponse(t *testing.T) {
	err := ErrorResponse{
		Code:    "01",
		Message: "Error message",
	}
	expected := "01: Error message"
	if err.Error() != expected {
		t.Errorf("ErrorResponse.Error() = %v, want %v", err.Error(), expected)
	}
}
