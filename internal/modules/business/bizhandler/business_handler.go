// internal/modules/business/bizhandler/business_handler.go

package bizhandler

import (
	"context"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BusinessHandler struct {
	service *bizservice.BusinessService
}

func NewBusinessHandler(service *bizservice.BusinessService) *BusinessHandler {
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

type GetBusinessesByTypeRequest struct {
	BusinessType string `form:"business_type" binding:"required" example:"private_company"`
	Page         int    `form:"page" default:"1" minimum:"1" example:"1"`
	PageSize     int    `form:"page_size" default:"20" minimum:"1" maximum:"100" example:"20"`
}

type GetBusinessesBySectorRequest struct {
	Sector   string `form:"sector" binding:"required" example:"technology"`
	Page     int    `form:"page" default:"1" minimum:"1" example:"1"`
	PageSize int    `form:"page_size" default:"20" minimum:"1" maximum:"100" example:"20"`
}

// ================================================
// PUBLIC HANDLERS (No authentication needed)
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

// GetBusinessByName godoc
// @Summary Get business by name
// @Description Get business details by name
// @Tags Business
// @Produce json
// @Param name path string true "Business Name" example:"Tech Hub Ltd"
// @Success 200 {object} response.BaseResponse{data=models.Business}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/name/{name} [get]
func (h *BusinessHandler) GetBusinessByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Business name is required", nil)
		return
	}

	businesses, _, err := h.service.SearchBusinesses(name, "", 1, 1)
	if err != nil {
		response.InternalError(c, "Failed to get business", gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(businesses) == 0 {
		response.NotFound(c, "Business not found", nil)
		return
	}

	response.Success(c, "Business retrieved successfully", businesses[0])
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

// GetBusinessesByType godoc
// @Summary Get businesses by business type
// @Description Get all businesses of a specific type (e.g., private_company, sole_proprietorship)
// @Tags Business
// @Produce json
// @Param business_type query string true "Business type" example:"private_company" Enums(sole_proprietorship,sole_trader,partnership,limited_partnership,private_company,public_company,limited_by_guarantee,cooperative,sacco,ngo,cbo,trust,foundation,faith_based,franchise,epz,special_economic,state_corporation,government_agency)
// @Param page query int false "Page number" default(1) example:"1"
// @Param page_size query int false "Page size" default(20) minimum(1) maximum(100) example:"20"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/by-type [get]
func (h *BusinessHandler) GetBusinessesByType(c *gin.Context) {
	var req GetBusinessesByTypeRequest
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

	businesses, total, err := h.service.GetBusinessesByType(req.BusinessType, req.Page, req.PageSize)
	if err != nil {
		response.InternalError(c, "Failed to get businesses", gin.H{
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

// GetBusinessesBySector godoc
// @Summary Get businesses by sector
// @Description Get all businesses in a specific sector (e.g., technology, financial)
// @Tags Business
// @Produce json
// @Param sector query string true "Sector" example:"technology" Enums(retail,wholesale,fashion,beauty,food,health,agriculture,construction,real_estate,transport,logistics,hospitality,tourism,education,professional,financial,technology,telecom,energy,manufacturing,mining,automotive,entertainment,media,sports,creative,community,environment,security,cleaning,veterinary)
// @Param page query int false "Page number" default(1) example:"1"
// @Param page_size query int false "Page size" default(20) minimum(1) maximum(100) example:"20"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/by-sector [get]
func (h *BusinessHandler) GetBusinessesBySector(c *gin.Context) {
	var req GetBusinessesBySectorRequest
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

	businesses, total, err := h.service.GetBusinessesBySector(req.Sector, req.Page, req.PageSize)
	if err != nil {
		response.InternalError(c, "Failed to get businesses", gin.H{
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

// GetBusinessStats godoc
// @Summary Get business statistics
// @Description Get statistics for a business (requires admin or manager role)
// @Tags Business
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{id}/stats [get]
func (h *BusinessHandler) GetBusinessStats(c *gin.Context) {
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
	stats, err := h.service.GetBusinessStats(ctx, uuid.MustParse(userID), uid)
	if err != nil {
		if err.Error() == "insufficient permissions to view business stats" {
			response.Forbidden(c, "You don't have permission to view business stats", nil)
			return
		}
		response.InternalError(c, "Failed to get business stats", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business stats retrieved successfully", stats)
}

// VerifyBusiness godoc
// @Summary Verify a business
// @Description Verify a business (platform admin only)
// @Tags Business
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{id}/verify [post]
func (h *BusinessHandler) VerifyBusiness(c *gin.Context) {
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
	err = h.service.VerifyBusiness(ctx, uuid.MustParse(userID), uid)
	if err != nil {
		if err.Error() == "only platform admins can verify businesses" {
			response.Forbidden(c, "Only platform admins can verify businesses", nil)
			return
		}
		if err.Error() == "business not found" {
			response.NotFound(c, "Business not found", nil)
			return
		}
		response.InternalError(c, "Failed to verify business", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business verified successfully", nil)
}

// ================================================
// PROTECTED HANDLERS (Authentication required)
// ================================================

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
// @Description Update business details (requires admin or manager role)
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

	ctx := context.Background()
	updated, err := h.service.UpdateBusiness(ctx, uuid.MustParse(userID), uid, updates)
	if err != nil {
		if err.Error() == "insufficient permissions to update business" {
			response.Forbidden(c, "You don't have permission to update this business", nil)
			return
		}
		response.InternalError(c, "Failed to update business", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business updated successfully", updated)
}

// DeleteBusiness godoc
// @Summary Delete business
// @Description Delete a business (requires admin or owner role)
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

	ctx := context.Background()
	if err := h.service.DeleteBusiness(ctx, uuid.MustParse(userID), uid); err != nil {
		if err.Error() == "insufficient permissions to delete business" {
			response.Forbidden(c, "Only business owners can delete the business", nil)
			return
		}
		response.InternalError(c, "Failed to delete business", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business deleted successfully", nil)
}

// ================================================
// REFERENCE DATA HANDLERS
// ================================================

// GetBusinessTypes godoc
// @Summary Get all business types
// @Description Get all business types (legal structures) for dropdowns
// @Tags Business
// @Produce json
// @Success 200 {object} response.BaseResponse{data=[]models.BusinessType}
// @Router /api/v1/business-types [get]
func (h *BusinessHandler) GetBusinessTypes(c *gin.Context) {
	types, err := h.service.GetAllBusinessTypes()
	if err != nil {
		response.InternalError(c, "Failed to get business types", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business types retrieved successfully", types)
}

// GetBusinessSectors godoc
// @Summary Get all business sectors
// @Description Get all business sectors (industry categories) for dropdowns
// @Tags Business
// @Produce json
// @Success 200 {object} response.BaseResponse{data=[]models.BusinessSector}
// @Router /api/v1/business-sectors [get]
func (h *BusinessHandler) GetBusinessSectors(c *gin.Context) {
	sectors, err := h.service.GetAllBusinessSectors()
	if err != nil {
		response.InternalError(c, "Failed to get business sectors", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business sectors retrieved successfully", sectors)
}

// GetBusinessSubcategories godoc
// @Summary Get all business subcategories
// @Description Get all business subcategories for dropdowns
// @Tags Business
// @Produce json
// @Success 200 {object} response.BaseResponse{data=[]models.BusinessSubcategory}
// @Router /api/v1/business-subcategories [get]
func (h *BusinessHandler) GetBusinessSubcategories(c *gin.Context) {
	subcategories, err := h.service.GetAllBusinessSubcategories()
	if err != nil {
		response.InternalError(c, "Failed to get business subcategories", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business subcategories retrieved successfully", subcategories)
}

// GetBusinessSubcategoriesBySector godoc
// @Summary Get business subcategories by sector
// @Description Get business subcategories filtered by sector ID
// @Tags Business
// @Produce json
// @Param sectorId path string true "Sector ID"
// @Success 200 {object} response.BaseResponse{data=[]models.BusinessSubcategory}
// @Router /api/v1/business-subcategories/sector/{sectorId} [get]
func (h *BusinessHandler) GetBusinessSubcategoriesBySector(c *gin.Context) {
	sectorID := c.Param("sectorId")
	if sectorID == "" {
		response.BadRequest(c, "Sector ID is required", nil)
		return
	}

	uid, err := uuid.Parse(sectorID)
	if err != nil {
		response.BadRequest(c, "Invalid sector ID", nil)
		return
	}

	subcategories, err := h.service.GetBusinessSubcategoriesBySector(uid)
	if err != nil {
		response.InternalError(c, "Failed to get business subcategories", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Business subcategories retrieved successfully", subcategories)
}

// GetEstablishmentTypes godoc
// @Summary Get all establishment types
// @Description Get all establishment types for dropdowns
// @Tags Business
// @Produce json
// @Success 200 {object} response.BaseResponse{data=[]models.EstablishmentType}
// @Router /api/v1/establishment-types [get]
func (h *BusinessHandler) GetEstablishmentTypes(c *gin.Context) {
	types, err := h.service.GetAllEstablishmentTypes()
	if err != nil {
		response.InternalError(c, "Failed to get establishment types", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Establishment types retrieved successfully", types)
}