package authentication

import (
	"context"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authhandler"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authrepo"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/queue"
	"github.com/gin-gonic/gin"
)

type AuthModule struct {
	authHandler  *authhandler.AuthHandler
	tokenService *authservice.TokenService
	authService  *authservice.AuthService
}

// NewAuthModule creates a new auth module
func NewAuthModule(
	cfg *config.Config,
	permService *permissions.Service,
	businessService *bizservice.BusinessService,
) *AuthModule {
	// Initialize queue client
	queueClient := queue.NewClient(cfg.Redis.URL)

	// Initialize repository (internal only)
	repo := authrepo.NewAuthRepository(database.GetDB())

	// Initialize token service
	tokenService := authservice.NewTokenService(repo, cfg)

	// Initialize auth service with all dependencies
	authService := authservice.NewAuthService(
		repo,
		cfg,
		queueClient,
		permService,
		businessService,
		tokenService,
	)

	// Initialize handler with all dependencies
	// NewAuthHandler signature: (service, repo, cfg, permService, tokenService, businessService)
	authHandler := authhandler.NewAuthHandler(
		authService,
		repo,
		cfg,
		permService,
		tokenService,
		businessService,
	)

	return &AuthModule{
		authHandler:  authHandler,
		tokenService: tokenService,
		authService:  authService,
	}
}

// SetupRoutes registers all auth routes with the router
func (m *AuthModule) SetupRoutes(r *gin.RouterGroup) {
	RegisterAuthRoutes(r, m.authHandler)
}

// GetAuthHandler returns the auth handler
func (m *AuthModule) GetAuthHandler() *authhandler.AuthHandler {
	return m.authHandler
}

// GetTokenService returns the token service
func (m *AuthModule) GetTokenService() *authservice.TokenService {
	return m.tokenService
}

// GetAuthService returns the auth service
func (m *AuthModule) GetAuthService() *authservice.AuthService {
	return m.authService
}

// Init initializes the module
func (m *AuthModule) Init(ctx context.Context) error {
	log.Println("Auth module initialized")
	return nil
}

// Close closes the module and cleans up resources
func (m *AuthModule) Close() {
	log.Println("Auth module closed")
}