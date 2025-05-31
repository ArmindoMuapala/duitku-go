/*
Package duitku provides a Go client for the Duitku Payment Gateway API.

Duitku is an Indonesian payment gateway that supports various payment methods
including virtual accounts, e-wallets, retail outlets, credit cards, and more.

# Basic Usage

Initialize a client with your merchant code and API key:

	client := duitku.NewClient(duitku.Config{
		MerchantCode: "YOUR_MERCHANT_CODE",
		APIKey:       "YOUR_API_KEY",
		IsSandbox:    true, // Set to false for production
	})

# Getting Available Payment Methods

Get available payment methods for a specific amount:

	paymentMethods, err := client.GetPaymentMethods(10000)
	if err != nil {
		log.Fatalf("Error getting payment methods: %v", err)
	}

	for _, method := range paymentMethods {
		fmt.Printf("Payment Method: %s (%s)\n", method.PaymentName, method.PaymentMethod)
	}

# Creating a Transaction

Create a transaction with the minimum required fields:

	transaction := duitku.TransactionRequest{
		PaymentAmount:   40000,
		PaymentMethod:   duitku.PaymentMethodBCA, // BCA Virtual Account
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

# Creating a Transaction with Additional Details

Create a transaction with additional details like customer information and item details:

	transaction := duitku.TransactionRequest{
		PaymentAmount:   40000,
		PaymentMethod:   duitku.PaymentMethodBCA,
		MerchantOrderID: "ORDER123",
		ProductDetails:  "Test Product",
		CustomerVaName:  "John Doe",
		Email:           "customer@example.com",
		PhoneNumber:     "08123456789",
		CallbackURL:     "https://example.com/callback",
		ReturnURL:       "https://example.com/return",
		ExpiryPeriod:    10,
		
		// Add customer details
		CustomerDetail: &duitku.CustomerDetail{
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "customer@example.com",
			PhoneNumber: "08123456789",
			BillingAddress: &duitku.Address{
				FirstName:   "John",
				LastName:    "Doe",
				Address:     "Jl. Kembangan Raya",
				City:        "Jakarta",
				PostalCode:  "11530",
				Phone:       "08123456789",
				CountryCode: "ID",
			},
			ShippingAddress: &duitku.Address{
				FirstName:   "John",
				LastName:    "Doe",
				Address:     "Jl. Kembangan Raya",
				City:        "Jakarta",
				PostalCode:  "11530",
				Phone:       "08123456789",
				CountryCode: "ID",
			},
		},
		
		// Add item details
		ItemDetails: []duitku.ItemDetail{
			{
				Name:     "Product 1",
				Price:    10000,
				Quantity: 1,
			},
			{
				Name:     "Product 2",
				Price:    30000,
				Quantity: 1,
			},
		},
	}

# Checking Transaction Status

Check the status of a transaction:

	status, err := client.CheckTransaction("ORDER123")
	if err != nil {
		log.Fatalf("Error checking transaction: %v", err)
	}

	fmt.Printf("Status: %s (%s)\n", status.StatusMessage, status.StatusCode)

# Handling Callbacks

Handle callbacks from Duitku in your HTTP handler:

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		client.HandleCallback(w, r, func(data *duitku.CallbackData) error {
			// Check if payment is successful
			if data.IsSuccessful() {
				// Update order status in database
				fmt.Printf("Payment successful for order %s\n", data.MerchantOrderID)
			} else {
				// Handle failed payment
				fmt.Printf("Payment failed for order %s with code %s\n", 
					data.MerchantOrderID, data.ResultCode)
			}
			return nil
		})
	})

# Payment Methods

The package provides constants for all payment methods supported by Duitku:

	duitku.PaymentMethodBCA      // BCA Virtual Account
	duitku.PaymentMethodMandiri  // Mandiri Virtual Account
	duitku.PaymentMethodPermata  // Permata Virtual Account
	duitku.PaymentMethodBNI      // BNI Virtual Account
	duitku.PaymentMethodBRI      // BRI Virtual Account
	duitku.PaymentMethodCIMB     // CIMB Niaga Virtual Account
	duitku.PaymentMethodDanamon  // Danamon Virtual Account
	duitku.PaymentMethodMaybank  // Maybank Virtual Account
	duitku.PaymentMethodSahabat  // Sahabat Sampoerna Virtual Account
	duitku.PaymentMethodBSI      // BSI Virtual Account
	duitku.PaymentMethodOVO      // OVO
	duitku.PaymentMethodShopeePay // ShopeePay
	duitku.PaymentMethodLinkAja  // LinkAja
	duitku.PaymentMethodDANA     // DANA
	duitku.PaymentMethodQRIS     // QRIS
	duitku.PaymentMethodAlfamart // Alfamart
	duitku.PaymentMethodIndomaret // Indomaret
	duitku.PaymentMethodCreditCard // Credit Card
	duitku.PaymentMethodAkulaku  // Akulaku
	duitku.PaymentMethodKredivo  // Kredivo
	duitku.PaymentMethodAtome    // Atome

# Transaction Status Codes

The package provides constants for transaction status codes:

	duitku.StatusSuccess    // Success
	duitku.StatusPending    // Pending
	duitku.StatusFailed     // Failed
	duitku.StatusCancelled  // Cancelled
	duitku.StatusExpired    // Expired

# Subscription Frequency Types

For credit card subscription transactions, the package provides constants for frequency types:

	duitku.FrequencyDaily   // Daily
	duitku.FrequencyWeekly  // Weekly
	duitku.FrequencyMonthly // Monthly
	duitku.FrequencyYearly  // Yearly
*/
package duitku
