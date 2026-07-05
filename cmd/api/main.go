package main

import (
    "log"
    "time"

    "github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/auth"
    "github.com/gin-gonic/gin"

    //  Swagger dependencies - CORRECT
    _ "github.com/ALLAN-star-glitch/flownatty-backend/docs"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger annotations - CORRECT
// @title Flownatty Backend API
// @version 1.0
// @description Backend API engine for Flownatty platform.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1

func main() {
    // Load configuration
    cfg := config.Load()
    log.Printf("Starting Flownatty MVP1 API in %s mode", cfg.Environment)

    // Connect to database
    if err := database.Connect(cfg); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Initialize router
    router := gin.Default()

    // API v1 group
    v1 := router.Group("/api/v1")

    // Register auth module
    authModule := auth.NewAuthModule(cfg)
    authModule.RegisterRoutes(v1)

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.IndentedJSON(200, gin.H{
            "status":    "ok",
            "service":   "flownatty-api",
            "version":   "1.0.0",
            "timestamp": time.Now().UTC().Format(time.RFC3339),
        })
    })

    // Welcome endpoint
    router.GET("/", func(c *gin.Context) {
        c.IndentedJSON(200, gin.H{
            "message": "Welcome to Flownatty API",
            "version": "1.0.0",
            "status":  "running",
        })
    })

    // Swagger route - CORRECT
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Start server
    log.Printf("Server starting on port %s", cfg.Server.Port)
    if err := router.Run(":" + cfg.Server.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}