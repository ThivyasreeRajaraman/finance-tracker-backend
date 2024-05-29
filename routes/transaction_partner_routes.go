package routes

import (
	transactionpartnercontrollers "github.com/Thivyasree-Rajaraman/finance-tracker/controllers/transaction_partner"
	"github.com/gin-gonic/gin"
)

func SetupTransactionPartnerRoutes(router *gin.RouterGroup) {
	router.POST("/user/transactionpartner", transactionpartnercontrollers.FetchOrCreate)
	router.GET("/user/transactionpartner", transactionpartnercontrollers.Fetch)
}
