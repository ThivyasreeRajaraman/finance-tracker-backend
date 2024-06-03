package transactionhelpers

import (
	"net/http"

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

	transaction.CategoryID = &category.ID
	return nil
}

func HandleLendBorrowTransaction(userID uint, transaction *models.Transaction, transactionData helpers.TransactionData) error {
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

	return nil
}
