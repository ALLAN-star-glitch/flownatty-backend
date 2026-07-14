package usermanagement

import (
	"context"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement/handler"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement/repository"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement/service"
	"github.com/gin-gonic/gin"
)

type UserManagementModule struct {
	userHandler *handler.UserHandler
	userService *service.UserService
}

func NewUserManagementModule(cfg *config.Config) *UserManagementModule {
	db := database.GetDB()

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handler
	userHandler := handler.NewUserHandler(userService)

	return &UserManagementModule{
		userHandler: userHandler,
		userService: userService,
	}
}

// SetupRoutes registers all user management routes
func (m *UserManagementModule) SetupRoutes(
	r *gin.RouterGroup,
) {
	RegisterUserManagementRoutes(
		r,
		m.userHandler,
	)
}

// GetUserHandler returns the user handler
func (m *UserManagementModule) GetUserHandler() *handler.UserHandler {
	return m.userHandler
}

// GetUserService returns the user service
func (m *UserManagementModule) GetUserService() *service.UserService {
	return m.userService
}

// Init initializes the module
func (m *UserManagementModule) Init(ctx context.Context) error {
	log.Println("User management module initialized")
	return nil
}

// Close closes the module
func (m *UserManagementModule) Close() {
	log.Println("User management module closed")
}