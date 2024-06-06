package routes

import (
	usercontrollers "github.com/Thivyasree-Rajaraman/finance-tracker/controllers/user"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup) {
	router.PUT("/user", usercontrollers.Update)
	router.GET("/currencies", usercontrollers.FetchCurrencies)
	router.GET("/user", usercontrollers.Fetch)
	router.GET("/categories", usercontrollers.FetchCategories)
}
