package duitku

// Payment method constants
const (
	// Bank Transfer
	PaymentMethodBCA      = "BC" // BCA Virtual Account
	PaymentMethodMandiri  = "M1" // Mandiri Virtual Account
	PaymentMethodPermata  = "BT" // Permata Virtual Account
	PaymentMethodBNI      = "I1" // BNI Virtual Account
	PaymentMethodBRI      = "BR" // BRI Virtual Account
	PaymentMethodCIMB     = "B1" // CIMB Niaga Virtual Account
	PaymentMethodDanamon  = "DN" // Danamon Virtual Account
	PaymentMethodMaybank  = "VA" // Maybank Virtual Account
	PaymentMethodSahabat  = "SA" // Sahabat Sampoerna Virtual Account
	PaymentMethodBSI      = "S1" // BSI Virtual Account
	
	// E-Wallet
	PaymentMethodOVO      = "OV" // OVO
	PaymentMethodShopeePay = "SP" // ShopeePay
	PaymentMethodLinkAja  = "LA" // LinkAja
	PaymentMethodDANA     = "DA" // DANA
	
	// QRIS
	PaymentMethodQRIS     = "QR" // QRIS
	
	// Retail Outlets
	PaymentMethodAlfamart = "A1" // Alfamart
	PaymentMethodIndomaret = "IR" // Indomaret
	
	// Credit Card
	PaymentMethodCreditCard = "VC" // Credit Card
	
	// Paylater
	PaymentMethodAkulaku  = "AK" // Akulaku
	PaymentMethodKredivo  = "K1" // Kredivo
	PaymentMethodAtome    = "AT" // Atome
)

// Transaction status codes
const (
	StatusSuccess         = "00" // Success
	StatusPending         = "01" // Pending
	StatusFailed          = "02" // Failed
	StatusCancelled       = "03" // Cancelled
	StatusExpired         = "04" // Expired
)

// Subscription frequency types
const (
	FrequencyDaily   = 1
	FrequencyWeekly  = 2
	FrequencyMonthly = 3
	FrequencyYearly  = 4
)
