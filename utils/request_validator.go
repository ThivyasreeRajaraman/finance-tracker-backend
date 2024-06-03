package utils

import (
	"strings"
	"time"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
)

func IsValidCurrency(currency string) error {
	normalizedCurrency := strings.ToUpper(currency)
	if _, ok := ValidCurrencies[normalizedCurrency]; !ok {
		return CreateError("invalid currency code")
	}
	return nil
}

func IsValidTransactionType(transactionType string) error {
	if _, ok := ValidTransactionTypes[transactionType]; !ok {
		return CreateError("invalid transaction type")
	}
	return nil
}

func IsValidFrequency(frequency string) error {
	normalizedFrequency := strings.ToUpper(frequency)
	if _, ok := ValidFrequencies[normalizedFrequency]; !ok {
		return CreateError("invalid frequency")
	}
	return nil
}

func ValidateLendBorrowData(transactionData helpers.TransactionData) error {
	if transactionData.TransactionPartner == nil || *transactionData.TransactionPartner == "" {
		return CreateError("Transaction partner is required for lend/borrow transactions")
	}
	if transactionData.PaymentDueDate == nil || transactionData.PaymentDueDate.IsZero() {
		return CreateError("Payment due date is required for recurring expense transactions")
	}
	if transactionData.PaymentDueDate.Before(time.Now()) {
		return CreateError("Payment due date cannot be in the past")
	}
	return nil
}
