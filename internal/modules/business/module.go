package business

import (
	"context"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/handler"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/repository"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/service"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/gin-gonic/gin"
)

type BusinessModule struct {
	businessHandler *handler.BusinessHandler
	memberHandler   *handler.MemberHandler
	businessService *service.BusinessService
}

// NewBusinessModule creates a new business module
func NewBusinessModule(cfg *config.Config, permService *permissions.Service) *BusinessModule {
	db := database.GetDB()

	// Initialize repositories
	businessRepo := repository.NewBusinessRepository(db)
	productRepo := repository.NewProductRepository(db)
	memberRepo := repository.NewBusinessMemberRepository(db)

	// Initialize services
	businessService := service.NewBusinessService(
		businessRepo,
		productRepo,
		memberRepo,
		db,
	)

	// Initialize handlers
	businessHandler := handler.NewBusinessHandler(businessService)
	memberHandler := handler.NewMemberHandler(businessService)

	return &BusinessModule{
		businessHandler: businessHandler,
		memberHandler:   memberHandler,
		businessService: businessService,
	}
}

// SetupRoutes registers all business routes with the router
func (m *BusinessModule) SetupRoutes(
	r *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	enforcer *permissions.Enforcer,
) {
	RegisterBusinessRoutes(
		r,
		m.businessHandler,
		m.memberHandler,
		authMiddleware,
		enforcer,
	)
}

// GetBusinessService returns the business service
func (m *BusinessModule) GetBusinessService() *service.BusinessService {
	return m.businessService
}

// GetBusinessHandler returns the business handler
func (m *BusinessModule) GetBusinessHandler() *handler.BusinessHandler {
	return m.businessHandler
}

// GetMemberHandler returns the member handler
func (m *BusinessModule) GetMemberHandler() *handler.MemberHandler {
	return m.memberHandler
}

// Init initializes the module
func (m *BusinessModule) Init(ctx context.Context) error {
	log.Println("Business module initialized")
	return nil
}

// Close closes the module and cleans up resources
func (m *BusinessModule) Close() {
	log.Println("Business module closed")
}