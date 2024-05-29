package utils

import (
	"strings"
)

func IsValidCurrency(currency string) error {
	normalizedCurrency := strings.ToUpper(currency)
	if _, ok := validCurrencies[normalizedCurrency]; !ok {
		return CreateError("invalid currency code")
	}
	return nil
}

func IsValidFrequency(frequency string) error {
	normalizedFrequency := strings.ToUpper(frequency)
	if _, ok := validFrequencies[normalizedFrequency]; !ok {
		return CreateError("invalid frequency")
	}
	return nil
}
