package userhelper

import (
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/dgrijalva/jwt-go"
)

func CreateUser(user *models.User) error {
	return initializers.DB.Create(&user).Error
}

func Update(user *models.User) error {
	return initializers.DB.Save(user).Error
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
