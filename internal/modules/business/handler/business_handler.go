package handler

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/service"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BusinessHandler struct {
	service *service.BusinessService
}

func NewBusinessHandler(service *service.BusinessService) *BusinessHandler {
	return &BusinessHandler{service: service}
}

// ================================================
// REQUEST MODELS
// ================================================

type UpdateBusinessRequest struct {
	Name        string  `json:"name" example:"Tech Hub Ltd"`
	Category    string  `json:"category" example:"electronics" enums:"retail,fashion,beauty,electronics,food,health"`
	Description string  `json:"description" example:"Quality electronics and gadgets"`
	Logo        string  `json:"logo" example:"https://r2.cloudflare.com/businesses/logo.jpg"`
	Phone       string  `json:"phone" example:"+254745678901"`
	Email       string  `json:"email" example:"info@techhub.com"`
	Address     string  `json:"address" example:"Upper Hill, Nairobi"`
	Location    string  `json:"location" example:"Nairobi CBD"`
	Latitude    float64 `json:"latitude" example:"-1.2921"`
	Longitude   float64 `json:"longitude" example:"36.8219"`
}

type SearchBusinessRequest struct {
	Query    string `form:"q" example:"tech"`
	Category string `form:"category" example:"electronics" enums:"retail,fashion,beauty,electronics,food,health"`
	Page     int    `form:"page" default:"1" minimum:"1" example:"1"`
	PageSize int    `form:"page_size" default:"20" minimum:"1" maximum:"100" example:"20"`
}

// ================================================
// HANDLERS
// ================================================

// GetBusiness godoc
// @Summary Get business by ID
// @Description Get business details by ID
// @Tags Business
// @Produce json
// @Param id path string true "Business ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse{data=models.Business}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{id} [get]
func (h *BusinessHandler) GetBusiness(c *gin.Context) {
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

	business, err := h.service.GetBusinessByID(uid)
	if err != nil {
		response.InternalError(c, "Failed to get business", gin.H{
			"error": err.Error(),
		})
		return
	}

	if business == nil {
		response.NotFound(c, "Business not found", nil)
		return
	}

	response.Success(c, "Business retrieved successfully", business)
}

// GetMyBusiness godoc
// @Summary Get my business
// @Description Get the authenticated user's primary business
// @Tags Business
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=models.Business}
// @Failure 401 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/me [get]
func (h *BusinessHandler) GetMyBusiness(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	business, err := h.service.GetBusinessByUserID(uid)
	if err != nil {
		response.InternalError(c, "Failed to get business", gin.H{
			"error": err.Error(),
		})
		return
	}

	if business == nil {
		response.NotFound(c, "Business not found. Please complete business registration.", nil)
		return
	}

	response.Success(c, "Business retrieved successfully", business)
}

// GetMyBusinesses godoc
// @Summary Get all my businesses
// @Description Get all businesses the authenticated user belongs to
// @Tags Business
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 401 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/my [get]
func (h *BusinessHandler) GetMyBusinesses(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	businesses, err := h.service.GetBusinessesByUserID(uid)
	if err != nil {
		response.InternalError(c, "Failed to get businesses", gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(businesses) == 0 {
		response.NotFound(c, "No businesses found for this user", nil)
		return
	}

	response.Success(c, "Businesses retrieved successfully", gin.H{
		"businesses": businesses,
		"count":      len(businesses),
	})
}

// UpdateBusiness godoc
// @Summary Update business
// @Description Update business details
// @Tags Business
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param request body UpdateBusinessRequest true "Business update details"
// @Success 200 {object} response.BaseResponse{data=models.Business}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{id} [put]
func (h *BusinessHandler) UpdateBusiness(c *gin.Context) {
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

	// Verify the user is a member of this business
	isMember, err := h.service.IsUserMemberOfBusiness(uuid.MustParse(userID), uid)
	if err != nil {
		response.InternalError(c, "Failed to verify membership", gin.H{
			"error": err.Error(),
		})
		return
	}
	if !isMember {
		response.Forbidden(c, "You don't have permission to update this business", nil)
		return
	}

	var req UpdateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Logo != "" {
		updates["logo"] = req.Logo
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Latitude != 0 {
		updates["latitude"] = req.Latitude
	}
	if req.Longitude != 0 {
		updates["longitude"] = req.Longitude
	}

	if len(updates) == 0 {
		response.BadRequest(c, "No fields to update", nil)
		return
	}

	updated, err := h.service.UpdateBusiness(uid, updates)
	if err != nil {
		response.InternalError(c, "Failed to update business", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business updated successfully", updated)
}

// DeleteBusiness godoc
// @Summary Delete business
// @Description Delete a business (owners only)
// @Tags Business
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{id} [delete]
func (h *BusinessHandler) DeleteBusiness(c *gin.Context) {
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

	// Check if user is an owner of this business
	isOwner, err := h.service.IsUserBusinessOwner(uuid.MustParse(userID), uid)
	if err != nil {
		response.InternalError(c, "Failed to verify ownership", gin.H{
			"error": err.Error(),
		})
		return
	}
	if !isOwner {
		response.Forbidden(c, "Only business owners can delete the business", nil)
		return
	}

	if err := h.service.DeleteBusiness(uid); err != nil {
		response.InternalError(c, "Failed to delete business", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business deleted successfully", nil)
}

// SearchBusinesses godoc
// @Summary Search businesses
// @Description Search for businesses by name or category
// @Tags Business
// @Produce json
// @Param q query string false "Search query" example:"tech"
// @Param category query string false "Category filter" example:"electronics" Enums(retail,fashion,beauty,electronics,food,health)
// @Param page query int false "Page number" default(1) example:"1"
// @Param page_size query int false "Page size" default(20) minimum(1) maximum(100) example:"20"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/search [get]
func (h *BusinessHandler) SearchBusinesses(c *gin.Context) {
	var req SearchBusinessRequest
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

	businesses, total, err := h.service.SearchBusinesses(req.Query, req.Category, req.Page, req.PageSize)
	if err != nil {
		response.InternalError(c, "Failed to search businesses", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Businesses retrieved successfully", gin.H{
		"data":        businesses,
		"total":       total,
		"page":        req.Page,
		"page_size":   req.PageSize,
		"total_pages": (total + int64(req.PageSize) - 1) / int64(req.PageSize),
	})
}