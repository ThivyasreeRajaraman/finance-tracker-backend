package routes

import (
	transactioncontrollers "github.com/Thivyasree-Rajaraman/finance-tracker/controllers/transaction"
	"github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(router *gin.RouterGroup) {
	controller := transactioncontrollers.GetTransactionControllerInstance()
	router.POST("/user/transaction", controller.Create)
	router.PUT("/user/transaction/:transactionId", controller.Update)
	router.DELETE("/user/transaction/:transactionId", controller.Delete)
}
