package business

import (
	"context"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizhandler"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizrepository"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/gin-gonic/gin"
)

type BusinessModule struct {
	businessHandler *bizhandler.BusinessHandler
	memberHandler   *bizhandler.MemberHandler
	businessService *bizservice.BusinessService
}

// NewBusinessModule creates a new business module
func NewBusinessModule(
	cfg *config.Config,
	permService *permissions.Service,
	enforcer *permissions.Enforcer,
) *BusinessModule {
	db := database.GetDB()

	// Initialize repositories
	businessRepo := bizrepository.NewBusinessRepository(db)
	productRepo := bizrepository.NewProductRepository(db)
	memberRepo := bizrepository.NewBusinessMemberRepository(db)
	onboardingRepo := bizrepository.NewOnboardingRepository(db)

	// Initialize service with all dependencies
	businessService := bizservice.NewBusinessService(
		businessRepo,
		productRepo,
		onboardingRepo,
		memberRepo,
		enforcer,
		permService,
		db,
	)

	// Initialize handlers
	businessHandler := bizhandler.NewBusinessHandler(businessService)
	memberHandler := bizhandler.NewMemberHandler(businessService)

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
func (m *BusinessModule) GetBusinessService() *bizservice.BusinessService {
	return m.businessService
}

// GetBusinessHandler returns the business handler
func (m *BusinessModule) GetBusinessHandler() *bizhandler.BusinessHandler {
	return m.businessHandler
}

// GetMemberHandler returns the member handler
func (m *BusinessModule) GetMemberHandler() *bizhandler.MemberHandler {
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