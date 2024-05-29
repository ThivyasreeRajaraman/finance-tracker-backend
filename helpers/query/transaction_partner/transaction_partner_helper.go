package transactionpartnerhelper

import (
	"strings"
	"time"

	dbhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/common"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
)

func UpdateTransactionPartnerAmount(id uint, transactionType string, amount uint, duedate *time.Time) error {
	var partner models.TransactionPartner
	if err := initializers.DB.Where("id = ?", id).First(&partner).Error; err != nil {
		return utils.CreateError("Failed to find partner")
	}
	// Update the ClosingBalance based on the transaction type
	if duedate != nil {
		partner.DueDate = *duedate
	}
	switch transactionType {
	case "lend":
		partner.ClosingBalance -= int(amount) // negative closing balance for lend
	case "borrow":
		partner.ClosingBalance += int(amount) // positive closing balace for borrow
	default:
		return utils.CreateError("invalid transaction type")
	}
	if err := initializers.DB.Save(&partner).Error; err != nil {
		return utils.CreateError("Failed to update partner")
	}

	return nil
}

func Fetch(userID uint, partnerName *string) (*models.TransactionPartner, error) {
	var partner models.TransactionPartner
	err := initializers.DB.Where("user_id = ? AND LOWER(partner_name) = ?", userID, strings.ToLower(*partnerName)).First(&partner).Error
	if err != nil {
		return nil, err
	}
	return &partner, nil
}

func Create(userID uint, partnerName *string) (*models.TransactionPartner, error) {
	partner := models.TransactionPartner{
		PartnerName:    *partnerName,
		UserID:         userID,
		DueDate:        time.Now(),
		ClosingBalance: 0,
	}
	if err := dbhelper.GenericCreate(&partner); err != nil {
		return nil, err
	}
	return &partner, nil
}

func FetchAll(c *gin.Context, partners *[]models.TransactionPartner) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}
	if err := initializers.DB.Where("user_id = ?", userID).Find(&partners).Error; err != nil {
		return err
	}
	return nil
}
