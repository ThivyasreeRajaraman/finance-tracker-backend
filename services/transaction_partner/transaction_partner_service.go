package transactionpartnerservice

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	transactionpartnerhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/transaction_partner"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
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

func Fetch(c *gin.Context) ([]models.TransactionPartner, error) {
	var partners []models.TransactionPartner
	if err := transactionpartnerhelper.FetchAll(c, &partners); err != nil {
		return nil, err
	}
	return partners, nil
}

func NotifyUpcomingDueDate(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return
	}
	var upcomingLendOrBorrowTransactions []models.TransactionPartner
	if err := initializers.DB.Model(&models.TransactionPartner{}).
		Where("user_id = ? AND due_date BETWEEN ? AND ?", userID, time.Now(), time.Now().AddDate(0, 0, 5)).
		Preload("User").Find(&upcomingLendOrBorrowTransactions).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch upcoming lend or borrow transactions", err)
		return
	}
	var reminders []string
	for _, transaction := range upcomingLendOrBorrowTransactions {
		formattedDate, err := time.Parse("2006-01-02", transaction.DueDate)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, "Failed to parse next expense date", err)
			continue
		}
		daysUntilExpense := int(time.Until(formattedDate).Hours() / 24)
		if daysUntilExpense <= 5 {
			reminders = append(reminders, sendLendOrBorrowReminder(transaction, daysUntilExpense))
		}
	}
	reminderMessage := strings.Join(reminders, "\n")
	c.JSON(http.StatusOK, gin.H{"Reminder": reminderMessage})
}

func sendLendOrBorrowReminder(transaction models.TransactionPartner, daysUntilExpense int) string {
	var message, transaction_type, adj string
	fmt.Println("\n\ndata::", daysUntilExpense)
	if transaction.ClosingBalance > 0 {
		transaction_type = "Borrowed"
		adj = "from"
	} else if transaction.ClosingBalance < 0 {
		transaction_type = "Lent"
		adj = "to"
	}

	if daysUntilExpense == 0 {
		message = fmt.Sprintf("The amount of %d %s you %s %s %s is due today.",
			uint(math.Abs(float64(transaction.ClosingBalance))), *transaction.User.DefaultCurrency, transaction_type, adj, transaction.PartnerName)
	} else {
		fmt.Printf("data::%+v", transaction)
		message = fmt.Sprintf("The amount of %d %s you %s %s %s is due in %d day(s).",
			uint(math.Abs(float64(transaction.ClosingBalance))), *transaction.User.DefaultCurrency, transaction_type, adj, transaction.PartnerName, daysUntilExpense)
	}
	return message
}
