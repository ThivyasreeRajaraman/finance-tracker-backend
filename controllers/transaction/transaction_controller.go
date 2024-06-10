package transactioncontrollers

import (
	"net/http"
	"sort"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	transactionservices "github.com/Thivyasree-Rajaraman/finance-tracker/services/transaction"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

type TransactionController struct{}

type TransactionControllerInterface interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	FetchTransactionTypes(c *gin.Context)
	FetchTotal(c *gin.Context)
	FetchSingleTransaction(c *gin.Context)
	Fetch(c *gin.Context)
}

func GetTransactionControllerInstance() TransactionControllerInterface {
	return new(TransactionController)
}

func (controller *TransactionController) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return
	}
	var transactionData helpers.TransactionData

	if err := transactionservices.UnmarshalAndValidate(c, &transactionData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
		return
	}

	if err := transactionservices.CreateTransaction(c, transactionData, userID); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to create transaction", err)
		return
	}
}

func (controller *TransactionController) Update(c *gin.Context) {
	existingTransaction, err := transactionservices.GetTransactionFromPathParam(c)
	if err != nil {
		return
	}
	var transactionData helpers.TransactionUpdate
	if err := transactionservices.UnmarshalAndValidateSingleEntity(c, &transactionData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
		return
	}

	if err := transactionservices.Update(c, existingTransaction, transactionData); err != nil {
		return
	}
}

func (controller *TransactionController) Delete(c *gin.Context) {
	existingTransaction, err := transactionservices.GetTransactionFromPathParam(c)
	if err != nil {
		return
	}
	// Soft delete the transaction
	if err := transactionservices.Delete(c, existingTransaction); err != nil {
		return
	}
}

func (controller *TransactionController) FetchTransactionTypes(c *gin.Context) {
	transactionTypes := make([]string, 0, len(utils.ValidTransactionTypes))
	for transactionType := range utils.ValidTransactionTypes {
		transactionTypes = append(transactionTypes, transactionType)
	}
	sort.Strings(transactionTypes)
	c.JSON(http.StatusOK, gin.H{"transaction_types": transactionTypes})
}

func (controller *TransactionController) FetchTotal(c *gin.Context) {
	transactionservices.CalculateTotal(c)
}

func (controller *TransactionController) FetchSingleTransaction(c *gin.Context) {
	transactionId, err := utils.ParseUintParam(c, "transactionId")
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Invalid transaction id", err)
		return
	}
	var transaction models.Transaction
	transactionservices.FetchTransactionById(c, &transaction, transactionId)
	utils.SendResponse(c, "Transaction fetched successfully", "transaction", transaction)
}

func (controller *TransactionController) Fetch(c *gin.Context) {
	transactionModel := new(models.Transaction)
	userID, err := utils.GetUserID(c)
	if err != nil {
		return
	}
	transactionType := c.Param("transactionType")

	conditions := map[string]interface{}{
		"user_id":          userID,
		"transaction_type": transactionType,
	}
	if data := utils.List(c, transactionModel, conditions, nil, "created_at DESC"); data != nil {
		return
	}
}
