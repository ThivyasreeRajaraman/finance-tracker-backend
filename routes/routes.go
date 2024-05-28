package routes

import (
	"github.com/Thivyasree-Rajaraman/finance-tracker/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	SetupAuthRoutes(r)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	SetupUserRoutes(protected)
}
