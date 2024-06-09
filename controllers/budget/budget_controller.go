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
	Fetch(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	UnitFetch(c *gin.Context)
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

	controller.FetchCreatedBudget(c)
}

func (controller *BudgetController) FetchCreatedBudget(c *gin.Context) {

	var budgets []models.Budgets
	if err := budgetservices.Fetch(c, &budgets); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch budget", err)
	}

	budgetResponses, err := utils.CreateBudgetResponse(budgets)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to construct budget response", err)
	}

	utils.SendResponse(c, "Budget created successfully", "budget", budgetResponses)
}

func (controller *BudgetController) Fetch(c *gin.Context) {
	budgetModel := new(models.Budgets)
	userID, err := utils.GetUserID(c)
	if err != nil {
		return
	}
	conditions := map[string]interface{}{
		"user_id": userID,
	}
	if data := utils.List(c, budgetModel, conditions, "id ASC"); data != nil {
		return
	}
}

func (controller *BudgetController) Update(c *gin.Context) {
	existingBudget, err := budgetservices.GetBudgetFromPathParam(c)
	if err != nil {
		return
	}
	var budgetData helpers.BudgetData
	if err := budgetservices.UnmarshalAndValidateSingleEntity(c, &budgetData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
		return
	}

	if err := budgetservices.Update(c, existingBudget, budgetData); err != nil {
		return
	}
}

func (controller *BudgetController) Delete(c *gin.Context) {
	existingBudget, err := budgetservices.GetBudgetFromPathParam(c)
	if err != nil {
		return
	}

	// Soft delete the transaction
	if err := budgetservices.Delete(c, existingBudget); err != nil {
		return
	}
}

func (controller *BudgetController) UnitFetch(c *gin.Context) {
	existingBudget, err := budgetservices.GetBudgetFromPathParam(c)
	if err != nil {
		return
	}
	budgetResponses, err := utils.CreateBudgetResponse([]models.Budgets{*existingBudget})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to construct budget response", err)
	}
	utils.SendResponse(c, "Budget retrieved successfully", "budget", budgetResponses)
}
