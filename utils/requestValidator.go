package utils

import (
	"errors"
	"strings"
)

func IsValidCurrency(currency string) error {
	normalizedCurrency := strings.ToUpper(currency)
	if _, ok := validCurrencies[normalizedCurrency]; !ok {
		return errors.New("invalid currency code")
	}
	return nil
}
