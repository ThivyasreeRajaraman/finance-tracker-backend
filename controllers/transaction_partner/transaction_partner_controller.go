package transactionpartnercontrollers

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	transactionpartnerservice "github.com/Thivyasree-Rajaraman/finance-tracker/services/transaction_partner"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

type PartnerController struct{}

type PartnerControllerInterface interface {
	FetchOrCreate(c *gin.Context)
	Fetch(c *gin.Context)
}

func GetPartnerControllerInstance() PartnerControllerInterface {
	return new(PartnerController)
}

func (controller *PartnerController) FetchOrCreate(c *gin.Context) {
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

func (controller *PartnerController) Fetch(c *gin.Context) {

	partners, err := transactionpartnerservice.Fetch(c)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch partners", err)
	}
	partnerResponses, err := utils.CreatePartnerResponse(partners)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to construct partner response", err)
	}

	utils.SendResponse(c, "Transaction partners fetched successfully", "Transaction Partners", partnerResponses)
}
