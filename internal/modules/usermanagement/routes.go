package usermanagement

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement/userhandler"
	"github.com/gin-gonic/gin"
)

func RegisterUserManagementRoutes(
    router *gin.RouterGroup,
    userHandler *userhandler.UserHandler,
) {
    users := router.Group("/users")
    {
        users.GET("", userHandler.GetAllUsers)
        users.GET("/stats", userHandler.GetUserStats)
        users.GET("/:id", userHandler.GetUserByID)
        users.PUT("/:id/role", userHandler.UpdateUserRole)
        users.DELETE("/:id", userHandler.DeleteUser)
        users.DELETE("/:id/hard", userHandler.HardDeleteUser)
        users.PUT("/:id/suspend", userHandler.SuspendUser)
        users.PUT("/:id/activate", userHandler.ActivateUser)
    }
}