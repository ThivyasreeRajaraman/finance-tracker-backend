package transactionhelpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

	if err := initializers.DB.Where("user_id = ?", userID).First(&transaction, transactionId).Error; err != nil {
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
	transaction.CategoryID = category.ID
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

	totalAmounts := make(map[string]float64)
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	monthStr := c.Query("month")
	yearStr := c.Query("year")
	fmt.Println("month,year:", monthStr, yearStr)
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		month = int(time.Now().Month())
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = time.Now().Year()
	}

	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	rows, err := initializers.DB.Model(&models.Transaction{}).
		Select("transaction_type, SUM(converted_amount) as total_amount").
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startOfMonth, endOfMonth).
		Group("transaction_type").
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for transactionType := range utils.ValidTransactionTypes {
		if transactionType != "recurringExpense" {
			totalAmounts[transactionType] = 0
		}
	}

	for rows.Next() {
		var transactionType string
		var totalAmount float64
		if err := rows.Scan(&transactionType, &totalAmount); err != nil {
			return err
		}
		fmt.Println("val:", transactionType, totalAmount)

		if _, ok := utils.ValidTransactionTypes[transactionType]; ok {
			totalAmounts[transactionType] = totalAmount
			fmt.Println("sum:", totalAmounts[transactionType])
		}
	}

	utils.SendResponse(c, "Total fetched successfully", "transaction_total", totalAmounts)
	return nil
}

func CalculateCategoryWiseAmounts(c *gin.Context) error {
	totalAmounts := make(map[string]map[string]float64)
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	monthStr := c.Query("month")
	yearStr := c.Query("year")
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		month = int(time.Now().Month())
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = time.Now().Year()
	}

	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	totalAmounts["budget"] = make(map[string]float64)

	budgetAmounts, err := GetBudgetAmount(c, userID)
	if err != nil {
		return err
	}

	for categoryName, amount := range budgetAmounts {
		totalAmounts["budget"][categoryName] = amount
	}

	rows, err := initializers.DB.Model(&models.Transaction{}).
		Select("transaction_type, category_id, SUM(converted_amount) as total_amount").
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startOfMonth, endOfMonth).
		Group("transaction_type, category_id").
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for transactionType := range utils.ValidTransactionTypes {
		totalAmounts[transactionType] = make(map[string]float64)
	}

	for rows.Next() {
		var transactionType string
		var categoryID *uint
		var totalAmount float64
		if err := rows.Scan(&transactionType, &categoryID, &totalAmount); err != nil {
			return err
		}

		fmt.Println("data:::", transactionType, totalAmount)
		if _, ok := utils.ValidTransactionTypes[transactionType]; ok && categoryID != nil {
			categoryName := GetCategoryName(c, categoryID)
			if categoryName != "" {
				totalAmounts[transactionType][categoryName] = totalAmount
			}
		}
	}
	fmt.Println("tot::", totalAmounts)
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

func GetBudgetAmount(c *gin.Context, userID uint) (map[string]float64, error) {
	budgetAmounts := make(map[string]float64)

	rows, err := initializers.DB.Model(&models.Budgets{}).
		Select("category_id, SUM(amount) as total_amount").
		Where("user_id = ?", userID).
		Group("category_id").
		Rows()
	if err != nil {
		return budgetAmounts, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryID uint
		var totalAmount float64
		if err := rows.Scan(&categoryID, &totalAmount); err != nil {
			return budgetAmounts, err
		}
		categoryName := GetCategoryName(c, &categoryID)
		if categoryName != "" {
			if _, ok := budgetAmounts[categoryName]; !ok {
				budgetAmounts[categoryName] = totalAmount
			}
		}
	}

	return budgetAmounts, nil
}

func ConvertCurrency(amount float64, currency string, defaultCurrency *string) (float64, error) {
	apiKey := os.Getenv("API_KEY")
	convertURL := fmt.Sprintf("https://anyapi.io/api/v1/exchange/convert?apiKey=%s&base=%s&to=%s&amount=%d", apiKey, currency, *defaultCurrency, uint(amount))
	fmt.Println("args:base:", currency, "to:", *defaultCurrency, "amount:", amount)
	fmt.Println("url:", convertURL)

	var result struct {
		Base            string  `json:"base"`
		To              string  `json:"to"`
		Amount          uint    `json:"amount"`
		ConvertedAmount float64 `json:"converted"`
		Rate            float64 `json:"rate"`
		LatestUpdate    int64   `json:"latestUpdate"`
	}

	maxRetries := 3
	retryDelay := time.Second

	for retries := 0; retries < maxRetries; retries++ {
		resp, err := http.Get(convertURL)
		if err != nil {
			return 0, fmt.Errorf("failed to make API request: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Println("Rate limit exceeded. Retrying...")
			resp.Body.Close() // Close the response body before retrying
			time.Sleep(retryDelay)
			retryDelay *= 2
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("API request failed with status: %s", resp.Status)
		}

		fmt.Println("resp::::", resp)
		fmt.Println("body::", resp.Body)

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return 0, fmt.Errorf("failed to decode API response: %w", err)
		}

		if result.ConvertedAmount == 0 {
			return 0, fmt.Errorf("API returned zero converted amount, result: %+v", result)
		}

		fmt.Println("convAmount:", result.ConvertedAmount, result)
		fmt.Println("convAmount:", uint(result.ConvertedAmount))

		return result.ConvertedAmount, nil
	}

	return 0, fmt.Errorf("failed to convert currency after %d attempts", maxRetries)
}
