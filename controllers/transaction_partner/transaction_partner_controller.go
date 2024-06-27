package transactionpartnercontrollers

import (
	"net/http"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	transactionpartnerservice "github.com/Thivyasree-Rajaraman/finance-tracker/services/transaction_partner"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

type PartnerController struct{}

type PartnerControllerInterface interface {
	FetchOrCreate(c *gin.Context)
	Fetch(c *gin.Context)
	FetchPartners(c *gin.Context)
	Remind(c *gin.Context)
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

	partnerModel := new(models.TransactionPartner)
	userID, err := utils.GetUserID(c)
	transactionType := c.Param("transactionType")
	if err != nil {
		return
	}
	conditions := map[string]interface{}{
		"user_id": userID,
	}
	var greaterThanConditions, lesserThanConditions, unequalThanConditions map[string]interface{}
	switch transactionType {
	case "lend":
		lesserThanConditions = map[string]interface{}{
			"closing_balance": 0,
		}
	case "borrow":
		greaterThanConditions = map[string]interface{}{
			"closing_balance": 0,
		}
	case "all":
		unequalThanConditions = map[string]interface{}{
			"closing_balance": 0,
		}
	default:
		return
	}
	if data := utils.List(c, partnerModel, conditions, unequalThanConditions, greaterThanConditions, lesserThanConditions, "due_date ASC"); data != nil {
		return
	}
}

func (controller *PartnerController) FetchPartners(c *gin.Context) {

	partners, err := transactionpartnerservice.Fetch(c)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch partners", err)
	}
	partnerNames := make([]string, len(partners))
	for i, partner := range partners {
		partnerNames[i] = partner.PartnerName
	}
	utils.SendResponse(c, "Transaction partners fetched successfully", "partners", partnerNames)
}

func (controller *PartnerController) Remind(c *gin.Context) {
	transactionpartnerservice.NotifyUpcomingDueDate(c)
}
