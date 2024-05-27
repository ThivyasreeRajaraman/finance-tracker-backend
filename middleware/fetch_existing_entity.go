package middleware

import (
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	budgetservices "github.com/Thivyasree-Rajaraman/finance-tracker/services/budget"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func GetBudgetFromPathParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		budgetID, err := utils.ParseUintParam(c, "budgetId")
		if err != nil {
			return
		}
		var existingBudget models.Budgets
		if err := budgetservices.FetchBudgetById(c, &existingBudget, budgetID); err != nil {
			return
		}
		c.Set("existingBudget", existingBudget)
	}
}
