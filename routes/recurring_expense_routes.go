package routes

import (
	recurringexpensecontrollers "github.com/Thivyasree-Rajaraman/finance-tracker/controllers/recurring_expense"
	"github.com/gin-gonic/gin"
)

func SetupRecurringExpenseRoutes(router *gin.RouterGroup) {
	// router.PUT("/user/recurringExpense", recurringexpensecontrollers.Create)
	controller := recurringexpensecontrollers.GetRecurringExpenseControllerInstance()
	router.POST("/user/recurringExpense", controller.Create)
	router.PUT("/user/recurringExpense/:recurringExpenseId", controller.Update)
	router.DELETE("/user/recurringExpense/:recurringExpenseId", controller.Delete)
	router.GET("/user/recurringExpense", controller.Fetch)
	router.GET("/user/recurringExpense/:recurringExpenseId", controller.FetchSingleEntity)
	router.GET("/user/recurringExpense/reminder", controller.Remind)
	router.GET("/recurringExpense/Frequencies", controller.FetchFrequencies)
	router.PUT("/user/recurringExpense/:recurringExpenseId/updateNextExpenseDate", controller.UpdateNextExpenseDate)
}
