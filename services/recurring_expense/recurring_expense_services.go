package recurringexpenseservices

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	if recurringExpenseData.NextExpenseDate == "" {
		return utils.CreateError("Payment start date is required for recurring expense transactions")
	}
	if recurringExpenseData.Currency == "" {
		return utils.CreateError("Currency cannot be empty")
	}
	if err := utils.IsValidCurrency(recurringExpenseData.Currency); err != nil {
		return utils.CreateError("Invalid Currency Code")
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
		Currency:        recurringExpenseData.Currency,
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
	if recurringExpenseData.Frequency != nil {
		if *recurringExpenseData.Frequency == "" {
			return utils.CreateError("frequency cannot be empty")
		}
		err := utils.IsValidFrequency(*recurringExpenseData.Frequency)
		if err != nil {
			return err
		}
	}
	if recurringExpenseData.NextExpenseDate != nil {
		var formattedDate string
		parsedDate, err := time.Parse("2006-01-02", *recurringExpenseData.NextExpenseDate)
		if err != nil {
			return fmt.Errorf("failed to parse next_expense_date: %v", err)
		}
		fmt.Println("parsed date::", parsedDate)
		formattedDate = parsedDate.Format("2006-01-02")
		fmt.Println("formatted Date::", formattedDate)
		recurringExpenseData.NextExpenseDate = &formattedDate
		fmt.Println(" recc formatted Date::", recurringExpenseData.NextExpenseDate)
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
	if recurringExpenseData.Active != nil {
		existingExpense.Active = *recurringExpenseData.Active
	}
	fmt.Println("data bef::", existingExpense)

	if err := dbhelper.GenericUpdate(existingExpense); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update budget", err)
		return err
	}
	fmt.Println("data aft::", existingExpense)

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
	var reminders []map[string]interface{}

	var upcomingRecurringExpenses []models.RecurringExpense
	if err := initializers.DB.Model(&models.RecurringExpense{}).
		Where("user_id = ? AND active = ? AND next_expense_date BETWEEN ? AND ?", userID, true, time.Now().Format("2006-01-02"), time.Now().AddDate(0, 0, 5).Format("2006-01-02")).
		Preload("User").Preload("Category").
		Find(&upcomingRecurringExpenses).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch upcoming recurring expenses", err)
		return
	}

	updatePastRecurringExpenses(c, userID)

	for _, expense := range upcomingRecurringExpenses {
		nextExpenseDate, err := time.Parse("2006-01-02", expense.NextExpenseDate)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, "Failed to parse next expense date", err)
			continue
		}

		daysUntilExpense := int(time.Until(nextExpenseDate).Hours() / 24)
		if daysUntilExpense <= 5 {
			reminder := sendRecurringExpenseReminder(expense, daysUntilExpense)
			reminders = append(reminders, reminder)
		}
	}
	sortedReminders, err := sortReminders(reminders)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to sort expenses", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Reminder": sortedReminders})
}

func sendRecurringExpenseReminder(expense models.RecurringExpense, daysUntilExpense int) map[string]interface{} {
	var message string
	if daysUntilExpense == 0 {
		message = fmt.Sprintf("Your recurring expense for %s of %d %s is due today!",
			expense.Category.Name, expense.Amount, expense.Currency)
	} else {
		message = fmt.Sprintf("Your recurring expense for %s of %d %s is due in %d day(s)!",
			expense.Category.Name, expense.Amount, expense.Currency, daysUntilExpense)
	}

	return map[string]interface{}{"reminders": message, "id": expense.ID}
}

func GetExistingExpense(c *gin.Context, categoryName string) (uint, error) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.Zero, err
	}
	category, err := budgetservices.GetOrCreateCategory(c, userID, &categoryName, "recurringExpense")
	if err != nil {
		return utils.Zero, err
	}

	recurringExpenseID, err := recurringexpensehelper.GetExistingExpense(c, category.ID)
	if err != nil {
		return utils.Zero, err
	}
	return recurringExpenseID, nil
}

func sortReminders(reminders []map[string]interface{}) ([]map[string]interface{}, error) {
	sort.SliceStable(reminders, func(i, j int) bool {
		return getDueTime(reminders[i]["reminders"].(string)) < getDueTime(reminders[j]["reminders"].(string))
	})

	return reminders, nil
}

func getDueTime(reminder string) int {
	if strings.Contains(reminder, "today") {
		return 0
	}
	parts := strings.Split(reminder, " ")
	for i, part := range parts {
		if part == "day(s)!" || part == "day(s)," {
			dueTime, _ := strconv.Atoi(parts[i-1])
			fmt.Println("dueeeee", parts[i-1])
			return dueTime
		}
	}
	return -1
}

func updatePastRecurringExpenses(c *gin.Context, userID uint) {
	var pastRecurringExpenses []models.RecurringExpense
	fmt.Println("timeeee::", time.Now().Format("2006-01-02"))
	if err := initializers.DB.Model(&models.RecurringExpense{}).
		Where("user_id = ? AND active = ? AND next_expense_date < ?", userID, true, time.Now().Format("2006-01-02")).
		Find(&pastRecurringExpenses).Error; err != nil {
		fmt.Println("error here")
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update past expense date", err)
		return
	}
	fmt.Println("no errorr")

	for _, expense := range pastRecurringExpenses {
		nextExpenseDate, err := time.Parse("2006-01-02", expense.NextExpenseDate)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, "Failed to update past expense date", err)
			continue
		}
		var newNextExpenseDate time.Time
		if nextExpenseDate.Day() > time.Now().Day() {
			newNextExpenseDate = time.Date(time.Now().Year(), time.Now().Month(), nextExpenseDate.Day(), 0, 0, 0, 0, nextExpenseDate.Location())
		} else {
			newNextExpenseDate = time.Date(time.Now().Year(), time.Now().Month(), nextExpenseDate.Day(), 0, 0, 0, 0, nextExpenseDate.Location())
			newNextExpenseDate = newNextExpenseDate.AddDate(0, 1, 0)
		}

		if err := initializers.DB.Model(&expense).Update("next_expense_date", newNextExpenseDate.Format("2006-01-02")).Error; err != nil {
			utils.HandleError(c, http.StatusInternalServerError, "Failed to update past expense date", err)
			continue
		}
	}
}

func UpdateNextExpenseDate(c *gin.Context, existingExpense *models.RecurringExpense) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	nextExpenseDate, err := time.Parse("2006-01-02", existingExpense.NextExpenseDate)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to parse next expense date", err)
		return err
	}
	switch existingExpense.Frequency {
	case "MONTHLY":
		existingExpense.NextExpenseDate = nextExpenseDate.AddDate(0, 1, 0).Format("2006-01-02")
	case "WEEKLY":
		existingExpense.NextExpenseDate = nextExpenseDate.AddDate(0, 0, 7).Format("2006-01-02")
	case "YEARLY":
		existingExpense.NextExpenseDate = nextExpenseDate.AddDate(1, 0, 0).Format("2006-01-02")
	}
	if err := dbhelper.GenericUpdate(existingExpense); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update budget", err)
		return err
	}

	transaction := models.Transaction{
		UserID:          userID,
		TransactionType: "recurringExpense",
		CategoryID:      existingExpense.CategoryID,
		Amount:          existingExpense.Amount,
		Currency:        existingExpense.Currency,
	}
	if err := dbhelper.GenericCreate(&transaction); err != nil {
		return err
	}
	return nil
}
