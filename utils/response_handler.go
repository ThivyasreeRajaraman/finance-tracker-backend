package utils

import (
	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
)

func CreateBudgetResponse(budgets []models.Budgets) ([]helpers.BudgetResponse, error) {
	budgetResponses := make([]helpers.BudgetResponse, 0)
	for _, budget := range budgets {
		var defaultCurrency string
		if budget.User.DefaultCurrency != nil {
			defaultCurrency = *budget.User.DefaultCurrency
		}
		budgetResponse := helpers.BudgetResponse{
			ID:              budget.ID,
			UserID:          budget.UserID,
			Name:            budget.User.Name,
			Email:           budget.User.Email,
			DefaultCurrency: defaultCurrency,
			CategoryID:      budget.CategoryID,
			CategoryName:    budget.Category.Name,
			Amount:          budget.Amount,
			Threshold:       budget.Threshold,
		}
		budgetResponses = append(budgetResponses, budgetResponse)
	}
	return budgetResponses, nil
}

func CreateTransactionResponse(transaction models.Transaction) (helpers.TransactionResponse, error) {

	var defaultCurrency string
	var categoryID uint
	if transaction.CategoryID != nil {
		categoryID = *transaction.CategoryID
	}
	if transaction.User.DefaultCurrency != nil {
		defaultCurrency = *transaction.User.DefaultCurrency
	}
	transactionResponse := helpers.TransactionResponse{
		ID:                 transaction.ID,
		UserID:             transaction.UserID,
		Name:               transaction.User.Name,
		DefaultCurrency:    defaultCurrency,
		CategoryID:         &categoryID,
		CategoryName:       &transaction.Category.Name,
		Amount:             transaction.Amount,
		TransactionType:    transaction.TransactionType,
		TransactionPartner: &transaction.TransactionPartner.PartnerName,
	}
	return transactionResponse, nil
}

func CreatePartnerResponse(partners []models.TransactionPartner) ([]helpers.TransactionPartnerResponse, error) {
	partnerResponses := make([]helpers.TransactionPartnerResponse, 0)
	for _, partner := range partners {
		var transaction string
		if partner.ClosingBalance > 0 {
			transaction = "Borrowed"
		} else if partner.ClosingBalance < 0 {
			transaction = "Lent"
		} else {
			transaction = "Nil"
		}
		partnerResponse := helpers.TransactionPartnerResponse{
			PartnerName:     partner.PartnerName,
			Amount:          abs(partner.ClosingBalance),
			TransactionType: transaction,
			DueDate:         partner.DueDate,
		}
		partnerResponses = append(partnerResponses, partnerResponse)
	}
	return partnerResponses, nil
}
