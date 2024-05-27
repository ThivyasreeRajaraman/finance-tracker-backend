package budgetcontrollers

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	budgetservices "github.com/Thivyasree-Rajaraman/finance-tracker/services/budget"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

type BudgetController struct{}

type BudgetControllerInterface interface {
	Create(c *gin.Context)
	Fetch(c *gin.Context, createdSuccessfully bool)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func GetBudgetControllerInstance() BudgetControllerInterface {
	return new(BudgetController)
}

func (controller *BudgetController) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return
	}

	var budgetData []helpers.BudgetData
	if err := budgetservices.UnmarshalAndValidate(c, &budgetData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
		return
	}

	if err := budgetservices.CreateBudgets(c, budgetData, userID); err != nil {
		return
	}

	controller.Fetch(c, true)
}

func (controller *BudgetController) Fetch(c *gin.Context, createdSuccessfully bool) {

	var budgets []models.Budgets
	if err := budgetservices.Fetch(c, &budgets); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch budget", err)
	}

	var message string
	if createdSuccessfully {
		message = "Budget created successfully"
	} else {
		message = "Budget retrieved successfully"
	}

	budgetResponses, err := utils.CreateBudgetResponse(budgets)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to construct budget response", err)
	}

	utils.SendResponse(c, message, "budget", budgetResponses)
}

func (controller *BudgetController) Update(c *gin.Context) {
	var budgetData helpers.BudgetData
	if err := budgetservices.UnmarshalAndValidateSingleEntity(c, &budgetData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
		return
	}

	existingBudgetInterface, _ := c.Get("existingBudget")
	existingBudget, ok := existingBudgetInterface.(models.Budgets)
	if !ok {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get existing budgettt", nil)
		return
	}

	if err := budgetservices.Update(c, &existingBudget, budgetData); err != nil {
		return
	}
}

func (controller *BudgetController) Delete(c *gin.Context) {
	existingBudgetInterface, _ := c.Get("existingBudget")
	existingBudget, ok := existingBudgetInterface.(models.Budgets)
	if !ok {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get existing budget", nil)
		return
	}

	// Soft delete the transaction
	if err := budgetservices.Delete(c, &existingBudget); err != nil {
		return
	}
}
