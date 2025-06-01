# Duitku Go 💸

[![Go](https://github.com/fatkulnurk/duitku-go/actions/workflows/go.yml/badge.svg)](https://github.com/fatkulnurk/duitku-go/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/fatkulnurk/duitku-go)](https://goreportcard.com/report/github.com/fatkulnurk/duitku-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/fatkulnurk/duitku-go.svg)](https://pkg.go.dev/github.com/fatkulnurk/duitku-go)

A simple and lightweight Duitku.com Payment Gateway SDK for Golang — built with only Go's standard library. No external dependencies, making it ideal for minimal and secure payment gateway integrations. This package implements [Duitku API v2](https://docs.duitku.com/api/en/#introduction).

## Table of Contents
- [🌟 Overview](#-overview)
- [📋 Requirements](#-requirements)
- [📦 Installation](#-installation)
- [📚 Usage](#-usage)
- [✨ Features](#-features)
  - [Payment Methods 💳](#payment-methods-)
  - [Transactions 🛒](#transactions-)
  - [Callbacks 📡](#callbacks-)
  - [Additional Features 🔧](#additional-features-)
  - [Payment Methods Supported 💰](#payment-methods-supported-)
- [🚀 Advanced Usage](#-advanced-usage)
  - [Subscription Payments](#subscription-payments)
- [🧪 Example Application](#-example-application)
- [👥 Contributing](#-contributing)
- [📄 License](#-license)

## 🌟 Overview

Duitku is a payment gateway service that provides various payment methods for Indonesian merchants. This Go package provides a clean, idiomatic interface to integrate with Duitku's payment services.

[🔼 Jump to Table of Contents](#table-of-contents)

## 📋 Requirements

- Go 1.20 or higher
- Duitku merchant account and API credentials

[🔼 Jump to Table of Contents](#table-of-contents)

## 📦 Installation

```bash
go get github.com/fatkulnurk/duitku-go
```

[🔼 Jump to Table of Contents](#table-of-contents)

## 📚 Usage

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
		LogEveryRequestAndResponse: true, // Set to false to disable logging every request and response
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
		PaymentMethod:   duitku.PaymentMethodCreditCard, // Credit Card
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

[🔼 Jump to Table of Contents](#table-of-contents)

## ✨ Features

### Payment Methods 💳

- ✅ **Get Available Payment Methods** - [API Reference](https://docs.duitku.com/api/en/#get-payment-method)
  - Retrieve all available payment methods based on transaction amount
  - Supports all payment methods including VA, e-wallets, QRIS, retail outlets, and more
  - Automatically filters methods based on minimum transaction amount

[🔼 Jump to Table of Contents](#table-of-contents)

### Transactions 🛒

- ✅ **Create Transaction** - [API Reference](https://docs.duitku.com/api/en/#create-invoice-request-body-json)
  - Create payment invoices with complete customer details
  - Support for item details and additional parameters
  - Configurable expiry period
  - Returns payment URL and reference for customer redirection

- ✅ **Check Transaction Status** - [API Reference](https://docs.duitku.com/api/en/#check-transaction)
  - Verify payment status using merchant order ID
  - Get detailed transaction information including amount and reference

[🔼 Jump to Table of Contents](#table-of-contents)

### Callbacks 📡

- ✅ **Handle Payment Notifications** - [API Reference](https://docs.duitku.com/api/en/#callback)
  - Secure callback handling with signature verification
  - Automatic parsing of callback data
  - Easy-to-use handler function for processing successful payments

[🔼 Jump to Table of Contents](#table-of-contents)

### Additional Features 🔧

- ✅ **Environment Support**
  - Sandbox environment for testing - [Sandbox Dashboard](https://sandbox.duitku.com/)
  - Production environment for live transactions

- ✅ **Security**
  - SHA256 and MD5 signature generation and verification
  - Secure API key handling

- ✅ **Comprehensive Testing**
  - Full test suite with high coverage
  - CI/CD integration with GitHub Actions

[🔼 Jump to Table of Contents](#table-of-contents)

### Payment Methods Supported 💰

| Category | Payment Method | Code | Status |
|----------|---------------|------|--------|
| **Credit Card** | Credit Card | VC | ✅ |
| **Virtual Account** | BCA VA | BC | ✅ |
| | Mandiri VA | M2 | ✅ |
| | Permata VA | BT | ✅ |
| | BNI VA | I1 | ✅ |
| | BRI VA | BR | ✅ |
| | CIMB Niaga VA | B1 | ✅ |
| | Danamon VA | DM | ✅ |
| | Maybank VA | VA | ✅ |
| | Sahabat Sampoerna VA | S1 | ✅ |
| | BSI VA | BV | ✅ |
| | Bank Neo Commerce VA | NC | ✅ |
| | Bank Artha Graha VA | AG | ✅ |
| **E-Wallet** | OVO | OV | ✅ |
| | OVO Account Link | OL | ✅ |
| | ShopeePay Apps | SA | ✅ |
| | Shopee Account Link | SL | ✅ |
| | LinkAja (Fixed Fee) | LF | ✅ |
| | LinkAja (Percentage Fee) | LA | ✅ |
| | DANA | DA | ✅ |
| | Jenius Pay | JP | ✅ |
| **QRIS** | QRIS ShopeePay | SP | ✅ |
| | QRIS Nobu | QN | ✅ |
| | QRIS Dana | DQ | ✅ |
| | QRIS Gudang Voucher | GQ | ✅ |
| | QRIS Nusapay | SQ | ✅ |
| **Retail Outlets** | Alfamart/Pegadaian/POS | FT | ✅ |
| | Indomaret | IR | ✅ |
| **Paylater** | Indodana Paylater | ID | ✅ |
| | Atome | AT | ✅ |

[🔼 Jump to Table of Contents](#table-of-contents)

## 🚀 Advanced Usage

### Subscription Payments

- ✅ **Subscription Support** - [API Reference](https://docs.duitku.com/api/en/#subscription)
  - Create recurring payment schedules
  - Support for daily, weekly, monthly, and yearly billing cycles
  - Configurable frequency and interval settings

```go
// Create a subscription transaction
isSubscription := true
transaction := duitku.TransactionRequest{
    // Basic transaction details
    PaymentAmount:   50000,
    PaymentMethod:   duitku.PaymentMethodCreditCard,
    MerchantOrderID: "SUB123",
    ProductDetails:  "Monthly Subscription",
    CustomerVaName:  "John Doe",
    Email:           "customer@example.com",
    CallbackURL:     "https://example.com/callback",
    ReturnURL:       "https://example.com/return",
    ExpiryPeriod:    60,
    
    // Enable subscription
    IsSubscription: &isSubscription,
    
    // Subscription details
    SubscriptionDetail: &duitku.SubscriptionDetail{
        Description:     "Monthly Premium Plan",
        FrequencyType:    duitku.FrequencyMonthly,
        FrequencyInterval: 1,
        TotalNoOfCycles:   12,
    },
}
```

[🔼 Jump to Table of Contents](#table-of-contents)

## 🧪 Example Application

This package includes a fully functional example application that demonstrates how to use all the features of the Duitku Go client. The example application is a simple web server that allows you to:

- Create payment transactions
- Handle payment callbacks
- Check transaction status
- View payment details

To run the example application:

```bash
# Set your Duitku credentials
export DUITKU_MERCHANT_CODE="your_merchant_code"
export DUITKU_API_KEY="your_api_key"

# Run the example application
cd example
go run main.go
```

Then open your browser to http://localhost:8080 to see the example application in action.

[🔼 Jump to Table of Contents](#table-of-contents)

## 👥 Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests to ensure everything works (`go test ./...`)
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

Please make sure your code passes all tests and only use Go's standard library.

[🔼 Jump to Table of Contents](#table-of-contents)

## 📄 License

MIT License - See [LICENSE](LICENSE) for details.

[🔼 Jump to Table of Contents](#table-of-contents)
