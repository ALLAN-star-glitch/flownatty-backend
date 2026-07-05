package main

import (
	"context"
	"log"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/auth/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/auth"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/auth/middleware"
	"github.com/gin-gonic/gin"

	// Swagger dependencies
	_ "github.com/ALLAN-star-glitch/flownatty-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger annotations
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
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// ================================================
	// 1. Initialize Permissions Module
	// ================================================
	permModule, err := permissions.NewModule(db, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize permissions module: %v", err)
	}

	// Seed permissions (if not already seeded)
	ctx := context.Background()
	if err := permModule.Init(ctx); err != nil {
		log.Fatalf("Failed to seed permissions: %v", err)
	}

	enforcer := permModule.GetEnforcer()
	permService := permModule.GetService()

	// ================================================
	// 2. Initialize Auth Module with Permissions
	// ================================================
	authModule := auth.NewAuthModule(cfg, permService)

	// ================================================
	// 3. Initialize Router
	// ================================================
	router := gin.Default()

	// API v1 group
	v1 := router.Group("/api/v1")

	// ================================================
	// 4. Public Routes (No Auth Required)
	// ================================================
	authModule.SetupRoutes(v1)

	// ================================================
	// 5. Protected Routes (Auth + Authorization)
	// ================================================
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		// ================================================
		// 5a. Consumer Routes
		// ================================================
		consumer := protected.Group("/consumer")
		consumer.Use(permissions.AuthorizationMiddleware(enforcer))
		{
			consumer.GET("/products", func(c *gin.Context) {
				userID := c.GetString(permissions.ContextKeyUserID)
				c.JSON(200, gin.H{
					"message": "Consumer products list",
					"user_id": userID,
				})
			})

			consumer.GET("/orders", func(c *gin.Context) {
				userID := c.GetString(permissions.ContextKeyUserID)
				c.JSON(200, gin.H{
					"message": "Consumer orders",
					"user_id": userID,
				})
			})
		}

		// ================================================
		// 5b. Business Routes
		// ================================================
		business := protected.Group("/business")
		business.Use(permissions.AuthorizationMiddleware(enforcer))
		{
			// Routes that require business owner role
			owner := business.Group("/:businessId")
			owner.Use(permissions.RequireBusinessOwner(enforcer))
			{
				owner.GET("/products", func(c *gin.Context) {
					businessID := c.GetString(permissions.ContextKeyBusinessID)
					c.JSON(200, gin.H{
						"message":     "Business products",
						"business_id": businessID,
					})
				})

				owner.POST("/products", func(c *gin.Context) {
					businessID := c.GetString(permissions.ContextKeyBusinessID)
					c.JSON(201, gin.H{
						"message":     "Product created",
						"business_id": businessID,
					})
				})

				owner.GET("/orders", func(c *gin.Context) {
					businessID := c.GetString(permissions.ContextKeyBusinessID)
					c.JSON(200, gin.H{
						"message":     "Business orders",
						"business_id": businessID,
					})
				})

				owner.GET("/dashboard", func(c *gin.Context) {
					businessID := c.GetString(permissions.ContextKeyBusinessID)
					c.JSON(200, gin.H{
						"message":     "Business dashboard",
						"business_id": businessID,
					})
				})
			}

			// Routes that require any business access (owner or staff)
			staff := business.Group("/:businessId")
			staff.Use(permissions.RequireBusinessAccess(enforcer))
			{
				staff.GET("/chat", func(c *gin.Context) {
					businessID := c.GetString(permissions.ContextKeyBusinessID)
					c.JSON(200, gin.H{
						"message":     "Business chat",
						"business_id": businessID,
					})
				})
			}
		}

		// ================================================
		// 5c. Admin Routes
		// ================================================
		admin := protected.Group("/admin")
		admin.Use(permissions.AuthorizationMiddleware(enforcer))
		admin.Use(permissions.RequirePlatformRole(enforcer, permissions.RoleAdmin, permissions.RoleSuperAdmin))
		{
			admin.GET("/users", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "User list (admin only)",
				})
			})

			admin.GET("/businesses", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Business list (admin only)",
				})
			})
		}
	}

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

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}