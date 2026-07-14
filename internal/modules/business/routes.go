package business

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/handler"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/gin-gonic/gin"
)

func RegisterBusinessRoutes(
	router *gin.RouterGroup,
	businessHandler *handler.BusinessHandler,
	memberHandler *handler.MemberHandler,
	authMiddleware gin.HandlerFunc,
	enforcer *permissions.Enforcer,
) {
	// ================================================
	// PUBLIC ROUTES (No auth required)
	// ================================================
	public := router.Group("/businesses")
	{
		public.GET("/search", businessHandler.SearchBusinesses)
		public.GET("/:id", businessHandler.GetBusiness)
	}

	// ================================================
	// PROTECTED ROUTES (Auth + Authorization required)
	// ================================================
	protected := router.Group("/businesses")
	protected.Use(authMiddleware)                               // AuthMiddleware from authentication module
	protected.Use(permissions.AuthorizationMiddleware(enforcer)) // Authorization middleware from permissions package
	{
		// Business CRUD
		protected.GET("/me", businessHandler.GetMyBusiness)
		protected.GET("/my", businessHandler.GetMyBusinesses)
		protected.PUT("/:id", businessHandler.UpdateBusiness)
		protected.DELETE("/:id", businessHandler.DeleteBusiness)

		// Business Members
		protected.GET("/:id/members", memberHandler.GetBusinessMembers)
		protected.POST("/:id/members", memberHandler.AddBusinessMember)
		protected.PUT("/:id/members/:memberId", memberHandler.UpdateBusinessMemberRole)
		protected.DELETE("/:id/members/:memberId", memberHandler.RemoveBusinessMember)
	}
}