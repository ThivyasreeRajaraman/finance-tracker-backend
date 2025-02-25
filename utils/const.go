package utils

import "errors"

const (
	TokenExpirationDays = 30
	EmailKey            = "email"
	Zero                = 0
	EmptyString         = ""
)

var (
	ErrEmptyRequestBody = errors.New("empty request body. Atleast provide one field")
	ErrInvalidDataType  = errors.New("invalid datatype")
)

var ValidCurrencies = map[string]bool{
	"INR": true,
	"USD": true,
	"EUR": true,
	"GBP": true,
	"JPY": true,
	"AUD": true,
	"CAD": true,
	"CHF": true,
	"CNY": true,
	"HKD": true,
	"NZD": true,
	"SEK": true,
	"KRW": true,
	"SGD": true,
	"NOK": true,
	"MXN": true,
	"RUB": true,
	"BRL": true,
	"TRY": true,
	"ZAR": true,
}

var ValidTransactionTypes = map[string]bool{
	"income":           true,
	"expense":          true,
	"lend":             true,
	"borrow":           true,
	"recurringExpense": true,
}

var ValidFrequencies = map[string]bool{
	"WEEKLY":  true,
	"MONTHLY": true,
	"YEARLY":  true,
}
