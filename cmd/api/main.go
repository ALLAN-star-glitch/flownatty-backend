package main

import (
	"log"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/gin-gonic/gin"
	// Database import commented out for now
	// "github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting Flownatty MVP1 API in %s mode", cfg.Environment)

	// Database connection - commented out for now
	// if err := database.Connect(cfg); err != nil {
	//     log.Fatalf("Failed to connect to database: %v", err)
	// }

	// Initialize router
	router := gin.Default()

	// Health check endpoint with IndentedJSON
	router.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{
			"status":    "ok",
			"service":   "flownatty-api",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Welcome endpoint with IndentedJSON
	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{
			"message": "Welcome to Flownatty API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}