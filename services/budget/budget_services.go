package budgetservices

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	budgethelpers "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/budget"
	categoryhelpers "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/category"
	dbhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/common"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func UnmarshalAndValidate(c *gin.Context, budgetData *[]helpers.BudgetData) error {
	if err := utils.UnmarshalData(c, budgetData); err != nil {
		return err
	}

	for _, d := range *budgetData {
		if d.CategoryName == nil || *d.CategoryName == "" {
			return utils.CreateError("category_name cannot be empty")
		}
		if d.Amount == nil || *d.Amount <= 0 {
			return utils.CreateError("amount must be greater than zero")
		}
		if d.Threshold == nil || *d.Threshold <= 0 {
			return utils.CreateError("threshold must be greater than zero")
		}
		if d.Threshold != nil && d.Amount != nil && *d.Threshold > *d.Amount {
			return utils.CreateError("threshold must be less than budget amount")
		}
	}
	return nil
}

func UnmarshalAndValidateSingleEntity(c *gin.Context, budgetData *helpers.BudgetData) error {
	if err := utils.UnmarshalData(c, budgetData); err != nil {
		return err
	}
	if budgetData.CategoryName != nil && *budgetData.CategoryName == "" {
		return utils.CreateError("category_name cannot be empty")
	}
	if budgetData.Amount != nil && *budgetData.Amount <= 0 {
		return utils.CreateError("amount must be greater than zero")
	}
	if budgetData.Threshold != nil && *budgetData.Threshold <= 0 {
		return utils.CreateError("threshold must be greater than zero")
	}
	return nil
}

func GetOrCreateCategory(c *gin.Context, userID uint, categoryName *string, transactionType string) (*models.Categories, error) {
	category, err := categoryhelpers.GetOrCreateCategory(userID, categoryName, transactionType)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func CreateBudgets(c *gin.Context, budgetData []helpers.BudgetData, userID uint) error {
	for _, data := range budgetData {
		category, err := GetOrCreateCategory(c, userID, data.CategoryName, "budget")
		if err != nil {
			return err
		}

		budget := models.Budgets{
			UserID:     userID,
			CategoryID: category.ID,
			Amount:     *data.Amount,
			Threshold:  *data.Threshold,
		}
		if err := Create(c, &budget); err != nil {
			return err
		}
	}
	return nil
}

func Create(c *gin.Context, budget *models.Budgets) error {
	if err := dbhelper.GenericCreate(budget); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to create budget", err)
		return err
	}
	return nil
}

func Fetch(c *gin.Context, budgets *[]models.Budgets) error {
	if err := budgethelpers.Fetch(c, budgets); err != nil {
		return err
	}
	return nil
}

func FetchBudgetById(c *gin.Context, budget *models.Budgets, budgetID uint) error {
	if err := budgethelpers.FetchByID(c, budget, budgetID); err != nil {
		return err
	}
	return nil
}

func preloadBudgetAssociations(c *gin.Context, budget *models.Budgets) error {
	if err := initializers.DB.Preload("Category").Preload("User").First(budget, budget.ID).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to preload user and category association", err)
		return err
	}
	return nil
}

func UpdateExistingBudget(c *gin.Context, existingBudget *models.Budgets, budgetData helpers.BudgetData, categoryID uint) error {
	if categoryID != 0 && existingBudget.CategoryID != categoryID {
		existingBudget.CategoryID = categoryID
	}
	if budgetData.Amount != nil {
		existingBudget.Amount = *budgetData.Amount
	}
	if budgetData.Threshold != nil {
		existingBudget.Threshold = *budgetData.Threshold
	}

	if err := dbhelper.GenericUpdate(existingBudget); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update budget", err)
		return err
	}

	if err := preloadBudgetAssociations(c, existingBudget); err != nil {
		return err
	}

	utils.SendResponse(c, "Budget updated successfully", "budget", existingBudget)
	return nil
}

func Update(c *gin.Context, existingBudget *models.Budgets, budgetData helpers.BudgetData) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}
	var category *models.Categories
	if budgetData.CategoryName != nil {
		category, err = GetOrCreateCategory(c, userID, budgetData.CategoryName, "budget")
		if err != nil {
			return err
		}
	}
	// Update the existing budget
	if err := UpdateExistingBudget(c, existingBudget, budgetData, category.ID); err != nil {
		return err
	}
	return nil
}

func Delete(c *gin.Context, existingBudget *models.Budgets) error {
	if err := dbhelper.GenericDelete(existingBudget); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to delete budget", err)
		return err
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Budget deleted successfully"})
	return nil
}

func GetBudgetFromPathParam(c *gin.Context) (*models.Budgets, error) {
	budgetID, err := utils.ParseUintParam(c, "budgetId")
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Invalid budget id", err)
		return nil, err
	}
	var existingBudget models.Budgets
	if err := FetchBudgetById(c, &existingBudget, budgetID); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch budget", err)
		return nil, err
	}
	return &existingBudget, nil
}
