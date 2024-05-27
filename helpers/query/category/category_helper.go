package categoryhelpers

import (
	"strings"

	dbhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/common"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
)

func GetOrCreateCategory(userID uint, categoryName *string, transactionType string) (*models.Categories, error) {
	var targetType string
	if transactionType == "budget" {
		targetType = "expense"
	} else if transactionType == "expense" {
		targetType = "budget"
	}
	var err error
	var existingCategory models.Categories
	if transactionType == "budget" || transactionType == "expense" {
		err = initializers.DB.Where("LOWER(name) = ? AND (type = ? OR type = ?) AND (user_id IS NULL OR user_id = ?)", strings.ToLower(*categoryName), transactionType, targetType, userID).First(&existingCategory).Error
	} else {
		// Income
		err = initializers.DB.Where("LOWER(name) = ? AND type = ? AND (user_id IS NULL OR user_id = ?)", strings.ToLower(*categoryName), transactionType, userID).First(&existingCategory).Error
	}

	if err != nil {
		// Category does not exist -> create a new one
		newCategory := models.Categories{
			Name:   *categoryName,
			Type:   transactionType,
			UserID: &userID,
		}
		if err = dbhelper.GenericCreate(&newCategory); err != nil {
			return nil, err
		}
		return &newCategory, nil
	}
	// Category already exists
	return &existingCategory, nil
}
