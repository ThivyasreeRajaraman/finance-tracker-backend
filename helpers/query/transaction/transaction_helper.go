package transactionhelpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	categoryhelpers "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/category"
	transactionpartnerhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/transaction_partner"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func FetchByID(c *gin.Context, transaction *models.Transaction, transactionId uint) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	if err := initializers.DB.Preload("Category").Preload("TransactionPartner").Where("user_id = ?", userID).First(&transaction, transactionId).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve existing transaction", err)
		return err
	}
	return nil
}
func HandleIncomeExpenseTransaction(userID uint, transaction *models.Transaction, transactionData helpers.TransactionData) error {
	if transactionData.CategoryName == nil || *transactionData.CategoryName == "" {
		return utils.CreateError("Category name is required for income/expense transactions")
	}

	category, err := categoryhelpers.GetOrCreateCategory(userID, transactionData.CategoryName, transactionData.TransactionType)
	if err != nil {
		return utils.CreateError("Failed to get or create category")
	}
	if err := utils.IsValidCurrency(transactionData.Currency); err != nil {
		return utils.CreateError("Invalid Currency Code")
	}
	transaction.CategoryID = &category.ID
	return nil
}

func HandleLendBorrowTransaction(userID uint, transaction *models.Transaction, transactionData helpers.TransactionData) error {
	if err := utils.ValidateLendBorrowData(transactionData); err != nil {
		return err
	}

	partner, err := transactionpartnerhelper.FetchOrCreate(userID, transactionData.TransactionPartner)
	if err != nil {
		return utils.CreateError("Failed to get transaction partner")
	}

	transaction.TransactionPartnerID = &partner.ID
	if err := transactionpartnerhelper.UpdateTransactionPartnerAmount(*transaction.TransactionPartnerID, transactionData.TransactionType, transactionData.Amount, transactionData.PaymentDueDate); err != nil {
		return err
	}
	if err := utils.IsValidCurrency(transactionData.Currency); err != nil {
		return utils.CreateError("Invalid Currency Code")
	}

	return nil
}

func CheckThreshold(c *gin.Context, transaction helpers.TransactionResponse, userID uint) string {
	var totalAmount uint
	var budget models.Budgets

	startOfMonth := time.Now().UTC().Truncate(time.Hour*24).AddDate(0, 0, -time.Now().Day()+1)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	if err := initializers.DB.
		Model(&models.Transaction{}).
		Select("SUM(amount) as totalAmount").
		Where("user_id = ? AND category_id = ? AND transaction_type = ? AND created_at BETWEEN ? AND ?", userID, transaction.CategoryID, "expense", startOfMonth, endOfMonth).
		Scan(&totalAmount).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve transactions", err)
	}

	err := initializers.DB.Model(&models.Budgets{}).
		Where("user_id = ? AND category_id = ?", userID, transaction.CategoryID).First(&budget).Error
	if err == nil {
		if totalAmount+transaction.Amount > budget.Threshold {
			return "Expense threshold reached for this category"
		}
	}
	return ""
}

func CalculateTotalAmounts(c *gin.Context) error {

	totalAmounts := make(map[string]uint)
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}
	userInterface, _ := c.Get("currentUser")
	user, ok := userInterface.(models.User)
	if !ok {
		err := utils.CreateError("invalid user data")
		utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return err
	}
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	rows, err := initializers.DB.Model(&models.Transaction{}).
		Select("transaction_type, SUM(amount) as total_amount, currency").
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startOfMonth, endOfMonth).
		Group("transaction_type, currency").
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for transactionType := range utils.ValidTransactionTypes {
		totalAmounts[transactionType] = 0
	}
	totalAmounts["budget"] = 0

	for rows.Next() {
		var transactionType, currency string
		var totalAmount uint
		if err := rows.Scan(&transactionType, &totalAmount, &currency); err != nil {
			return err
		}
		fmt.Println("val:", transactionType, currency, totalAmount)
		if currency != *user.DefaultCurrency {
			convertedAmount, convErr := convertCurrency(totalAmount, currency, user.DefaultCurrency)
			if convErr != nil {
				return convErr
			}
			totalAmount = convertedAmount
			fmt.Println("conv:", totalAmount)
		}

		if _, ok := utils.ValidTransactionTypes[transactionType]; ok {
			totalAmounts[transactionType] += totalAmount
			fmt.Println("sum:", totalAmounts[transactionType])
		}
	}

	// budgetRows, err := initializers.DB.Model(&models.Budgets{}).
	// 	Select("SUM(amount), currency").
	// 	Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startOfMonth, endOfMonth).
	// 	Group("currency").
	// 	Rows()
	// if err != nil {
	// 	return err
	// }
	// defer budgetRows.Close()

	// for budgetRows.Next() {
	// 	var totalBudgetAmount uint
	// 	var currency string
	// 	if err := budgetRows.Scan(&totalBudgetAmount, &currency); err != nil {
	// 		return err
	// 	}

	// 	if currency != *user.DefaultCurrency {
	// 		convertedBudgetAmount, convErr := convertCurrency(totalBudgetAmount, currency, user.DefaultCurrency)
	// 		if convErr != nil {
	// 			return convErr
	// 		}
	// 		totalBudgetAmount = convertedBudgetAmount
	// 	}

	// 	totalAmounts["budget"] += totalBudgetAmount
	// }

	utils.SendResponse(c, "Total fetched successfully", "transaction_total", totalAmounts)
	return nil
}

func CalculateCategoryWiseAmounts(c *gin.Context) error {
	totalAmounts := make(map[string]map[string]uint)
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}
	userInterface, _ := c.Get("currentUser")
	user, ok := userInterface.(models.User)
	if !ok {
		err := utils.CreateError("invalid user data")
		utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return err
	}
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
	totalAmounts["budget"] = make(map[string]uint)

	budgetAmounts, err := GetBudgetAmount(c, userID, startOfMonth, endOfMonth)
	if err != nil {
		return err
	}

	for categoryName, amount := range budgetAmounts {
		totalAmounts["budget"][categoryName] = amount
	}

	rows, err := initializers.DB.Model(&models.Transaction{}).
		Select("transaction_type, category_id, SUM(amount) as total_amount, currency").
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startOfMonth, endOfMonth).
		Group("transaction_type, category_id, currency").
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for transactionType := range utils.ValidTransactionTypes {
		totalAmounts[transactionType] = make(map[string]uint)
	}

	for rows.Next() {
		var transactionType string
		var categoryID *uint
		var totalAmount uint
		var currency string
		if err := rows.Scan(&transactionType, &categoryID, &totalAmount, &currency); err != nil {
			return err
		}
		if _, ok := utils.ValidTransactionTypes[transactionType]; ok && categoryID != nil {
			categoryName := GetCategoryName(c, categoryID)
			if categoryName != "" {
				if currency != *user.DefaultCurrency {
					convertedAmount, convErr := convertCurrency(totalAmount, currency, user.DefaultCurrency)
					if convErr != nil {
						return convErr
					}
					totalAmount = convertedAmount
				}
				totalAmounts[transactionType][categoryName] += totalAmount
			}
		}
	}

	utils.SendResponse(c, "Total fetched successfully", "transaction_total_by_category", totalAmounts)
	return nil
}

func GetCategoryName(c *gin.Context, categoryID *uint) string {
	var category models.Categories
	err := initializers.DB.Model(&models.Categories{}).Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve existing category", err)
		return ""
	}
	return category.Name
}

// func GetBudgetAmount(c *gin.Context, userID uint, startOfMonth time.Time, endOfMonth time.Time) (map[string]uint, error) {
// 	budgetAmounts := make(map[string]uint)
// 	rows, err := initializers.DB.Model(&models.Budgets{}).
// 		Select("category_id, SUM(amount) as total_amount, currency").
// 		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startOfMonth, endOfMonth).
// 		Group("category_id, currency").
// 		Rows()
// 	if err != nil {
// 		return budgetAmounts, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var categoryID uint
// 		var totalAmount uint
// 		var currency string
// 		userInterface, _ := c.Get("currentUser")
// 		user, ok := userInterface.(models.User)
// 		if !ok {
// 			err := utils.CreateError("invalid user data")
// 			utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
// 			return nil, err
// 		}
// 		if err := rows.Scan(&categoryID, &totalAmount, &currency); err != nil {
// 			return budgetAmounts, err
// 		}
// 		categoryName := GetCategoryName(c, &categoryID)
// 		if categoryName != "" {
// 			if currency != *user.DefaultCurrency {
// 				convertedAmount, convErr := convertCurrency(totalAmount, currency, user.DefaultCurrency)
// 				if convErr != nil {
// 					return budgetAmounts, convErr
// 				}
// 				totalAmount = convertedAmount
// 			}
// 			budgetAmounts[categoryName] = totalAmount
// 		}
// 	}

// 	return budgetAmounts, nil
// }

func GetBudgetAmount(c *gin.Context, userID uint, startOfMonth time.Time, endOfMonth time.Time) (map[string]uint, error) {
	budgetAmounts := make(map[string]uint)

	rows, err := initializers.DB.Model(&models.Budgets{}).
		Select("category_id, SUM(amount) as total_amount, currency").
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startOfMonth, endOfMonth).
		Group("category_id, currency").
		Rows()
	if err != nil {
		return budgetAmounts, err
	}
	defer rows.Close()

	userInterface, _ := c.Get("currentUser")
	user, ok := userInterface.(models.User)
	if !ok {
		err := utils.CreateError("invalid user data")
		utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return nil, err
	}

	for rows.Next() {
		var categoryID uint
		var totalAmount uint
		var currency string
		if err := rows.Scan(&categoryID, &totalAmount, &currency); err != nil {
			return budgetAmounts, err
		}
		categoryName := GetCategoryName(c, &categoryID)
		if categoryName != "" {
			if _, ok := budgetAmounts[categoryName]; !ok {
				budgetAmounts[categoryName] = 0
			}
			if currency != *user.DefaultCurrency {
				convertedAmount, convErr := convertCurrency(totalAmount, currency, user.DefaultCurrency)
				if convErr != nil {
					return budgetAmounts, convErr
				}
				totalAmount = uint(convertedAmount)
			}
			budgetAmounts[categoryName] += totalAmount
		}
	}

	return budgetAmounts, nil
}

func convertCurrency(amount uint, currency string, defaultCurrency *string) (uint, error) {
	apiKey := "93umfet0jubnfrs7dem1gju1bm5a0q1g8j4k63989o86co1jk7n33o"
	convertURL := fmt.Sprintf("https://anyapi.io/api/v1/exchange/convert?apiKey=%s&base=%s&to=%s&amount=%d", apiKey, currency, *defaultCurrency, amount)
	fmt.Println("args:base:", currency, "to:", *defaultCurrency, "amount:", amount)
	fmt.Println("url:", convertURL)
	resp, err := http.Get(convertURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Base            string  `json:"base"`
		To              string  `json:"to"`
		Amount          uint    `json:"amount"`
		ConvertedAmount float64 `json:"converted"`
		Rate            float64 `json:"rate"`
		LatestUpdate    int64   `json:"latestUpdate"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	fmt.Println("convAmount:", result.ConvertedAmount)
	fmt.Println("convAmount:", uint(result.ConvertedAmount))

	return uint(result.ConvertedAmount), nil
}
