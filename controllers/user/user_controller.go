package usercontrollers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	userservices "github.com/Thivyasree-Rajaraman/finance-tracker/services/user"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func Update(c *gin.Context) {
	userInterface, _ := c.Get("currentUser")
	user, err := userInterface.(models.User)
	if !err {
		utils.HandleError(c, http.StatusBadRequest, "Invalid user data", nil)
	}

	if err := utils.UnmarshalData(c, &user); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
		return
	}

	if err := userservices.UpdateUser(c, &user); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	utils.SendResponse(c, "User updated successfully", "user", user)
}

func FetchCurrencies(c *gin.Context) {
	currencies := make([]string, 0, len(utils.ValidCurrencies))
	for currency := range utils.ValidCurrencies {
		currencies = append(currencies, currency)
	}
	sort.Strings(currencies)
	c.JSON(http.StatusOK, gin.H{"currencies": currencies})
}

func Fetch(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return
	}
	userModel := new(models.User)
	conditions := map[string]interface{}{
		"id": userID,
	}
	if data := utils.List(c, userModel, conditions, nil, nil, nil, ""); data != nil {
		return
	}
}

func FetchCategories(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	transactionType := c.Param("transactionType")
	fmt.Println("transaction:", transactionType)
	if err != nil {
		return
	}
	if err := userservices.FetchCategories(c, userID, transactionType); err != nil {
		return
	}

}
