package usercontrollers

import (
	"net/http"

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
