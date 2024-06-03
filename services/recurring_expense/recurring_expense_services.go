package recurringexpenseservices

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	dbhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/common"
	recurringexpensehelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/recurring_expense"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	budgetservices "github.com/Thivyasree-Rajaraman/finance-tracker/services/budget"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func UnmarshalAndValidate(c *gin.Context, recurringExpenseData *helpers.RecurringExpenseData) error {
	if err := utils.UnmarshalData(c, recurringExpenseData); err != nil {
		return err
	}

	if recurringExpenseData.CategoryName == "" {
		return utils.CreateError("category_name cannot be empty")
	}
	if recurringExpenseData.Amount <= 0 {
		return utils.CreateError("amount must be greater than zero")
	}
	if recurringExpenseData.Frequency == "" {
		return utils.CreateError("frequency cannot be empty")
	}
	if recurringExpenseData.NextExpenseDate.IsZero() {
		return utils.CreateError("Payment start date is required for recurring expense transactions")
	}
	if recurringExpenseData.NextExpenseDate.Before(time.Now()) {
		return utils.CreateError("Payment start date cannot be in the past")
	}
	if err := utils.IsValidFrequency(recurringExpenseData.Frequency); err != nil {
		return err
	}
	return nil
}

func Create(c *gin.Context, recurringExpenseData helpers.RecurringExpenseData) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}
	category, err := budgetservices.GetOrCreateCategory(c, userID, &recurringExpenseData.CategoryName, "recurringExpense")
	if err != nil {
		return err
	}

	recurringExpense := models.RecurringExpense{
		UserID:          userID,
		CategoryID:      category.ID,
		Amount:          recurringExpenseData.Amount,
		Frequency:       recurringExpenseData.Frequency,
		NextExpenseDate: recurringExpenseData.NextExpenseDate,
	}
	if err := dbhelper.GenericCreate(&recurringExpense); err != nil {
		return err
	}

	utils.SendResponse(c, "Recurring expense created successfully", "recurring_expense", recurringExpense)
	return nil
}

func UnmarshalAndValidateForUpdate(c *gin.Context, recurringExpenseData *helpers.UpdateRecurringExpenseData) error {
	if err := utils.UnmarshalData(c, recurringExpenseData); err != nil {
		return err
	}

	if recurringExpenseData.CategoryName != nil && *recurringExpenseData.CategoryName == "" {
		return utils.CreateError("category_name cannot be empty")
	}
	if recurringExpenseData.Amount != nil && *recurringExpenseData.Amount <= 0 {
		return utils.CreateError("amount must be greater than zero")
	}
	if recurringExpenseData.NextExpenseDate != nil && (*recurringExpenseData.NextExpenseDate).Before(time.Now()) {
		return utils.CreateError("Payment start date cannot be in the past")
	}
	if recurringExpenseData.Frequency != nil {
		if *recurringExpenseData.Frequency == "" {
			return utils.CreateError("frequency cannot be empty")
		}
		err := utils.IsValidFrequency(*recurringExpenseData.Frequency)
		if err != nil {
			return err
		}
	}
	return nil
}

func Update(c *gin.Context, recurringExpenseData helpers.UpdateRecurringExpenseData) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}
	existingExpense, err := GetExpenseFromPathParam(c)
	if err != nil {
		return err
	}

	var category *models.Categories
	if recurringExpenseData.CategoryName != nil {
		category, err = budgetservices.GetOrCreateCategory(c, userID, recurringExpenseData.CategoryName, "recurringExpense")
		if err != nil {
			return err
		}
		existingExpense.CategoryID = category.ID
	}
	if recurringExpenseData.Amount != nil {
		existingExpense.Amount = *recurringExpenseData.Amount
	}
	if recurringExpenseData.Frequency != nil {
		existingExpense.Frequency = *recurringExpenseData.Frequency
	}
	if recurringExpenseData.NextExpenseDate != nil {
		existingExpense.NextExpenseDate = *recurringExpenseData.NextExpenseDate
	}

	if err := dbhelper.GenericUpdate(existingExpense); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update budget", err)
		return err
	}

	if err := initializers.DB.Preload("Category").Preload("User").First(existingExpense, existingExpense.ID).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to preload user and category association", err)
		return err
	}

	utils.SendResponse(c, "Recurring expense updated successfully", "recurring expense", existingExpense)
	return nil
}

func GetExpenseFromPathParam(c *gin.Context) (*models.RecurringExpense, error) {
	expenseID, err := utils.ParseUintParam(c, "recurringExpenseId")
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Invalid recurring expense id", err)
		return nil, err
	}
	var existingExpense models.RecurringExpense
	if err := FetchById(c, &existingExpense, expenseID); err != nil {
		return nil, err
	}
	return &existingExpense, nil
}

func FetchById(c *gin.Context, existingExpense *models.RecurringExpense, expenseID uint) error {
	if err := recurringexpensehelper.FetchByID(c, existingExpense, expenseID); err != nil {
		return err
	}
	return nil
}

func Delete(c *gin.Context, existingExpense *models.RecurringExpense) error {
	if err := dbhelper.GenericDelete(existingExpense); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to delete recurring expense", err)
		return err
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Recurring Expense deleted successfully"})
	return nil
}

func SendRecurringExpenseReminders(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return
	}
	var upcomingRecurringExpenses []models.RecurringExpense
	if err := initializers.DB.Model(&models.RecurringExpense{}).
		Where("user_id = ? AND next_expense_date BETWEEN ? AND ?", userID, time.Now(), time.Now().AddDate(0, 0, 5)).
		Preload("User").Preload("Category").
		Find(&upcomingRecurringExpenses).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch upcoming recurring expenses", err)
		return
	}

	for _, expense := range upcomingRecurringExpenses {
		daysUntilExpense := int(time.Until(expense.NextExpenseDate).Hours() / 24)
		if daysUntilExpense <= 5 {
			sendRecurringExpenseReminder(c, expense, daysUntilExpense)
		}
		if expense.NextExpenseDate.Format("2006-01-02") == time.Now().Format("2006-01-02") {
			switch expense.Frequency {
			case "MONTHLY":
				expense.NextExpenseDate = expense.NextExpenseDate.AddDate(0, 1, 0)
			case "WEEKLY":
				expense.NextExpenseDate = expense.NextExpenseDate.AddDate(0, 0, 7)
			case "YEARLY":
				expense.NextExpenseDate = expense.NextExpenseDate.AddDate(1, 0, 0)
			}
			if err := dbhelper.GenericUpdate(&expense); err != nil {
				utils.HandleError(c, http.StatusInternalServerError, "Failed to update next expense date", err)
				continue
			}
		}
	}
}

func sendRecurringExpenseReminder(c *gin.Context, expense models.RecurringExpense, daysUntilExpense int) {
	var message string
	fmt.Println("\n\ndata::", daysUntilExpense)
	if daysUntilExpense == 0 {
		message = fmt.Sprintf("Your recurring expense for %s of %d %s is due today.",
			expense.Category.Name, expense.Amount, *expense.User.DefaultCurrency)
	} else {
		fmt.Printf("data::%+v", expense)
		message = fmt.Sprintf("Your recurring expense for %s of %d %s is due in %d day(s).",
			expense.Category.Name, expense.Amount, *expense.User.DefaultCurrency, daysUntilExpense)
	}

	c.JSON(http.StatusOK, gin.H{"Reminder": message})
}
