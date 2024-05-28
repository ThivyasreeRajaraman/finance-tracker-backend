package routes

import (
	usercontrollers "github.com/Thivyasree-Rajaraman/finance-tracker/controllers/user"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup) {
	router.PUT("/user", usercontrollers.Update)
}
