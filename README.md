# Duitku Go

A simple and lightweight Duitku.com Payment Gateway SDK for Golang — built with only Go's standard library. No external dependencies, making it ideal for minimal and secure payment gateway integrations.

See [Duitku Payment Gateway API](https://docs.duitku.com/api/id/) for more information about the API.

> **NOTE: This package is currently under development and NOT ready for production use.**

## Installation

```bash
go get github.com/fatkulnurk/duitku-go
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/fatkulnurk/duitku-go"
)

func main() {
	// Initialize the client
	client := duitku.NewClient(duitku.Config{
		MerchantCode: "YOUR_MERCHANT_CODE",
		APIKey:       "YOUR_API_KEY",
		IsSandbox:    true, // Set to false for production
	})

	// Get available payment methods
	paymentMethods, err := client.GetPaymentMethods(10000)
	if err != nil {
		log.Fatalf("Error getting payment methods: %v", err)
	}

	for _, method := range paymentMethods {
		fmt.Printf("Payment Method: %s (%s)\n", method.PaymentName, method.PaymentMethod)
	}

	// Create a transaction
	transaction := duitku.TransactionRequest{
		PaymentAmount:   40000,
		PaymentMethod:   "VC", // Credit Card
		MerchantOrderID: "ORDER123",
		ProductDetails:  "Test Product",
		CustomerVaName:  "John Doe",
		Email:           "customer@example.com",
		CallbackURL:     "https://example.com/callback",
		ReturnURL:       "https://example.com/return",
		ExpiryPeriod:    10, // 10 minutes
	}

	result, err := client.CreateTransaction(transaction)
	if err != nil {
		log.Fatalf("Error creating transaction: %v", err)
	}

	fmt.Printf("Payment URL: %s\n", result.PaymentURL)
	fmt.Printf("Reference: %s\n", result.Reference)

	// Check transaction status
	status, err := client.CheckTransaction("ORDER123")
	if err != nil {
		log.Fatalf("Error checking transaction: %v", err)
	}

	fmt.Printf("Status: %s (%s)\n", status.StatusMessage, status.StatusCode)
}
```

## Features

- Get available payment methods
- Create transactions
- Check transaction status
- Handle callbacks
- Support for all Duitku payment methods
- Sandbox and production environments

## License

MIT
