package transactionhelpers

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func Fetch(c *gin.Context, transaction *models.Transaction) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	if err := initializers.DB.Preload("User").Preload("Category").Where("user_id = ?", userID).Find(&transaction).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to load transaction", err)
		return err
	}
	return nil
}

func FetchByID(c *gin.Context, transaction *models.Transaction, transactionId uint) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	if err := initializers.DB.Where("user_id = ?", userID).First(&transaction, transactionId).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve existing transaction", err)
		return err
	}
	return nil
}
