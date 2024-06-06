package userhelper

import (
	"net/http"

	dbhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/common"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Create(user *models.User) error {
	return dbhelper.GenericCreate(user)
}

func Update(user *models.User) error {
	return dbhelper.GenericUpdate(user)
}

func SearchByEmail(user *models.User) error {
	return initializers.DB.Where("email = ?", user.Email).First(user).Error
}

func FetchUserByClaims(claims jwt.MapClaims) (models.User, error) {
	var user models.User

	query := initializers.DB.Model(&user)

	if userIDFloat, ok := claims["user_id"].(float64); ok {
		userID := uint(userIDFloat)
		query = query.Where("id = ?", userID)
	}

	if email, ok := claims["email"].(string); ok {
		query = query.Where("email = ?", email)
	}

	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func FetchCategories(c *gin.Context, userID uint) error {
	var categories []models.Categories
	err := initializers.DB.Model(&models.Categories{}).
		Select("name").Where("user_id = ? OR user_id IS NULL", userID).
		Find(&categories).Error
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch categories", err)
		return err
	}
	names := make([]string, len(categories))
	for i, category := range categories {
		names[i] = category.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"data": names,
	})
	return nil
}
