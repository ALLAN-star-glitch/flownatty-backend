// internal/app/app.go

package app

import (
	"context"
	"log"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authmidleware"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// Swagger
	_ "github.com/ALLAN-star-glitch/flownatty-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// App represents the application with all modules and dependencies
type App struct {
	Config   *config.Config
	DB       *gorm.DB
	Router   *gin.Engine

	// Modules
	PermissionModule    *permissions.Module
	BusinessModule      *business.BusinessModule
	AuthModule          *authentication.AuthModule
	UserManagementModule *usermanagement.UserManagementModule

	// Services (convenience access)
	PermissionService *permissions.Service
	Enforcer          *permissions.Enforcer
	BusinessService   *bizservice.BusinessService
	AuthService       *authservice.AuthService
	TokenService      *authservice.TokenService
}

// NewApp creates a new application instance with all dependencies initialized
func NewApp() (*App, error) {
	// 1. Load configuration
	cfg := config.Load()
	log.Printf("Starting Flownatty MVP1 API in %s mode", cfg.Environment)

	// 2. Initialize Redis
	if err := redis.Init(cfg.Redis.URL); err != nil {
		return nil, err
	}
	log.Println("Redis initialized successfully")

	// 3. Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}
	log.Println("Database connected successfully")

	// 4. Initialize Permissions Module
	permModule, err := permissions.NewModule(db, cfg)
	if err != nil {
		return nil, err
	}
	enforcer := permModule.GetEnforcer()
	permService := permModule.GetService()

	// 5. Initialize Business Module
	businessModule := business.NewBusinessModule(cfg, permService, enforcer)
	businessService := businessModule.GetBusinessService()

	// 6. Initialize Auth Module
	authModule := authentication.NewAuthModule(cfg, permService, businessService)
	authService := authModule.GetAuthService()
	tokenService := authModule.GetTokenService()

	// 7. Initialize User Management Module
	userManagementModule := usermanagement.NewUserManagementModule(cfg)

	// 8. Setup Router
	router := gin.Default()

	return &App{
		Config:               cfg,
		DB:                   db,
		Router:               router,
		PermissionModule:     permModule,
		BusinessModule:       businessModule,
		AuthModule:           authModule,
		UserManagementModule: userManagementModule,
		PermissionService:    permService,
		Enforcer:             enforcer,
		BusinessService:      businessService,
		AuthService:          authService,
		TokenService:         tokenService,
	}, nil
}

// Init initializes all modules (seeds data, etc.)
func (app *App) Init(ctx context.Context) error {
	// Seed permissions
	if err := app.PermissionModule.Init(ctx); err != nil {
		return err
	}

	// Initialize other modules
	if err := app.BusinessModule.Init(ctx); err != nil {
		return err
	}
	if err := app.AuthModule.Init(ctx); err != nil {
		return err
	}
	if err := app.UserManagementModule.Init(ctx); err != nil {
		return err
	}

	log.Println("All modules initialized successfully")
	return nil
}

// SetupRoutes registers all routes
func (app *App) SetupRoutes() {
	v1 := app.Router.Group("/api/v1")

	// ================================================
	// 1. Public Routes (No Authentication Required)
	// ================================================
	app.AuthModule.SetupRoutes(v1)

	// ================================================
	// 2. Protected Routes (Authentication Required)
	// ================================================
	authMW := authmidleware.AuthMiddleware(app.Config.JWT.Secret)
	app.BusinessModule.SetupRoutes(v1, authMW, app.Enforcer)

	// ================================================
	// 3. Admin Routes (Authentication + Admin Role Required)
	// ================================================
	adminGroup := v1.Group("/admin")
	adminGroup.Use(authMW)
	adminGroup.Use(permissions.RequirePlatformRole(app.Enforcer, permissions.RoleAdmin, permissions.RoleSuperAdmin))
	app.UserManagementModule.SetupRoutes(adminGroup)

	// ================================================
	// 4. Health Check & Welcome Endpoints
	// ================================================
	app.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "flownatty-api",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	app.Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Flownatty API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// ================================================
	// 5. Swagger Documentation
	// ================================================
	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Routes registered successfully")
}

// Run starts the server
func (app *App) Run() error {
	log.Printf("Server starting on port %s", app.Config.Server.Port)
	return app.Router.Run(":" + app.Config.Server.Port)
}

// Close cleans up all resources
func (app *App) Close() {
	log.Println("Closing application...")
	
	// Close modules
	if app.PermissionModule != nil {
		app.PermissionModule.Close()
	}
	if app.BusinessModule != nil {
		app.BusinessModule.Close()
	}
	if app.AuthModule != nil {
		app.AuthModule.Close()
	}
	if app.UserManagementModule != nil {
		app.UserManagementModule.Close()
	}
	
	// Close database connection
	if app.DB != nil {
		if sqlDB, err := app.DB.DB(); err == nil {
			sqlDB.Close()
		}
	}
	
	log.Println("Application closed")
}