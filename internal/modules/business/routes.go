package business

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizhandler"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/gin-gonic/gin"
)

func RegisterBusinessRoutes(
	router *gin.RouterGroup,
	businessHandler *bizhandler.BusinessHandler,
	memberHandler *bizhandler.MemberHandler,
	authMiddleware gin.HandlerFunc,
	enforcer *permissions.Enforcer,
) {
	// ================================================
	// PUBLIC REFERENCE ROUTES (No auth required)
	// ================================================
	reference := router.Group("/")
	{
		reference.GET("/business-types", businessHandler.GetBusinessTypes)
		reference.GET("/business-sectors", businessHandler.GetBusinessSectors)
		reference.GET("/business-subcategories", businessHandler.GetBusinessSubcategories)
		reference.GET("/business-subcategories/sector/:sectorId", businessHandler.GetBusinessSubcategoriesBySector)
		reference.GET("/establishment-types", businessHandler.GetEstablishmentTypes)
	}

	// ================================================
	// PUBLIC BUSINESS ROUTES (No auth required)
	// ================================================
	public := router.Group("/businesses")
	{
		public.GET("/search", businessHandler.SearchBusinesses)
		public.GET("/:id", businessHandler.GetBusiness)
	}

	// ================================================
	// "ME" ROUTES (Auth + Business Access Check)
	// ================================================
	meGroup := router.Group("/businesses")
	meGroup.Use(authMiddleware)                                 // AuthMiddleware from authentication module
	meGroup.Use(permissions.RequireBusinessAccess(enforcer))    // ✅ Check if user has ANY business role
	{
		meGroup.GET("/me", businessHandler.GetMyBusiness)
		meGroup.GET("/my", businessHandler.GetMyBusinesses)
	}

	// ================================================
	// BUSINESS MEMBER ROUTES (Auth + Business Access Check)
	// ================================================
	memberGroup := router.Group("/businesses")
	memberGroup.Use(authMiddleware)
	memberGroup.Use(permissions.RequireBusinessAccess(enforcer))
	{
		memberGroup.GET("/:id/members", memberHandler.GetBusinessMembers)
		memberGroup.POST("/:id/members", memberHandler.AddBusinessMember)
		memberGroup.PUT("/:id/members/:memberId", memberHandler.UpdateBusinessMemberRole)
		memberGroup.DELETE("/:id/members/:memberId", memberHandler.RemoveBusinessMember)
	}

	// ================================================
	// BUSINESS-SPECIFIC ROUTES (Auth + Authorization)
	// ================================================
	protected := router.Group("/businesses")
	protected.Use(authMiddleware)                               // AuthMiddleware from authentication module
	protected.Use(permissions.AuthorizationMiddleware(enforcer)) // Authorization middleware from permissions package
	{
		// Business CRUD (requires specific business permissions)
		protected.PUT("/:id", businessHandler.UpdateBusiness)
		protected.DELETE("/:id", businessHandler.DeleteBusiness)
	}
}