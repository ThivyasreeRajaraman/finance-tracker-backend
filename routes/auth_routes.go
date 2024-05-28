package routes

import (
	controllers "github.com/Thivyasree-Rajaraman/finance-tracker/controllers/auth_controller"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	// Sign in
	router.GET("/auth/google/callback", controllers.GoogleCallbackHandler)
}
