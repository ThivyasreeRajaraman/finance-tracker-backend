package userhelper

import (
	dbhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/common"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/dgrijalva/jwt-go"
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
