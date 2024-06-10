package transactionhelpers

import (
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
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	rows, err := initializers.DB.Model(&models.Transaction{}).
		Select("transaction_type, SUM(amount) as total_amount").
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startOfMonth, endOfMonth).
		Group("transaction_type").
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
		var transactionType string
		var totalAmount float64
		if err := rows.Scan(&transactionType, &totalAmount); err != nil {
			return err
		}
		if _, ok := utils.ValidTransactionTypes[transactionType]; ok {
			totalAmounts[transactionType] = totalAmount
		}
	}
	var totalBudgetAmount float64
	if err := initializers.DB.Model(&models.Budgets{}).
		Select("SUM(amount)").
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startOfMonth, endOfMonth).
		Scan(&totalBudgetAmount).Error; err != nil {
		return err
	}
	totalAmounts["budget"] = totalBudgetAmount

	utils.SendResponse(c, "Total fetched successfully", "transaction_total", totalAmounts)
	return nil
}
