package transactionservices

import (
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

func CreateTransaction(c *gin.Context, transactionData helpers.TransactionData, userID uint) error {
	transaction := models.Transaction{
		UserID:          userID,
		TransactionType: transactionData.TransactionType,
		Amount:          transactionData.Amount,
	}

	switch transactionData.TransactionType {
	case "income", "expense":
		if transactionData.CategoryName == nil || *transactionData.CategoryName == "" {
			return utils.CreateError("Category name is required for income/expense transactions")
		}
		category, err := categoryhelpers.GetOrCreateCategory(userID, transactionData.CategoryName, transactionData.TransactionType)
		if err != nil {
			return utils.CreateError("Failed to get or create category")
		}
		transaction.CategoryID = &category.ID

	case "lend", "borrow":
		if err := utils.ValidateLendBorrowData(transactionData); err != nil {
			return err
		}
		partner, err := transactionpartnerhelper.Fetch(userID, transactionData.TransactionPartner)
		if err != nil {
			return utils.CreateError("Failed to get transaction partner")
		}
		transaction.TransactionPartnerID = &partner.ID
		if err := transactionpartnerhelper.UpdateTransactionPartnerAmount(*transaction.TransactionPartnerID, transactionData.TransactionType, transactionData.Amount, transactionData.PaymentDueDate); err != nil {
			return err
		}

	default:
		return utils.CreateError("Invalid transaction type")
	}
	// Save transaction to DB
	if err := dbhelper.GenericCreate(&transaction); err != nil {
		return err
	}
	if err := Fetch(c, &transaction); err != nil {
		return utils.CreateError("Failed to create transaction")
	}
	return nil
}

func Fetch(c *gin.Context, transaction *models.Transaction) error {
	if err := transactionhelpers.Fetch(c, transaction); err != nil {
		return err
	}
	transactionResponses, err := utils.CreateTransactionResponse(*transaction)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to construct transaction response", err)
		return err
	}

	utils.SendResponse(c, "Transaction created successfully", "transaction", transactionResponses)
	return nil
}

func FetchTransactionById(c *gin.Context, transaction *models.Transaction, transactionId uint) error {
	if err := transactionhelpers.FetchByID(c, transaction, transactionId); err != nil {
		return err
	}
	return nil
}

func preloadTransactionAssociations(c *gin.Context, transaction *models.Transaction) error {
	if err := initializers.DB.Preload("Category").Preload("User").First(transaction, transaction.ID).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to preload user and category association", err)
		return err
	}
	return nil
}

func UpdateExistingTransaction(c *gin.Context, existingTransaction *models.Transaction, transactionData helpers.TransactionUpdate, categoryID uint) error {
	if categoryID != 0 && existingTransaction.CategoryID != &categoryID {
		existingTransaction.CategoryID = &categoryID
	}
	if transactionData.Amount != nil {
		existingTransaction.Amount = *transactionData.Amount
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
	utils.SendResponse(c, "Transaction updated successfully", "transaction", existingTransaction)
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
