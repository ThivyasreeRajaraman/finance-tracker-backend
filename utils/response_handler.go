package utils

import (
	"math"

	"net/http"
	"strconv"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	dbhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/common"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/gin-gonic/gin"
)

func CreateBudgetResponse(budgets []models.Budgets) ([]helpers.BudgetResponse, error) {
	budgetResponses := make([]helpers.BudgetResponse, 0)
	for _, budget := range budgets {

		budgetResponse := helpers.BudgetResponse{
			ID:           budget.ID,
			UserID:       budget.UserID,
			Name:         budget.User.Name,
			Email:        budget.User.Email,
			CategoryID:   budget.CategoryID,
			CategoryName: budget.Category.Name,
			Amount:       budget.Amount,
			Threshold:    budget.Threshold,
			Currency:     budget.Currency,
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
		Currency:           transaction.Currency,
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
			Amount:          uint(math.Abs(float64(partner.ClosingBalance))),
			TransactionType: transaction,
			DueDate:         partner.DueDate,
		}
		partnerResponses = append(partnerResponses, partnerResponse)
	}
	return partnerResponses, nil
}

func List[T any](c *gin.Context, model *T, conditions, unequalConditions, greaterThanConditions, lesserThanConditions map[string]interface{}, orderBy string) []T {
	page, err := strconv.Atoi(c.Query("pageNumber"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("perPageCount"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	data, totalCount, err := dbhelper.FetchDataWithPagination(model, page, limit, conditions, unequalConditions, greaterThanConditions, lesserThanConditions, orderBy)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to fetch data", err)
		return nil
	}

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data":       data,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
		"totalCount": totalCount,
	})
	return data
}
