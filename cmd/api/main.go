package main

import (
	"context"
	"log"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	auth "github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication"
	authMiddleware "github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/middleware"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/repository"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/redis"
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
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting Flownatty MVP1 API in %s mode", cfg.Environment)

	// ================================================
	// Initialize Redis (Global)
	// ================================================
	if err := redis.Init(cfg.Redis.URL); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	log.Println("Redis initialized successfully")

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
	// 2. Initialize Business Module
	// ================================================
	businessModule := business.NewBusinessModule(cfg, permService)
	memberRepo := repository.NewBusinessMemberRepository(db)

	// ================================================
	// 3. Initialize User Management Module
	// ================================================
	userManagementModule := usermanagement.NewUserManagementModule(cfg)

	// ================================================
	// 4. Initialize Auth Module
	// ================================================
	authModule := auth.NewAuthModule(
		cfg,
		permService,
		memberRepo,
	)

	// ================================================
	// 5. Initialize Router
	// ================================================
	router := gin.Default()

	// API v1 group
	v1 := router.Group("/api/v1")

	// ================================================
	// 6. Public Routes (No Authentication Required)
	// ================================================
	authModule.SetupRoutes(v1)

	// ================================================
	// 7. Protected Routes (Authentication + Authorization Required)
	// ================================================
	// Create auth middleware to validate JWT tokens
	authMW := authMiddleware.AuthMiddleware(cfg.JWT.Secret)

	// Business routes: Uses both Authentication (JWT) and Authorization (Casbin)
	// Authorization middleware is applied inside the business module routes
	businessModule.SetupRoutes(v1, authMW, enforcer)

	// ================================================
	// 8. Admin Routes (Authentication + Admin Role Required)
	// ================================================
	adminGroup := v1.Group("/admin")
	adminGroup.Use(authMW)                                                               // Authentication: Validates JWT
	adminGroup.Use(permissions.RequirePlatformRole(enforcer, permissions.RoleAdmin, permissions.RoleSuperAdmin)) // Authorization: Admin role check
	{
		// User management routes (admin only)
		userManagementModule.SetupRoutes(adminGroup)
	}

	// ================================================
	// 9. Health Check & Welcome Endpoints
	// ================================================
	router.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{
			"status":    "ok",
			"service":   "flownatty-api",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{
			"message": "Welcome to Flownatty API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// ================================================
	// 10. Swagger Documentation
	// ================================================
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ================================================
	// 11. Start Server
	// ================================================
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}