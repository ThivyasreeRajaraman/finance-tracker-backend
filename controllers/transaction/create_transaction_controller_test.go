package transactioncontrollers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Thivyasree-Rajaraman/finance-tracker/helpers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Transaction", func() {
	controller := GetTransactionControllerInstance()
	var w *httptest.ResponseRecorder
	var c *gin.Context
	var currency = "INR"
	var mockUserData models.User
	var mockPartnerData models.TransactionPartner

	BeforeEach(func() {
		mockUserData = models.User{
			Email:           "testuser@example.com",
			Name:            "TestUser",
			DefaultCurrency: &currency,
		}
		initializers.DB.Create(&mockUserData)
		mockPartnerData = models.TransactionPartner{
			PartnerName:    "TestPartner",
			UserID:         mockUserData.ID,
			ClosingBalance: 0,
			DueDate:        time.Now().Format("2006-01-02"),
		}
		initializers.DB.Create(&mockPartnerData)
		gin.SetMode(gin.TestMode)
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	})

	Context("Valid Income/Expense Transaction Creation", func() {
		It("Should return the right response and transaction data", func() {
			c.Set("currentUser", mockUserData)
			shopping := "Shopping"
			expense := "expense"
			var amount uint = 6000

			transactionParams := helpers.TransactionData{
				TransactionType: expense,
				CategoryName:    &shopping,
				Amount:          amount,
			}

			requestBody, err := json.Marshal(transactionParams)
			Expect(err).To(BeNil())

			req, _ := http.NewRequest(http.MethodPost, "/api/user/transaction", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			controller.Create(c)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			fmt.Println("respose::", response)
			Expect(response["message"]).To(Equal("Transaction created successfully"))
			transaction := response["transaction"].(map[string]interface{})
			Expect(uint(transaction["user_id"].(float64))).To(Equal(mockUserData.ID))
			Expect(transaction["name"]).To(Equal(mockUserData.Name))
			Expect(transaction["category_name"]).To(Equal(*transactionParams.CategoryName))
			Expect(uint(transaction["amount"].(float64))).To(Equal(transactionParams.Amount))
			Expect(transaction["default_currency"]).To(Equal(*mockUserData.DefaultCurrency))

		})
	})

	Context("Valid Lend/Borrow Transaction Creation", func() {
		It("Should return the right response and transaction data", func() {
			c.Set("currentUser", mockUserData)
			partner := "TestPartner"
			lend := "lend"
			dueDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
			var amount uint = 6000

			transactionParams := helpers.TransactionData{
				TransactionType:    lend,
				TransactionPartner: &partner,
				Amount:             amount,
				PaymentDueDate:     &dueDate,
			}

			requestBody, err := json.Marshal(transactionParams)
			Expect(err).To(BeNil())

			req, _ := http.NewRequest(http.MethodPost, "/api/user/transaction", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			controller.Create(c)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			fmt.Println("respose::", response)
			Expect(response["message"]).To(Equal("Transaction created successfully"))
			transaction := response["transaction"].(map[string]interface{})
			Expect(uint(transaction["user_id"].(float64))).To(Equal(mockUserData.ID))
			Expect(transaction["name"]).To(Equal(mockUserData.Name))
			Expect(transaction["transaction_partner"]).To(Equal(*transactionParams.TransactionPartner))
			Expect(uint(transaction["amount"].(float64))).To(Equal(transactionParams.Amount))
			Expect(transaction["default_currency"]).To(Equal(*mockUserData.DefaultCurrency))
		})
	})

	Context("Invalid request data", func() {
		It("Should return error when request body is invalid", func() {
			c.Set("currentUser", mockUserData)
			shopping := "Shopping"
			expense := "expense"
			var amount uint = 0

			transactionParams := helpers.TransactionData{
				TransactionType: expense,
				CategoryName:    &shopping,
				Amount:          amount,
			}

			requestBody, err := json.Marshal(transactionParams)
			Expect(err).To(BeNil())

			req, _ := http.NewRequest(http.MethodPost, "/api/user/transaction", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			controller.Create(c)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			fmt.Println("respose::", response)
			Expect(response["error"]).To(Equal("Failed to unmarshal request body"))
			Expect(response["details"]).To(Equal("amount must be greater than zero"))

		})
	})

})
