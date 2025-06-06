package duitku

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// SandboxBaseURL is the base URL for the Duitku sandbox environment
	SandboxBaseURL = "https://sandbox.duitku.com/webapi/api"
	// ProductionBaseURL is the base URL for the Duitku production environment
	ProductionBaseURL = "https://passport.duitku.com/webapi/api"
)

// Config holds the configuration for the Duitku client
type Config struct {
	// MerchantCode is the merchant code provided by Duitku
	MerchantCode string
	// APIKey is the API key provided by Duitku
	APIKey string
	// IsSandbox determines whether to use the sandbox or production environment
	IsSandbox bool
	// HTTPClient is an optional custom HTTP client
	HTTPClient *http.Client
	// Logger is an optional custom logger
	Logger *log.Logger
	// Log every request and response
	LogEveryRequestAndResponse bool
}

// Client is the Duitku API client
type Client struct {
	config                     Config
	baseURL                    string
	httpClient                 *http.Client
	logger                     *log.Logger
	logEveryRequestAndResponse bool
}

// NewClient creates a new Duitku client with the provided configuration
func NewClient(config Config) *Client {
	baseURL := ProductionBaseURL
	if config.IsSandbox {
		baseURL = SandboxBaseURL
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	logger := config.Logger
	if logger == nil {
		logger = log.New(os.Stdout, "github.com/fatkulnurk/duitku-go: ", log.LstdFlags)
	}

	return &Client{
		config:                     config,
		baseURL:                    baseURL,
		httpClient:                 httpClient,
		logger:                     logger,
		logEveryRequestAndResponse: config.LogEveryRequestAndResponse,
	}
}

// createSignatureSHA256 creates a SHA256 signature from the provided parameters
func (c *Client) createSignatureSHA256(params ...string) string {
	var combined string
	for _, param := range params {
		combined += param
	}
	combined += c.config.APIKey

	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// createSignatureMD5 creates an MD5 signature from the provided parameters
func (c *Client) createSignatureMD5(params ...string) string {
	var combined string
	for _, param := range params {
		combined += param
	}
	combined += c.config.APIKey

	hash := md5.Sum([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// doRequest performs an HTTP request to the Duitku API
func (c *Client) doRequest(method, endpoint string, body interface{}, result interface{}) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	if c.logEveryRequestAndResponse {
		c.logger.Println("-=-=-=-= [client.go][doRequest] -=-=-=-=")
		c.logger.Printf("Request method: %s \n", method)
		c.logger.Printf("Request url: %s \n", url)
		c.logger.Printf("Request body: %+v \n", body)
		c.logger.Printf("Response status code: %d \n", resp.StatusCode)
		c.logger.Printf("Response status: %s \n", resp.Status)

		// Read the response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.logger.Printf("Error reading response body: %v\n", err)
		} else {
			c.logger.Printf("Response body: %s\n", string(respBody))
		}

		// Create a new reader with the same data for further processing
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

		c.logger.Println("-=-=-=-= [client.go][doRequest] -=-=-=-=")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("error decoding error response: %w", err)
		}
		return fmt.Errorf("API error: %s (code: %s)", errorResp.Message, errorResp.Code)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
	}

	return nil
}

// ErrorResponse represents an error response from the Duitku API
type ErrorResponse struct {
	Code    string `json:"responseCode"`
	Message string `json:"responseMessage"`
}

// Error returns the error message
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
