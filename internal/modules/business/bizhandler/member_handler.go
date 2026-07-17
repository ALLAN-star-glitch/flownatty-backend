package bizhandler

import (
	"context"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberHandler struct {
	service *bizservice.BusinessService
}

func NewMemberHandler(service *bizservice.BusinessService) *MemberHandler {
	return &MemberHandler{service: service}
}

// ================================================
// REQUEST MODELS
// ================================================

type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role   string `json:"role" binding:"required,oneof=business_admin product_manager order_manager content_manager service_manager customer_support" example:"product_manager" enums:"business_admin,product_manager,order_manager,content_manager,service_manager,customer_support"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=business_admin product_manager order_manager content_manager service_manager customer_support" example:"product_manager" enums:"business_admin,product_manager,order_manager,content_manager,service_manager,customer_support"`
}

// ================================================
// HANDLERS
// ================================================

// GetBusinessMembers godoc
// @Summary Get business members
// @Description Get all members of a business
// @Tags Business Members
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse{data=[]models.BusinessMember}
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{id}/members [get]
func (h *MemberHandler) GetBusinessMembers(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Business ID is required", nil)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(c, "Invalid business ID", nil)
		return
	}

	ctx := context.Background()

	// ✅ Pass context to service
	members, err := h.service.GetBusinessMembers(ctx, uuid.MustParse(userID), uid)
	if err != nil {
		if err.Error() == "insufficient permissions to view members" {
			response.Forbidden(c, "You don't have permission to view members", nil)
			return
		}
		response.InternalError(c, "Failed to get business members", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business members retrieved successfully", members)
}

// AddBusinessMember godoc
// @Summary Add business member
// @Description Add a user as a member of a business (admins only)
// @Tags Business Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param request body AddMemberRequest true "Member details"
// @Success 201 {object} response.BaseResponse{data=models.BusinessMember}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 409 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{id}/members [post]
func (h *MemberHandler) AddBusinessMember(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Business ID is required", nil)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(c, "Invalid business ID", nil)
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	memberUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	ctx := context.Background()

	// ✅ Pass context to service
	member, err := h.service.AddBusinessMember(ctx, uuid.MustParse(userID), uid, memberUserID, req.Role)
	if err != nil {
		if err.Error() == "only business admins can add members" {
			response.Forbidden(c, "Only business admins can add members", nil)
			return
		}
		if err.Error() == "user is already a member of this business" {
			response.Conflict(c, err.Error(), nil)
			return
		}
		response.InternalError(c, "Failed to add business member", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Created(c, "Business member added successfully", member)
}

// UpdateBusinessMemberRole godoc
// @Summary Update member role
// @Description Update a member's role in a business (admins only)
// @Tags Business Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param memberId path string true "Member User ID" example:"550e8400-e29b-41d4-a716-446655440001"
// @Param request body UpdateMemberRoleRequest true "New role"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{id}/members/{memberId} [put]
func (h *MemberHandler) UpdateBusinessMemberRole(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	businessID := c.Param("id")
	if businessID == "" {
		response.BadRequest(c, "Business ID is required", nil)
		return
	}

	memberID := c.Param("memberId")
	if memberID == "" {
		response.BadRequest(c, "Member ID is required", nil)
		return
	}

	bizUID, err := uuid.Parse(businessID)
	if err != nil {
		response.BadRequest(c, "Invalid business ID", nil)
		return
	}

	memUID, err := uuid.Parse(memberID)
	if err != nil {
		response.BadRequest(c, "Invalid member ID", nil)
		return
	}

	var req UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := context.Background()

	// ✅ Pass context to service
	err = h.service.UpdateBusinessMemberRole(ctx, uuid.MustParse(userID), bizUID, memUID, req.Role)
	if err != nil {
		if err.Error() == "only business admins can update roles" {
			response.Forbidden(c, "Only business admins can update member roles", nil)
			return
		}
		if err.Error() == "user is not a member of this business" {
			response.NotFound(c, "User is not a member of this business", nil)
			return
		}
		response.InternalError(c, "Failed to update member role", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business member role updated successfully", nil)
}

// RemoveBusinessMember godoc
// @Summary Remove business member
// @Description Remove a member from a business (admins only)
// @Tags Business Members
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param memberId path string true "Member User ID" example:"550e8400-e29b-41d4-a716-446655440001"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{id}/members/{memberId} [delete]
func (h *MemberHandler) RemoveBusinessMember(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	businessID := c.Param("id")
	if businessID == "" {
		response.BadRequest(c, "Business ID is required", nil)
		return
	}

	memberID := c.Param("memberId")
	if memberID == "" {
		response.BadRequest(c, "Member ID is required", nil)
		return
	}

	bizUID, err := uuid.Parse(businessID)
	if err != nil {
		response.BadRequest(c, "Invalid business ID", nil)
		return
	}

	memUID, err := uuid.Parse(memberID)
	if err != nil {
		response.BadRequest(c, "Invalid member ID", nil)
		return
	}

	// Don't allow removing self
	if uuid.MustParse(userID) == memUID {
		response.BadRequest(c, "You cannot remove yourself from the business", nil)
		return
	}

	ctx := context.Background()

	// ✅ Pass context to service
	err = h.service.RemoveBusinessMember(ctx, uuid.MustParse(userID), bizUID, memUID)
	if err != nil {
		if err.Error() == "only business admins can remove members" {
			response.Forbidden(c, "Only business admins can remove members", nil)
			return
		}
		if err.Error() == "cannot remove yourself" {
			response.BadRequest(c, "You cannot remove yourself from the business", nil)
			return
		}
		if err.Error() == "target user is not a member of this business" {
			response.NotFound(c, "User is not a member of this business", nil)
			return
		}
		response.InternalError(c, "Failed to remove business member", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business member removed successfully", nil)
}