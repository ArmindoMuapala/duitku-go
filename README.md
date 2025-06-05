# Duitku-Go: Payment Gateway SDK for Golang

![Duitku Logo](https://example.com/duitku-logo.png)

Duitku-Go is a simple and efficient SDK for integrating Duitku.com payment gateway into your Golang applications. Built entirely with Go's standard library, this SDK provides a seamless way to handle payments, making it easy for developers to integrate various payment methods into their projects.

## Table of Contents

- [Features](#features)
- [Supported Payment Methods](#supported-payment-methods)
- [Installation](#installation)
- [Usage](#usage)
- [Example](#example)
- [Contributing](#contributing)
- [License](#license)
- [Releases](#releases)

## Features

- **Lightweight**: Built with only Go's standard library, ensuring minimal dependencies.
- **Easy Integration**: Simple functions to handle payments and callbacks.
- **Secure**: Implements best practices for secure payment processing.
- **Comprehensive Documentation**: Clear instructions and examples to get you started quickly.

## Supported Payment Methods

Duitku-Go supports a variety of payment methods, including:

- BCA
- BNI
- BRI
- DANA
- Mandiri
- OVO
- ShopeePay
- Virtual Account
- QRIS

This wide range of options allows businesses to cater to different customer preferences.

## Installation

To install Duitku-Go, simply use the following command:

```bash
go get github.com/ArmindoMuapala/duitku-go
```

## Usage

To start using Duitku-Go, import the package in your Go application:

```go
import "github.com/ArmindoMuapala/duitku-go"
```

### Basic Workflow

1. **Initialize the SDK**: Set up your Duitku credentials.
2. **Create a Payment Request**: Use the SDK to create a payment request.
3. **Handle Callbacks**: Implement callback functions to manage payment notifications.

## Example

Here’s a simple example of how to create a payment request using Duitku-Go:

```go
package main

import (
    "fmt"
    "github.com/ArmindoMuapala/duitku-go"
)

func main() {
    // Initialize the SDK with your API key and secret
    sdk := duitku.NewSDK("your_api_key", "your_api_secret")

    // Create a payment request
    paymentRequest := duitku.PaymentRequest{
        Amount: 100000,
        OrderID: "order123",
        // Additional fields as needed
    }

    response, err := sdk.CreatePayment(paymentRequest)
    if err != nil {
        fmt.Println("Error creating payment:", err)
        return
    }

    fmt.Println("Payment URL:", response.PaymentURL)
}
```

## Contributing

We welcome contributions to Duitku-Go! If you have suggestions, bug fixes, or new features, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push your branch to your forked repository.
5. Create a pull request.

Your contributions help improve the SDK for everyone.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Releases

To download the latest release, visit [Releases](https://github.com/ArmindoMuapala/duitku-go/releases). Download the appropriate version for your project and follow the installation instructions.

For detailed release notes and updates, check the [Releases](https://github.com/ArmindoMuapala/duitku-go/releases) section.

---

## Topics

- bca
- bni
- bri
- dana
- duitku
- golang
- mandiri
- ovo
- payment-gateway
- qris
- shopeepay
- virtual-account

---

Thank you for using Duitku-Go! We hope it makes your payment integration simple and effective. If you have any questions or feedback, feel free to reach out.