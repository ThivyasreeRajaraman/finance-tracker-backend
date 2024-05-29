package routes

import (
	transactioncontrollers "github.com/Thivyasree-Rajaraman/finance-tracker/controllers/transaction"
	"github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(router *gin.RouterGroup) {
	router.POST("/user/transaction", transactioncontrollers.Create)
	router.PUT("/user/transaction/:transactionId", transactioncontrollers.Update)
}
