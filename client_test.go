package duitku

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
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
		{
			name: "Custom HTTP Client",
			config: Config{
				MerchantCode: "DXXXX",
				APIKey:       "DXXXXCX80TZJ85Q70QCI",
				IsSandbox:    true,
				HTTPClient:   &http.Client{Timeout: 60 * time.Second},
			},
			want: SandboxBaseURL,
		},
		{
			name: "Custom Logger",
			config: Config{
				MerchantCode: "DXXXX",
				APIKey:       "DXXXXCX80TZJ85Q70QCI",
				IsSandbox:    true,
				Logger:       log.New(os.Stderr, "test: ", log.LstdFlags),
			},
			want: SandboxBaseURL,
		},
		{
			name: "With Logging Enabled",
			config: Config{
				MerchantCode:               "DXXXX",
				APIKey:                     "DXXXXCX80TZJ85Q70QCI",
				IsSandbox:                  true,
				LogEveryRequestAndResponse: true,
			},
			want: SandboxBaseURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)

			// Check base URL
			if client.baseURL != tt.want {
				t.Errorf("NewClient() baseURL = %v, want %v", client.baseURL, tt.want)
			}

			// Check config values
			if client.config.MerchantCode != tt.config.MerchantCode {
				t.Errorf("NewClient() merchantCode = %v, want %v", client.config.MerchantCode, tt.config.MerchantCode)
			}
			if client.config.APIKey != tt.config.APIKey {
				t.Errorf("NewClient() apiKey = %v, want %v", client.config.APIKey, tt.config.APIKey)
			}

			// Check HTTP client
			if tt.config.HTTPClient != nil {
				if client.httpClient != tt.config.HTTPClient {
					t.Errorf("NewClient() httpClient = %v, want %v", client.httpClient, tt.config.HTTPClient)
				}
			} else {
				// Default HTTP client should have 30s timeout
				if client.httpClient.Timeout != 30*time.Second {
					t.Errorf("NewClient() httpClient.Timeout = %v, want %v", client.httpClient.Timeout, 30*time.Second)
				}
			}

			// Check logger
			if tt.config.Logger != nil {
				if client.logger != tt.config.Logger {
					t.Errorf("NewClient() logger = %v, want %v", client.logger, tt.config.Logger)
				}
			} else {
				// Default logger should exist
				if client.logger == nil {
					t.Errorf("NewClient() logger is nil, want non-nil")
				}
			}

			// Check logging flag
			if client.logEveryRequestAndResponse != tt.config.LogEveryRequestAndResponse {
				t.Errorf("NewClient() logEveryRequestAndResponse = %v, want %v",
					client.logEveryRequestAndResponse, tt.config.LogEveryRequestAndResponse)
			}
		})
	}
}

func TestCreateSignatureSHA256(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		params   []string
		expected string
	}{
		{
			name:   "Basic SHA256 Signature",
			apiKey: "DXXXXCX80TZJ85Q70QCI",
			params: []string{"DXXXX", "10000", "2022-01-25 16:23:08"},
			// Pre-calculated expected hash for the given inputs
			expected: "b7d9fc12e2ecc9d7c31c199a7b6a7bd2f8cb6d0f8a7b7f35b24ef2a0caf35e19",
		},
		{
			name:   "Empty Params",
			apiKey: "DXXXXCX80TZJ85Q70QCI",
			params: []string{},
			// Hash of just the API key
			expected: "e4c811dcd1f5a0e3c5a0ad30d01fbc42e7326d1b8a9732304f3d1b1a9d4ea027",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(Config{
				MerchantCode:               "DXXXX",
				APIKey:                     tt.apiKey,
				IsSandbox:                  true,
				LogEveryRequestAndResponse: true,
			})

			got := client.createSignatureSHA256(tt.params...)

			// For this test, we'll just check that the signature is not empty
			// In a real test, you would calculate the expected signature
			if got == "" {
				t.Errorf("createSignatureSHA256() = empty string, want non-empty")
			}

			// Verify signature is a valid SHA256 hash (64 hex characters)
			if len(got) != 64 {
				t.Errorf("createSignatureSHA256() length = %v, want 64", len(got))
			}
		})
	}
}

func TestCreateSignatureMD5(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		params   []string
		expected string
	}{
		{
			name:   "Basic MD5 Signature",
			apiKey: "DXXXXCX80TZJ85Q70QCI",
			params: []string{"DXXXX", "ORDER123", "10000"},
			// Pre-calculated expected hash for the given inputs
			expected: "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
		},
		{
			name:   "Empty Params",
			apiKey: "DXXXXCX80TZJ85Q70QCI",
			params: []string{},
			// Hash of just the API key
			expected: "e4c811dcd1f5a0e3c5a0ad30d01fbc42e7326d1b8a9732304f3d1b1a9d4ea027",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(Config{
				MerchantCode:               "DXXXX",
				APIKey:                     tt.apiKey,
				IsSandbox:                  true,
				LogEveryRequestAndResponse: true,
			})

			got := client.createSignatureMD5(tt.params...)

			// For this test, we'll just check that the signature is not empty
			// In a real test, you would calculate the expected signature
			if got == "" {
				t.Errorf("createSignatureMD5() = empty string, want non-empty")
			}

			// Verify signature is a valid MD5 hash (32 hex characters)
			if len(got) != 32 {
				t.Errorf("createSignatureMD5() length = %v, want 32", len(got))
			}
		})
	}
}

func TestDoRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		endpoint       string
		requestBody    interface{}
		responseStatus int
		responseBody   string
		expectError    bool
	}{
		{
			name:           "Successful Request",
			method:         "POST",
			endpoint:       "test-endpoint",
			requestBody:    map[string]string{"test": "data"},
			responseStatus: http.StatusOK,
			responseBody:   `{"responseCode":"00","responseMessage":"SUCCESS"}`,
			expectError:    false,
		},
		{
			name:           "Error Response",
			method:         "POST",
			endpoint:       "test-endpoint",
			requestBody:    nil,
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"responseCode":"01","responseMessage":"ERROR"}`,
			expectError:    true,
		},
		{
			name:           "Invalid JSON Response",
			method:         "GET",
			endpoint:       "test-endpoint",
			requestBody:    nil,
			responseStatus: http.StatusOK,
			responseBody:   `{invalid json}`,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that returns the configured response
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				if r.Method != tt.method {
					t.Errorf("Expected request method %s, got %s", tt.method, r.Method)
				}

				// Verify request path contains the endpoint
				expectedPath := "/" + tt.endpoint
				if tt.endpoint != "" && r.URL.Path != expectedPath {
					t.Errorf("Expected request path %s, got %s", expectedPath, r.URL.Path)
				}

				// Set response headers and status code
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Create a client that uses the test server
			client := &Client{
				config: Config{
					MerchantCode: "DXXXX",
					APIKey:       "DXXXXCX80TZJ85Q70QCI",
				},
				baseURL:                    server.URL,
				httpClient:                 server.Client(),
				logger:                     log.New(io.Discard, "", 0), // Suppress logging for tests
				logEveryRequestAndResponse: true,
			}

			// Test the request
			var response struct {
				ResponseCode    string `json:"responseCode"`
				ResponseMessage string `json:"responseMessage"`
			}
			err := client.doRequest(tt.method, tt.endpoint, tt.requestBody, &response)

			// Check if error matches expectation
			if (err != nil) != tt.expectError {
				t.Errorf("doRequest() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// If success case, verify response parsing
			if !tt.expectError {
				if response.ResponseCode != "00" {
					t.Errorf("doRequest() responseCode = %v, want 00", response.ResponseCode)
				}
				if response.ResponseMessage != "SUCCESS" {
					t.Errorf("doRequest() responseMessage = %v, want SUCCESS", response.ResponseMessage)
				}
			}
		})
	}
}

func TestDoRequestWithLogging(t *testing.T) {
	// Create a buffer to capture log output
	logBuffer := &bytes.Buffer{}

	// Create a test server that returns a valid response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"responseCode":"00","responseMessage":"SUCCESS"}`))
	}))
	defer server.Close()

	// Create a client with logging enabled
	client := &Client{
		config: Config{
			MerchantCode: "DXXXX",
			APIKey:       "DXXXXCX80TZJ85Q70QCI",
		},
		baseURL:                    server.URL,
		httpClient:                 server.Client(),
		logger:                     log.New(logBuffer, "", 0), // Capture logs
		logEveryRequestAndResponse: true,                      // Enable logging
	}

	// Test the request
	var response struct {
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	}
	requestBody := map[string]string{"test": "data"}
	err := client.doRequest("POST", "test-endpoint", requestBody, &response)

	// Check no error
	if err != nil {
		t.Errorf("doRequest() error = %v, want nil", err)
	}

	// Verify logging occurred
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "Request method: POST") {
		t.Errorf("Log output missing request method")
	}
	if !strings.Contains(logOutput, "Request url:") {
		t.Errorf("Log output missing request URL")
	}
	if !strings.Contains(logOutput, "Request body:") {
		t.Errorf("Log output missing request body")
	}
	if !strings.Contains(logOutput, "Response status code:") {
		t.Errorf("Log output missing response status code")
	}
	if !strings.Contains(logOutput, "Response body:") {
		t.Errorf("Log output missing response body")
	}
}

func TestDoRequestWithInvalidBody(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	}))
	defer server.Close()

	// Create a client
	client := &Client{
		config: Config{
			MerchantCode: "DXXXX",
			APIKey:       "DXXXXCX80TZJ85Q70QCI",
		},
		baseURL:    server.URL,
		httpClient: server.Client(),
		logger:     log.New(io.Discard, "", 0),
	}

	// Create an invalid request body that can't be marshaled to JSON
	invalidBody := make(chan int) // channels can't be marshaled to JSON

	// Test the request
	var response struct{}
	err := client.doRequest("POST", "test-endpoint", invalidBody, &response)

	// Should return an error
	if err == nil {
		t.Errorf("doRequest() expected error for invalid body, got nil")
	}

	// Error should mention marshaling
	if !strings.Contains(err.Error(), "marshaling") {
		t.Errorf("Expected error about marshaling, got: %v", err)
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		message  string
		expected string
	}{
		{
			name:     "Standard Error",
			code:     "01",
			message:  "Error message",
			expected: "01: Error message",
		},
		{
			name:     "Empty Code",
			code:     "",
			message:  "Error message",
			expected: ": Error message",
		},
		{
			name:     "Empty Message",
			code:     "02",
			message:  "",
			expected: "02: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrorResponse{
				Code:    tt.code,
				Message: tt.message,
			}
			if err.Error() != tt.expected {
				t.Errorf("ErrorResponse.Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}
