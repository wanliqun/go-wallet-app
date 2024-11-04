package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wanliqun/go-wallet-app/config"
	"github.com/wanliqun/go-wallet-app/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	db := config.SetupDatabase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Initialize router
	router := gin.Default()

	// Setup routes
	routes.SetupRouter(router, db)

	// Run server
	log.Printf("Starting server on port %s", config.AppConfig.Server.Port)
	router.Run(":" + config.AppConfig.Server.Port)
}
