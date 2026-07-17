// internal/modules/business/bizhandler/product_handler.go

package bizhandler

import (
	"context"
	"strconv"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	service *bizservice.ProductService
}

func NewProductHandler(service *bizservice.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// ================================================
// REQUEST MODELS
// ================================================

type CreateProductRequest struct {
	Name           string     `json:"name" binding:"required" example:"iPhone 15 Pro"`
	Description    string     `json:"description" example:"Latest iPhone with advanced camera"`
	Price          float64    `json:"price" binding:"required" example:"120000"`
	ImageURL       string     `json:"image_url" example:"https://example.com/iphone.jpg"`
	Stock          int        `json:"stock" example:"10"`
	CategoryID     string     `json:"category_id" binding:"required" example:"40000000-0000-0000-0000-000000000013"`
	SubcategoryID  *string    `json:"subcategory_id,omitempty" example:"50000000-0000-0000-0000-000000000032"`
}

type UpdateProductRequest struct {
	Name           string  `json:"name" example:"iPhone 15 Pro"`
	Description    string  `json:"description" example:"Latest iPhone with advanced camera"`
	Price          float64 `json:"price" example:"120000"`
	ImageURL       string  `json:"image_url" example:"https://example.com/iphone.jpg"`
	Stock          int     `json:"stock" example:"10"`
	CategoryID     string  `json:"category_id" example:"40000000-0000-0000-0000-000000000013"`
	SubcategoryID  *string `json:"subcategory_id,omitempty" example:"50000000-0000-0000-0000-000000000032"`
	IsActive       *bool   `json:"is_active" example:"true"`
}

type SearchProductRequest struct {
	Query    string `form:"q" example:"phone"`
	Page     int    `form:"page" default:"1" minimum:"1" example:"1"`
	PageSize int    `form:"page_size" default:"20" minimum:"1" maximum:"100" example:"20"`
}

// ================================================
// PUBLIC HANDLERS (No authentication needed)
// ================================================

// GetProduct godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.BaseResponse{data=models.Product}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Product ID is required", nil)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	product, err := h.service.GetProductByID(uid)
	if err != nil {
		response.InternalError(c, "Failed to get product", gin.H{
			"error": err.Error(),
		})
		return
	}

	if product == nil {
		response.NotFound(c, "Product not found", nil)
		return
	}

	response.Success(c, "Product retrieved successfully", product)
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Get all products in a category
// @Tags Products
// @Produce json
// @Param categoryId path string true "Category ID"
// @Success 200 {object} response.BaseResponse{data=[]models.Product}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/products/category/{categoryId} [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryID := c.Param("categoryId")
	if categoryID == "" {
		response.BadRequest(c, "Category ID is required", nil)
		return
	}

	uid, err := uuid.Parse(categoryID)
	if err != nil {
		response.BadRequest(c, "Invalid category ID", nil)
		return
	}

	products, err := h.service.GetProductsByCategoryID(uid)
	if err != nil {
		response.InternalError(c, "Failed to get products", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Products retrieved successfully", products)
}

// GetProductsBySubcategory godoc
// @Summary Get products by subcategory
// @Description Get all products in a specific subcategory
// @Tags Products
// @Produce json
// @Param subcategoryId path string true "Subcategory ID"
// @Success 200 {object} response.BaseResponse{data=[]models.Product}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/products/subcategory/{subcategoryId} [get]
func (h *ProductHandler) GetProductsBySubcategory(c *gin.Context) {
	subcategoryID := c.Param("subcategoryId")
	if subcategoryID == "" {
		response.BadRequest(c, "Subcategory ID is required", nil)
		return
	}

	uid, err := uuid.Parse(subcategoryID)
	if err != nil {
		response.BadRequest(c, "Invalid subcategory ID", nil)
		return
	}

	products, err := h.service.GetProductsBySubcategoryID(uid)
	if err != nil {
		response.InternalError(c, "Failed to get products", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Products retrieved successfully", products)
}

// GetProductsByCategoryAndSubcategory godoc
// @Summary Get products by category and optional subcategory
// @Description Get products filtered by category and optionally subcategory with pagination
// @Tags Products
// @Produce json
// @Param categoryId path string true "Category ID"
// @Param subcategoryId query string false "Subcategory ID (optional)"
// @Param page query int false "Page number" default(1) example:"1"
// @Param page_size query int false "Page size" default(20) minimum(1) maximum(100) example:"20"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/products/category/{categoryId}/filter [get]
func (h *ProductHandler) GetProductsByCategoryAndSubcategory(c *gin.Context) {
	categoryID := c.Param("categoryId")
	if categoryID == "" {
		response.BadRequest(c, "Category ID is required", nil)
		return
	}

	uid, err := uuid.Parse(categoryID)
	if err != nil {
		response.BadRequest(c, "Invalid category ID", nil)
		return
	}

	var subcategoryID *uuid.UUID
	if subcategoryStr := c.Query("subcategoryId"); subcategoryStr != "" {
		parsed, err := uuid.Parse(subcategoryStr)
		if err == nil {
			subcategoryID = &parsed
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	products, total, err := h.service.GetProductsByCategoryAndSubcategory(uid, subcategoryID, page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to get products", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Products retrieved successfully", gin.H{
		"data":        products,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// SearchProducts godoc
// @Summary Search products
// @Description Search products by name or description
// @Tags Products
// @Produce json
// @Param q query string false "Search query" example:"phone"
// @Param page query int false "Page number" default(1) example:"1"
// @Param page_size query int false "Page size" default(20) minimum(1) maximum(100) example:"20"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	var req SearchProductRequest
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

	products, total, err := h.service.SearchProducts(req.Query, req.Page, req.PageSize)
	if err != nil {
		response.InternalError(c, "Failed to search products", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Products retrieved successfully", gin.H{
		"data":        products,
		"total":       total,
		"page":        req.Page,
		"page_size":   req.PageSize,
		"total_pages": (total + int64(req.PageSize) - 1) / int64(req.PageSize),
	})
}

// ================================================
// PROTECTED HANDLERS (Authentication required)
// ================================================

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product for a business (requires product manager or admin role)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param businessId path string true "Business ID"
// @Param request body CreateProductRequest true "Product details"
// @Success 201 {object} response.BaseResponse{data=models.Product}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{businessId}/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	businessID := c.Param("businessId")
	if businessID == "" {
		response.BadRequest(c, "Business ID is required", nil)
		return
	}

	bizUID, err := uuid.Parse(businessID)
	if err != nil {
		response.BadRequest(c, "Invalid business ID", nil)
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		response.BadRequest(c, "Invalid category ID", nil)
		return
	}

	var subcategoryID *uuid.UUID
	if req.SubcategoryID != nil && *req.SubcategoryID != "" {
		parsed, err := uuid.Parse(*req.SubcategoryID)
		if err != nil {
			response.BadRequest(c, "Invalid subcategory ID", nil)
			return
		}
		subcategoryID = &parsed
	}

	product := &models.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		ImageURL:      req.ImageURL,
		Stock:         req.Stock,
		CategoryID:    categoryID,
		SubcategoryID: subcategoryID,
	}

	ctx := context.Background()
	created, err := h.service.CreateProduct(ctx, uuid.MustParse(userID), bizUID, product)
	if err != nil {
		if err.Error() == "insufficient permissions to create product" {
			response.Forbidden(c, "You don't have permission to create products for this business", nil)
			return
		}
		response.InternalError(c, "Failed to create product", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Created(c, "Product created successfully", created)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product (requires product manager or admin role)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param businessId path string true "Business ID"
// @Param productId path string true "Product ID"
// @Param request body UpdateProductRequest true "Product update details"
// @Success 200 {object} response.BaseResponse{data=models.Product}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{businessId}/products/{productId} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	businessID := c.Param("businessId")
	if businessID == "" {
		response.BadRequest(c, "Business ID is required", nil)
		return
	}

	productID := c.Param("productId")
	if productID == "" {
		response.BadRequest(c, "Product ID is required", nil)
		return
	}

	bizUID, err := uuid.Parse(businessID)
	if err != nil {
		response.BadRequest(c, "Invalid business ID", nil)
		return
	}

	prodUID, err := uuid.Parse(productID)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	var req UpdateProductRequest
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
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price != 0 {
		updates["price"] = req.Price
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}
	if req.Stock != 0 {
		updates["stock"] = req.Stock
	}
	if req.CategoryID != "" {
		categoryID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			response.BadRequest(c, "Invalid category ID", nil)
			return
		}
		updates["category_id"] = categoryID
	}
	if req.SubcategoryID != nil && *req.SubcategoryID != "" {
		parsed, err := uuid.Parse(*req.SubcategoryID)
		if err != nil {
			response.BadRequest(c, "Invalid subcategory ID", nil)
			return
		}
		updates["subcategory_id"] = &parsed
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		response.BadRequest(c, "No fields to update", nil)
		return
	}

	ctx := context.Background()
	updated, err := h.service.UpdateProduct(ctx, uuid.MustParse(userID), bizUID, prodUID, updates)
	if err != nil {
		if err.Error() == "insufficient permissions to update product" {
			response.Forbidden(c, "You don't have permission to update this product", nil)
			return
		}
		if err.Error() == "product not found" {
			response.NotFound(c, "Product not found", nil)
			return
		}
		response.InternalError(c, "Failed to update product", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Product updated successfully", updated)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product (requires product manager or admin role)
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Param businessId path string true "Business ID"
// @Param productId path string true "Product ID"
// @Success 200 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{businessId}/products/{productId} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	businessID := c.Param("businessId")
	if businessID == "" {
		response.BadRequest(c, "Business ID is required", nil)
		return
	}

	productID := c.Param("productId")
	if productID == "" {
		response.BadRequest(c, "Product ID is required", nil)
		return
	}

	bizUID, err := uuid.Parse(businessID)
	if err != nil {
		response.BadRequest(c, "Invalid business ID", nil)
		return
	}

	prodUID, err := uuid.Parse(productID)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	ctx := context.Background()
	err = h.service.DeleteProduct(ctx, uuid.MustParse(userID), bizUID, prodUID)
	if err != nil {
		if err.Error() == "insufficient permissions to delete product" {
			response.Forbidden(c, "You don't have permission to delete this product", nil)
			return
		}
		if err.Error() == "product not found" {
			response.NotFound(c, "Product not found", nil)
			return
		}
		response.InternalError(c, "Failed to delete product", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Product deleted successfully", nil)
}

// GetBusinessProducts godoc
// @Summary Get business products
// @Description Get all products for a business (requires authentication)
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Param businessId path string true "Business ID"
// @Param page query int false "Page number" default(1) example:"1"
// @Param page_size query int false "Page size" default(20) minimum(1) maximum(100) example:"20"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 401 {object} response.BaseResponse
// @Failure 403 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/businesses/{businessId}/products [get]
func (h *ProductHandler) GetBusinessProducts(c *gin.Context) {
	userID := c.GetString(permissions.ContextKeyUserID)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated", nil)
		return
	}

	businessID := c.Param("businessId")
	if businessID == "" {
		response.BadRequest(c, "Business ID is required", nil)
		return
	}

	bizUID, err := uuid.Parse(businessID)
	if err != nil {
		response.BadRequest(c, "Invalid business ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	ctx := context.Background()
	products, total, err := h.service.GetProductsByBusinessIDWithPagination(ctx, uuid.MustParse(userID), bizUID, page, pageSize)
	if err != nil {
		if err.Error() == "insufficient permissions to view products" {
			response.Forbidden(c, "You don't have permission to view products for this business", nil)
			return
		}
		response.InternalError(c, "Failed to get products", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(c, "Products retrieved successfully", gin.H{
		"data":        products,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}