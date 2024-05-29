package transactionpartnerservice

import (
	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	transactionpartnerhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/transaction_partner"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func UnmarshalAndValidate(c *gin.Context, transactionPartnerData *helpers.TransactionPartnerData) error {
	if err := utils.UnmarshalData(c, transactionPartnerData); err != nil {
		return err
	}
	if transactionPartnerData.PartnerName == "" {
		return utils.CreateError("Transaction Partner name cannot be empty")
	}
	return nil
}

func GetOrCreatePartner(userID uint, partnerName *string) ([]helpers.TransactionPartnerResponse, error) {
	var partner *models.TransactionPartner
	var err error
	partner, err = transactionpartnerhelper.Fetch(userID, partnerName)
	if err != nil {
		partner, err = transactionpartnerhelper.Create(userID, partnerName)
		if err != nil {
			return nil, err
		}
	}
	partners := []models.TransactionPartner{*partner}
	partnerResponse, err := utils.CreatePartnerResponse(partners)
	if err != nil {
		return nil, err
	}
	return partnerResponse, nil
}

func Fetch(c *gin.Context, partners *[]models.TransactionPartner) error {
	if err := transactionpartnerhelper.FetchAll(c, partners); err != nil {
		return err
	}
	return nil
}
