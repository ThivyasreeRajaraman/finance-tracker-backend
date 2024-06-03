package recurringexpensecontrollers

import (
	"net/http"
	"sort"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	recurringexpenseservices "github.com/Thivyasree-Rajaraman/finance-tracker/services/recurring_expense"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

type RecurringExpenseController struct{}

type RecurringExpenseControllerInterface interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Fetch(c *gin.Context)
	Remind(c *gin.Context)
	FetchFrequencies(c *gin.Context)
}

func GetRecurringExpenseControllerInstance() RecurringExpenseControllerInterface {
	return new(RecurringExpenseController)
}

func (controller *RecurringExpenseController) Create(c *gin.Context) {
	var recurringExpenseData helpers.RecurringExpenseData
	if err := recurringexpenseservices.UnmarshalAndValidate(c, &recurringExpenseData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
		return
	}

	if err := recurringexpenseservices.Create(c, recurringExpenseData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to create recurring expense", err)
		return
	}
}

func (controller *RecurringExpenseController) Update(c *gin.Context) {
	var recurringExpenseData helpers.UpdateRecurringExpenseData
	if err := recurringexpenseservices.UnmarshalAndValidateForUpdate(c, &recurringExpenseData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
		return
	}

	if err := recurringexpenseservices.Update(c, recurringExpenseData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to update recurring expense", err)
		return
	}

}

func (controller *RecurringExpenseController) Delete(c *gin.Context) {
	existingExpense, err := recurringexpenseservices.GetExpenseFromPathParam(c)
	if err != nil {
		return
	}
	if err := recurringexpenseservices.Delete(c, existingExpense); err != nil {
		return
	}
}

func (controller *RecurringExpenseController) Fetch(c *gin.Context) {
	recurringExpenseModel := new(models.RecurringExpense)
	userID, err := utils.GetUserID(c)
	if err != nil {
		return
	}
	conditions := map[string]interface{}{
		"user_id": userID,
	}
	if data := utils.List(c, recurringExpenseModel, conditions); data != nil {
		return
	}
}

func (controller *RecurringExpenseController) Remind(c *gin.Context) {
	recurringexpenseservices.SendRecurringExpenseReminders(c)
}

func (controller *RecurringExpenseController) FetchFrequencies(c *gin.Context) {
	frequencies := make([]string, 0, len(utils.ValidFrequencies))
	for frequency := range utils.ValidFrequencies {
		frequencies = append(frequencies, frequency)
	}
	sort.Strings(frequencies)
	c.JSON(http.StatusOK, gin.H{"frequencies": frequencies})
}
