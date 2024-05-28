package transactionpartnercontrollers

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	transactionpartnerservice "github.com/Thivyasree-Rajaraman/finance-tracker/services/transaction_partner"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func FetchOrCreate(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return
	}

	var transactionPartnerData helpers.TransactionPartnerData

	if err := transactionpartnerservice.UnmarshalAndValidate(c, &transactionPartnerData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to unmarshal request body", err)
		return
	}

	partner, err := transactionpartnerservice.GetOrCreatePartner(userID, &transactionPartnerData.PartnerName)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to create transaction", err)
		return
	}
	utils.SendResponse(c, "Transaction partner fetched/created successfully", "transaction_partner", partner)
}

func Fetch(c *gin.Context) {

	var partners []models.TransactionPartner
	if err := transactionpartnerservice.Fetch(c, &partners); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch partners", err)
	}
	partnerResponses, err := utils.CreatePartnerResponse(partners)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to construct partner response", err)
	}

	utils.SendResponse(c, "Transaction partners fetched successfully", "Transaction Partners", partnerResponses)
}
