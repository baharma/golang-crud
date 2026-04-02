package main

import (
	"crud-test/internal/config"
	"crud-test/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	viperConfig := config.NewViper()

	// Ensure the application can reach the configured database before serving traffic.
	if _, err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Basic health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "CRUD API is running",
		})
	})

	// Get port from config or use default
	port := viperConfig.GetString("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
