package usercontrollers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update User Data", func() {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	var currency = "INR"
	var mockUserData = models.User{
		Email:           "thivyasree.ktr@shopup.org",
		Name:            "Sree",
		DefaultCurrency: currency,
	}
	initializers.DB.Create(&mockUserData)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	})

	Context("Valid User Update", func() {
		It("Should return the right response and update user", func() {
			c.Set("currentUser", mockUserData)
			req, _ := http.NewRequest(http.MethodPost, "/api/user/", bytes.NewBuffer([]byte(`{"name":"Thivya","default_currency":"INR"}`)))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			Update(c)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())

			user := response["user"].(map[string]interface{})
			actualID := uint(user["ID"].(float64))

			Expect(actualID).To(Equal(mockUserData.ID))
			Expect(response["message"]).To(Equal("User updated successfully"))
		})
	})

	Context("Invalid request data", func() {
		It("Should return error when request body is invalid", func() {
			c.Set("currentUser", mockUserData)
			req, _ := http.NewRequest(http.MethodPost, "/api/user/", bytes.NewBuffer([]byte(`{"name":1,"default_currency":"INR"}`)))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			Update(c)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response["error"]).To(Equal("Failed to unmarshal request body"))
		})
	})

	Context("Invalid currency code in request body", func() {
		It("Should return error when currency code is valid", func() {
			c.Set("currentUser", mockUserData)
			req, _ := http.NewRequest(http.MethodPost, "/api/user/", bytes.NewBuffer([]byte(`{"name":"Thivya","default_currency":"ABC"}`)))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			Update(c)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response["error"]).To(Equal("Failed to update user"))
		})
	})
})
