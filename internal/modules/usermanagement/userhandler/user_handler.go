package userhandler

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement/userservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *userservice.UserService
}

func NewUserHandler(service *userservice.UserService) *UserHandler {
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

type UserNotFoundResponse struct {
    Message    string `json:"message" example:"No users found matching the criteria"`
    SearchTerm string `json:"search_term,omitempty" example:"john"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=consumer business_owner admin super_admin" example:"business_owner"`
}

type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active" binding:"required" example:"false"`
}

// ================================================
// HANDLERS
// ================================================

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

    // Set default page number if not provided or invalid
    if req.Page == 0 {
        req.Page = 1
    }

    // Set default page size if not provided or invalid
    if req.PageSize == 0 {
        req.PageSize = 20
    }

    // Validate page size
    if req.PageSize > 100 {
        response.BadRequest(c, "Page size cannot exceed 100", gin.H{
            "max_page_size": 100,
        })
        return
    }

    users, total, err := h.service.GetAllUsers(req.Page, req.PageSize, req.Search)
    if err != nil {
        response.InternalError(c, "Failed to get users", gin.H{
            "error": err.Error(),
        })
        return
    }

    // Return 404 if no users found
    if len(users) == 0 {
        // Build a meaningful response
        notFoundDetails := gin.H{
            "message": "No users found matching the criteria",
        }
        
        // Include search term if provided
        if req.Search != "" {
            notFoundDetails["search_term"] = req.Search
        }
        
        // If page is beyond total pages (e.g., page 5 when only 3 pages exist)
        if total > 0 && req.Page > 0 {
            totalPages := (total + int64(req.PageSize) - 1) / int64(req.PageSize)
            if int64(req.Page) > totalPages {
                notFoundDetails["reason"] = "page exceeds total pages"
                notFoundDetails["total_pages"] = totalPages
                notFoundDetails["total_records"] = total
            }
        }
        
        response.NotFound(c, "No users found", notFoundDetails)
        return
    }

    // Calculate total pages
    totalPages := (total + int64(req.PageSize) - 1) / int64(req.PageSize)

    response.Success(c, "Users retrieved successfully", gin.H{
        "data":         users,
        "total":        total,
        "page":         req.Page,
        "page_size":    req.PageSize,
        "total_pages":  totalPages,
    })
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a specific user by ID (Admin only)
// @Tags User Management
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse{data=models.User}
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	uid, err := uuid.Parse(id) // Parse the user ID from the path parameter
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	user, err := h.service.GetUserByID(uid)
	if err != nil {
		response.InternalError(c, "Failed to get user", gin.H{
			"error": err.Error(),
		})
		return
	}

	if user == nil {
		response.NotFound(c, "User not found", nil)
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get statistics about users (total, consumers, business owners, admins)
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

// UpdateUserRole godoc
// @Summary Update user role
// @Description Update a user's role (Admin only)
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param request body UpdateUserRoleRequest true "New role"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/admin/users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.UpdateUserRole(uid, req.Role); err != nil {
		if err.Error() == "cannot demote the only super_admin" {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalError(c, "Failed to update user role", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "User role updated successfully", nil)
}

// DeleteUser godoc
// @Summary Delete user (soft delete)
// @Description Soft delete a user (Admin only)
// @Tags User Management
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	// Prevent admin from deleting themselves
	if userID == id {
		response.BadRequest(c, "You cannot delete your own account", nil)
		return
	}

	if err := h.service.DeleteUser(uid); err != nil {
		if err.Error() == "cannot delete the only super_admin" {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalError(c, "Failed to delete user", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "User deleted successfully", nil)
}

// HardDeleteUser godoc
// @Summary Permanently delete user
// @Description Permanently delete a user and all related data (Admin only)
// @Tags User Management
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/admin/users/{id}/hard [delete]
func (h *UserHandler) HardDeleteUser(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	// Prevent admin from deleting themselves
	if userID == id {
		response.BadRequest(c, "You cannot delete your own account", nil)
		return
	}

	if err := h.service.HardDeleteUser(uid); err != nil {
		if err.Error() == "cannot delete the only super_admin" {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalError(c, "Failed to permanently delete user", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "User permanently deleted successfully", nil)
}

// SuspendUser godoc
// @Summary Suspend user
// @Description Suspend a user account (Admin only)
// @Tags User Management
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/admin/users/{id}/suspend [put]
func (h *UserHandler) SuspendUser(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	// Prevent admin from suspending themselves
	if userID == id {
		response.BadRequest(c, "You cannot suspend your own account", nil)
		return
	}

	if err := h.service.SuspendUser(uid); err != nil {
		if err.Error() == "cannot suspend the only super_admin" {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalError(c, "Failed to suspend user", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "User suspended successfully", nil)
}

// ActivateUser godoc
// @Summary Activate user
// @Description Activate a suspended user account (Admin only)
// @Tags User Management
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/admin/users/{id}/activate [put]
func (h *UserHandler) ActivateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	if err := h.service.ActivateUser(uid); err != nil {
		response.InternalError(c, "Failed to activate user", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "User activated successfully", nil)
}
