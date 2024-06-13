package transactionservices

import (
	"fmt"
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	categoryhelpers "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/category"
	dbhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/common"
	transactionhelpers "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/transaction"
	transactionpartnerhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/transaction_partner"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UnmarshalAndValidate(c *gin.Context, transactionData *helpers.TransactionData) error {
	if err := utils.UnmarshalData(c, transactionData); err != nil {
		return err
	}
	if transactionData.TransactionType == "" {
		return utils.CreateError("transaction_type cannot be empty")
	}
	if transactionData.Amount <= 0 {
		return utils.CreateError("amount must be greater than zero")
	}
	if err := utils.IsValidTransactionType(transactionData.TransactionType); err != nil {
		return utils.CreateError("transaction type not found")
	}
	if err := utils.IsValidCurrency(transactionData.Currency); err != nil {
		return utils.CreateError("Invalid Currency Code")
	}
	return nil
}

func UnmarshalAndValidateSingleEntity(c *gin.Context, transactionData *helpers.TransactionUpdate) error {
	if err := utils.UnmarshalData(c, transactionData); err != nil {
		return err
	}
	if transactionData.CategoryName != nil && *transactionData.CategoryName == "" {
		return utils.CreateError("category_name cannot be empty")
	}
	if transactionData.Amount != nil && *transactionData.Amount <= 0 {
		return utils.CreateError("amount must be greater than zero")
	}
	if transactionData.TransactionType != nil && *transactionData.TransactionType == "" {
		return utils.CreateError("transaction type cannot be empty")
	}
	if transactionData.TransactionType != nil && *transactionData.TransactionType != "income" && *transactionData.TransactionType != "expense" {
		return utils.CreateError("transaction type can be income or expense only")
	}
	return nil
}

func getDefaultCurrencyForUser(c *gin.Context) (string, error) {
	userInterface, _ := c.Get("currentUser")
	user, ok := userInterface.(models.User)
	if !ok {
		err := utils.CreateError("invalid user data")
		utils.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return "", err
	}
	return *user.DefaultCurrency, nil
}

func GetConvertedCurrency(c *gin.Context, amount uint, currency string) (float64, error) {
	defaultCurrency, err := getDefaultCurrencyForUser(c)
	if err != nil {
		return utils.Zero, utils.CreateError("Failed to retrieve default currency")
	}

	convertedAmount, err := transactionhelpers.ConvertCurrency(float64(amount), currency, &defaultCurrency)
	if err != nil {
		return utils.Zero, utils.CreateError(fmt.Sprintf("Currency conversion failed: %v", err))
	}
	return convertedAmount, nil
}

func CreateTransaction(c *gin.Context, transactionData helpers.TransactionData, userID uint) error {
	transaction := models.Transaction{
		UserID:          userID,
		TransactionType: transactionData.TransactionType,
		Amount:          transactionData.Amount,
		Currency:        transactionData.Currency,
	}

	convertedAmount, err := GetConvertedCurrency(c, transactionData.Amount, transactionData.Currency)
	if err != nil {
		return utils.CreateError(fmt.Sprintf("Currency conversion failed: %v", err))
	}

	transaction.ConvertedAmount = uint(convertedAmount)

	switch transactionData.TransactionType {
	case "income", "expense":
		if err := transactionhelpers.HandleIncomeExpenseTransaction(userID, &transaction, transactionData); err != nil {
			return err
		}
	case "lend", "borrow":
		if err := transactionhelpers.HandleLendBorrowTransaction(userID, &transaction, transactionData); err != nil {
			return err
		}
	default:
		return utils.CreateError("Invalid transaction type")
	}

	// Save transaction to DB
	if err := dbhelper.GenericCreate(&transaction); err != nil {
		return err
	}
	if err := preloadTransactionAssociations(c, &transaction); err != nil {
		return err
	}
	transactionResponse, err := utils.CreateTransactionResponse(transaction)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to construct transaction response", err)
	}

	// Verify if the threshold has been reached
	if transactionData.TransactionType == "expense" {
		alert := transactionhelpers.CheckThreshold(c, transactionResponse, userID)
		if alert != "" {
			c.JSON(http.StatusOK, gin.H{
				"success":     true,
				"message":     "Transaction created successfully",
				"transaction": transactionResponse,
				"alert":       alert,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success":     true,
				"message":     "Transaction created successfully",
				"transaction": transactionResponse,
			})
		}
	} else {
		utils.SendResponse(c, "Transaction created successfully", "transaction", transactionResponse)
	}
	return nil
}

func FetchTransactionById(c *gin.Context, transaction *models.Transaction, transactionId uint) error {
	if err := transactionhelpers.FetchByID(c, transaction, transactionId); err != nil {
		return err
	}
	return nil
}

func preloadTransactionAssociations(c *gin.Context, transaction *models.Transaction) error {
	if err := initializers.DB.Preload("Category").Preload("User").Preload("TransactionPartner").First(transaction, transaction.ID).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to preload user and category association", err)
		return err
	}
	return nil
}

func UpdateExistingTransaction(c *gin.Context, existingTransaction *models.Transaction, transactionData helpers.TransactionUpdate, categoryID uint) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}
	if categoryID != 0 && existingTransaction.CategoryID != categoryID {
		existingTransaction.CategoryID = categoryID
	}
	if transactionData.Amount != nil {
		existingTransaction.Amount = *transactionData.Amount
		convertedAmount, err := GetConvertedCurrency(c, *transactionData.Amount, existingTransaction.Currency)
		if err != nil {
			return utils.CreateError(fmt.Sprintf("Currency conversion failed: %v", err))
		}
		existingTransaction.ConvertedAmount = uint(convertedAmount)
	}
	if transactionData.TransactionType != nil {
		existingTransaction.TransactionType = *transactionData.TransactionType
	}

	if err := dbhelper.GenericUpdate(existingTransaction); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update transaction", err)
		return err
	}
	if err := preloadTransactionAssociations(c, existingTransaction); err != nil {
		return err
	}
	transactionResponse, err := utils.CreateTransactionResponse(*existingTransaction)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to construct transaction response", err)
	}
	if transactionData.Amount != nil {
		if existingTransaction.TransactionType == "expense" {
			// transactionhelpers.CheckThreshold(c, transactionResponse, userID)
			alert := transactionhelpers.CheckThreshold(c, transactionResponse, userID)
			if alert != "" {
				c.JSON(http.StatusOK, gin.H{
					"success":     true,
					"message":     "Transaction updated successfully",
					"transaction": transactionResponse,
					"alert":       alert,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"success":     true,
					"message":     "Transaction updated successfully",
					"transaction": transactionResponse,
				})
			}
		} else {
			utils.SendResponse(c, "Transaction updated successfully", "transaction", transactionResponse)
		}
	} else {
		utils.SendResponse(c, "Transaction updated successfully", "transaction", transactionResponse)
	}
	return nil
}

func Update(c *gin.Context, existingTransaction *models.Transaction, transactionData helpers.TransactionUpdate) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}
	var category *models.Categories
	if existingTransaction.TransactionType != "income" && existingTransaction.TransactionType != "expense" {
		utils.HandleError(c, http.StatusInternalServerError, "Only income/expense transactions can be updated", err)
		return err
	}
	if transactionData.CategoryName != nil {
		category, err = categoryhelpers.GetOrCreateCategory(userID, transactionData.CategoryName, existingTransaction.TransactionType)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, "Failed to get or create category", err)
			return err
		}
	} else {
		category = &models.Categories{Model: gorm.Model{ID: 0}}
	}
	// Update the existing transaction
	if err := UpdateExistingTransaction(c, existingTransaction, transactionData, category.ID); err != nil {
		return err
	}
	return nil
}

func GetTransactionFromPathParam(c *gin.Context) (*models.Transaction, error) {
	transactionId, err := utils.ParseUintParam(c, "transactionId")
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Invalid transaction id", err)
		return nil, err
	}
	var existingTransaction models.Transaction
	if err := FetchTransactionById(c, &existingTransaction, transactionId); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch transaction", err)
		return nil, err
	}

	return &existingTransaction, nil
}

func Delete(c *gin.Context, transaction *models.Transaction) error {

	if transaction.TransactionType == "lend" || transaction.TransactionType == "borrow" {
		var targetType string
		if transaction.TransactionType == "lend" {
			targetType = "borrow"
		} else {
			targetType = "lend"
		}
		if err := transactionpartnerhelper.UpdateTransactionPartnerAmount(*transaction.TransactionPartnerID, targetType, transaction.Amount, nil); err != nil {
			return err
		}
	}

	if err := dbhelper.GenericDelete(transaction); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to delete transaction", err)
		return err
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Transaction deleted successfully"})
	return nil
}

func CalculateTotal(c *gin.Context) error {
	if err := transactionhelpers.CalculateTotalAmounts(c); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve total amount of transaction", err)
		return err
	}
	return nil
}

func CalculateCategoryWiseTotal(c *gin.Context) error {
	if err := transactionhelpers.CalculateCategoryWiseAmounts(c); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to retrieve total amount of transaction", err)
		return err
	}
	return nil
}
