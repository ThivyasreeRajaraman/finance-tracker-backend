package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Thivyasree-Rajaraman/finance-tracker/routes"
)

func ApplyCorsConfig(router *gin.Engine) {
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))
}

func main() {
	r := gin.Default()

	// Apply CORS configuration
	ApplyCorsConfig(r)

	// Setup routes
	routes.SetupRoutes(r)

	// Start the server
	if err := r.Run(); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
