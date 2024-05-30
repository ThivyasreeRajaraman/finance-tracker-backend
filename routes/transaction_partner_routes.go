package routes

import (
	transactionpartnercontrollers "github.com/Thivyasree-Rajaraman/finance-tracker/controllers/transaction_partner"
	"github.com/gin-gonic/gin"
)

func SetupTransactionPartnerRoutes(router *gin.RouterGroup) {
	controller := transactionpartnercontrollers.GetPartnerControllerInstance()
	router.POST("/user/transactionpartner", controller.FetchOrCreate)
	router.GET("/user/transactionpartner", controller.Fetch)
	router.GET("/user/lendOrBorrowDuedate/reminder", controller.Remind)
}
