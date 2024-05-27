package routes

import (
	budgetcontrollers "github.com/Thivyasree-Rajaraman/finance-tracker/controllers/budget"
	"github.com/Thivyasree-Rajaraman/finance-tracker/middleware"
	"github.com/gin-gonic/gin"
)

func SetupBudgetRoutes(router *gin.RouterGroup) {
	controller := budgetcontrollers.GetBudgetControllerInstance()
	router.POST("/user/budget", controller.Create)
	router.GET("/user/budget", func(c *gin.Context) {
		controller.Fetch(c, false)
	})
	router.Use(middleware.GetBudgetFromPathParam())
	router.PUT("/user/budget/:budgetId", controller.Update)
	router.DELETE("/user/budget/:budgetId", controller.Delete)
}
