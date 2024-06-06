package recurringexpensehelper

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func FetchByID(c *gin.Context, existingExpense *models.RecurringExpense, expenseID uint) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	if err := initializers.DB.Preload("Category").Where("user_id = ?", userID).First(existingExpense, expenseID).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve existing recurring expense", err)
		return err
	}
	return nil
}
