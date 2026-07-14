package usermanagement


import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement/handler"
	"github.com/gin-gonic/gin"
)

func RegisterUserManagementRoutes(
	router *gin.RouterGroup,
	userHandler *handler.UserHandler,
) {
	// All routes are protected by admin middleware in main.go
	user := router.Group("/users")
	{
		user.GET("", userHandler.GetAllUsers)
		user.GET("/stats", userHandler.GetUserStats)
	}
}