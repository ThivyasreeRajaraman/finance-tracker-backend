package recurringexpensecontrollers

import (
	"fmt"
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
	FetchSingleEntity(c *gin.Context)
	Remind(c *gin.Context)
	FetchFrequencies(c *gin.Context)
	UpdateNextExpenseDate(c *gin.Context)
}

func GetRecurringExpenseControllerInstance() RecurringExpenseControllerInterface {
	return new(RecurringExpenseController)
}

func (controller *RecurringExpenseController) Create(c *gin.Context) {
	var recurringExpenseData helpers.RecurringExpenseData
	if err := recurringexpenseservices.UnmarshalAndValidate(c, &recurringExpenseData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
	}
	if err := recurringexpenseservices.Create(c, recurringExpenseData); err != nil {
		fmt.Println("err.Error()::", err.Error())
		if err.Error() == `ERROR: duplicate key value violates unique constraint "idx_category_id_user_id" (SQLSTATE 23505)` {
			fmt.Println("yessss")
			if expenseID, getExistingErr := recurringexpenseservices.GetExistingExpense(c, recurringExpenseData.CategoryName); expenseID != 0 {
				fmt.Println("second yess")
				c.JSON(http.StatusOK, gin.H{
					"success":     true,
					"status code": 200,
					"error":       "Recurring expense for the same category exists already!",
					"details":     err.Error(),
					"existingId":  expenseID,
				})
			} else if getExistingErr != nil {
				utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve existing expense", getExistingErr)
			}
		} else {
			fmt.Println("noooo")
			utils.HandleError(c, http.StatusBadRequest, "Failed to create recurring expense", err)
		}
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
	if data := utils.List(c, recurringExpenseModel, conditions, nil, nil, nil, "next_expense_date ASC"); data != nil {
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

func (controller *RecurringExpenseController) FetchSingleEntity(c *gin.Context) {
	existingExpense, err := recurringexpenseservices.GetExpenseFromPathParam(c)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch expense", err)
		return
	}

	response := helpers.RecurringExpenseResponse{
		Category:        existingExpense.Category.Name,
		Amount:          existingExpense.Amount,
		Frequency:       existingExpense.Frequency,
		NextExpenseDate: existingExpense.NextExpenseDate,
		Currency:        existingExpense.Currency,
	}

	utils.SendResponse(c, "Expense fetched successfully", "expense", response)
}

func (controller *RecurringExpenseController) UpdateNextExpenseDate(c *gin.Context) {
	existingExpense, err := recurringexpenseservices.GetExpenseFromPathParam(c)
	if err != nil {
		return
	}
	if err := recurringexpenseservices.UpdateNextExpenseDate(c, existingExpense); err != nil {
		return
	}
}
