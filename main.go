package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Thivyasree-Rajaraman/finance-tracker/routes"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(config))

	// Setup routes
	routes.SetupRoutes(r)
	if err := r.Run(); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
