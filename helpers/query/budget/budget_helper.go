package budgethelpers

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func Fetch(c *gin.Context, budgets *[]models.Budgets) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	if err := initializers.DB.Preload("User").Preload("Category").Where("user_id = ?", userID).Find(&budgets).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to load budgets", err)
		return err
	}
	return nil
}

func FetchByID(c *gin.Context, budget *models.Budgets, budgetID uint) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	if err := initializers.DB.Where("user_id = ?", userID).First(&budget, budgetID).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve existing budget", err)
		return err
	}
	return nil
}
