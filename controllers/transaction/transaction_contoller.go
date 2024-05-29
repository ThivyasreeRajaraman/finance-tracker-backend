package transactioncontrollers

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	transactionservices "github.com/Thivyasree-Rajaraman/finance-tracker/services/transaction"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

type TransactionController struct{}

type TransactionControllerInterface interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
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
