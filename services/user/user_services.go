package userservices

import (
	userhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/user"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func UpdateUser(c *gin.Context, user *models.User) error {
	// validate currency code
	if err := utils.IsValidCurrency(*user.DefaultCurrency); err != nil {
		return err
	}
	if err := userhelper.Update(user); err != nil {
		return err
	}

	// set user data in context
	c.Set("currentUser", user)
	return nil
}
