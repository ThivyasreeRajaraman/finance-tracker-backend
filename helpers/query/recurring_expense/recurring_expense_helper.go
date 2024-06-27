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

func GetExistingExpense(c *gin.Context, categoryID uint) (uint, error) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.Zero, err
	}

	var recurringExpense models.RecurringExpense
	if err := initializers.DB.Preload("Category").Where("user_id = ? AND category_id = ?", userID, categoryID).First(&recurringExpense).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve existing recurring expense", err)
		return utils.Zero, err
	}
	return recurringExpense.ID, nil
}
