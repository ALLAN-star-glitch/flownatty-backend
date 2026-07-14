package authentication

import (
	"context"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authrepo"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/handler"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/service"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/repository"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/queue"
	"github.com/gin-gonic/gin"
)

type AuthModule struct {
	authHandler  *handler.AuthHandler
	tokenService *service.TokenService
	authRepo     *authrepo.AuthRepository
	authService  *service.AuthService
	permService  *permissions.Service
}

// NewAuthModule creates a new auth module with permissions and token services
func NewAuthModule(
	cfg *config.Config,
	permService *permissions.Service,
	memberRepo *repository.BusinessMemberRepository,
) *AuthModule {
	// Initialize queue client
	queueClient := queue.NewClient(cfg.Redis.URL)

	// Initialize repository
	repo := authrepo.NewAuthRepository(database.GetDB())

	// Initialize auth service (no redis client needed - uses global redis)
	authService := service.NewAuthService(repo, cfg, queueClient)

	// Initialize token service
	tokenService := service.NewTokenService(repo, cfg)

	// Initialize handler with all dependencies
	authHandler := handler.NewAuthHandler(
		authService,
		repo,
		memberRepo,
		cfg,
		permService,
		tokenService,
	)

	return &AuthModule{
		authHandler:  authHandler,
		tokenService: tokenService,
		authRepo:     repo,
		authService:  authService,
		permService:  permService,
	}
}

// SetupRoutes registers all auth routes with the router
func (m *AuthModule) SetupRoutes(r *gin.RouterGroup) {
	RegisterAuthRoutes(r, m.authHandler)
}

// GetAuthHandler returns the auth handler
func (m *AuthModule) GetAuthHandler() *handler.AuthHandler {
	return m.authHandler
}

// GetTokenService returns the token service
func (m *AuthModule) GetTokenService() *service.TokenService {
	return m.tokenService
}

// GetAuthRepo returns the auth repository
func (m *AuthModule) GetAuthRepo() *authrepo.AuthRepository {
	return m.authRepo
}

// GetAuthService returns the auth service
func (m *AuthModule) GetAuthService() *service.AuthService {
	return m.authService
}

// GetPermissionService returns the permission service
func (m *AuthModule) GetPermissionService() *permissions.Service {
	return m.permService
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