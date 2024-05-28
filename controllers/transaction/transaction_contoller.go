package transactioncontrollers

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	transactionservices "github.com/Thivyasree-Rajaraman/finance-tracker/services/transaction"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
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

func Update(c *gin.Context) {
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
