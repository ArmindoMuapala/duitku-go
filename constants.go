package duitku

// Payment method constants
const (
	// Credit Card
	PaymentMethodCreditCard = "VC" // Credit Card

	// Virtual Account
	PaymentMethodBCA         = "BC" // BCA Virtual Account
	PaymentMethodMandiri     = "M2" // Mandiri Virtual Account
	PaymentMethodMaybank     = "VA" // Maybank Virtual Account
	PaymentMethodBNI         = "I1" // BNI Virtual Account
	PaymentMethodCIMB        = "B1" // CIMB Niaga Virtual Account
	PaymentMethodPermata     = "BT" // Permata Bank Virtual Account
	PaymentMethodATMBersama  = "A1" // ATM Bersama
	PaymentMethodArthaGraha  = "AG" // Bank Artha Graha / Bank Artha Graha Virtual Account
	PaymentMethodNeoCommerce = "NC" // Bank Neo Commerce / BNC / Bank Neo Commerce Virtual Account
	PaymentMethodBRI         = "BR" // BRIVA / Bank BRI / Bank Rakyat Indonesia Virtual Account
	PaymentMethodSahabat     = "S1" // Bank Sahabat Sampoerna / Bank Sahabat Sampoerna Virtual Account
	PaymentMethodDanamon     = "DM" // Danamon Virtual Account
	PaymentMethodBSI         = "BV" // BSI Virtual Account

	// Retail Outlets
	PaymentMethodAlfamart  = "FT" // Alfamart (Pegadaian/ALFA/Pos)
	PaymentMethodPegadaian = "FT" // Pegadaian (Pegadaian/ALFA/Pos)
	PaymentMethodPOS       = "FT" // POS (Pegadaian/ALFA/Pos)
	PaymentMethodIndomaret = "IR" // Indomaret

	// E-Wallet
	PaymentMethodOVO            = "OV" // OVO (Support Void)
	PaymentMethodShopeePay      = "SA" // Shopee Pay Apps (Support Void)
	PaymentMethodLinkAjaFixed   = "LF" // LinkAja Apps (Fixed Fee)
	PaymentMethodLinkAjaPercent = "LA" // LinkAja Apps (Percentage Fee)
	PaymentMethodDANA           = "DA" // DANA
	PaymentMethodShopeeLink     = "SL" // Shopee Pay Account Link
	PaymentMethodOVOLink        = "OL" // OVO Account Link
	PaymentMethodJeniusPay      = "JP" // Jenius Pay

	// QRIS
	PaymentMethodQrisShopeePay     = "SP" // QRIS ShopeePay
	PaymentMethodQrisNobu          = "QN" // QRIS Nobu
	PaymentMethodQrisDana          = "DQ" // QRIS Dana
	PaymentMethodQrisGudangVoucher = "GQ" // QRIS Gudang Voucher
	PaymentMethodQrisNusapay       = "SQ" // QRIS Nusapay

	// Paylater/Credit
	PaymentMethodIndodanaPaylater = "ID" // Indodana Paylater
	PaymentMethodAtome            = "AT" // ATOME
)

// Transaction status codes
const (
	StatusSuccess   = "00" // Success
	StatusPending   = "01" // Pending
	StatusFailed    = "02" // Failed
	StatusCancelled = "03" // Cancelled
	StatusExpired   = "04" // Expired
)

// Subscription frequency types
const (
	FrequencyDaily   = 1 // Daily
	FrequencyWeekly  = 2 // Weekly
	FrequencyMonthly = 3 // Monthly
	FrequencyYearly  = 4 // Yearly
)
