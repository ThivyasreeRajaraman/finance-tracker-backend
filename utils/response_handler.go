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
