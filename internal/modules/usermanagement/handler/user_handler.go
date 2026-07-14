package handler

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement/service"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// ================================================
// REQUEST MODELS
// ================================================

type GetUsersRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
	Search   string `form:"search" example:"john"`
}

// ================================================
// HANDLERS
// ================================================

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users with pagination and search (Admin only)
// @Tags User Management
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1) example:"1"
// @Param page_size query int false "Page size" default(20) minimum(1) maximum(100) example:"20"
// @Param search query string false "Search by name, email, or phone" example:"john"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/admin/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var req GetUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	users, total, err := h.service.GetAllUsers(req.Page, req.PageSize, req.Search)
	if err != nil {
		response.InternalError(c, "Failed to get users", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Users retrieved successfully", gin.H{
		"data":         users,
		"total":        total,
		"page":         req.Page,
		"page_size":    req.PageSize,
		"total_pages":  (total + int64(req.PageSize) - 1) / int64(req.PageSize),
	})
}

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get statistics about users (total, consumers, business owners)
// @Tags User Management
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/admin/users/stats [get]
func (h *UserHandler) GetUserStats(c *gin.Context) {
	stats, err := h.service.GetUserStats()
	if err != nil {
		response.InternalError(c, "Failed to get user stats", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "User stats retrieved successfully", stats)
}