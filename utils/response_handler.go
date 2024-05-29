package utils

import (
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

func List[T any](c *gin.Context, model *T, conditions map[string]interface{}) []T {
	page, err := strconv.Atoi(c.Query("pageNumber"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("perPageCount"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	data, totalCount, err := dbhelper.FetchDataWithPagination(model, page, limit, conditions)
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
