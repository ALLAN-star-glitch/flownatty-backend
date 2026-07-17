// internal/modules/business/bizservice/product_service.go

package bizservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizrepository"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/google/uuid"
)

type ProductService struct {
	repo        *bizrepository.ProductRepository
	enforcer    *permissions.Enforcer
	permService *permissions.Service
}

func NewProductService(
	repo *bizrepository.ProductRepository,
	enforcer *permissions.Enforcer,
	permService *permissions.Service,
) *ProductService {
	return &ProductService{
		repo:        repo,
		enforcer:    enforcer,
		permService: permService,
	}
}

// ================================================
// PRODUCT CRUD OPERATIONS
// ================================================

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, userID, businessID uuid.UUID, product *models.Product) (*models.Product, error) {
	// Check permission
	allowed, err := s.enforcer.Enforce(
		userID.String(),
		permissions.BusinessDomain(businessID.String()),
		permissions.ResourceProduct,
		permissions.ActionCreate,
	)
	if err != nil {
		return nil, fmt.Errorf("permission check failed: %w", err)
	}
	if !allowed {
		return nil, errors.New("insufficient permissions to create product")
	}

	product.BusinessID = businessID
	product.IsActive = true

	if err := s.repo.CreateProduct(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetProductByID gets a product by ID (public - no permission check)
func (s *ProductService) GetProductByID(id uuid.UUID) (*models.Product, error) {
	return s.repo.GetProductByID(id)
}

// GetProductsByBusinessID gets all products for a business
func (s *ProductService) GetProductsByBusinessID(ctx context.Context, userID, businessID uuid.UUID) ([]models.Product, error) {
	// Check permission (read access)
	allowed, err := s.enforcer.Enforce(
		userID.String(),
		permissions.BusinessDomain(businessID.String()),
		permissions.ResourceProduct,
		permissions.ActionRead,
	)
	if err != nil {
		return nil, fmt.Errorf("permission check failed: %w", err)
	}
	if !allowed {
		return nil, errors.New("insufficient permissions to view products")
	}

	return s.repo.GetProductsByBusinessID(businessID)
}

// GetProductsByCategoryID gets all products in a category (public)
func (s *ProductService) GetProductsByCategoryID(categoryID uuid.UUID) ([]models.Product, error) {
	return s.repo.GetProductsByCategoryID(categoryID)
}

// GetProductsBySubcategoryID gets all products in a subcategory (public)
func (s *ProductService) GetProductsBySubcategoryID(subcategoryID uuid.UUID) ([]models.Product, error) {
	return s.repo.GetProductsBySubcategoryID(subcategoryID)
}

// GetProductsByCategoryAndSubcategory gets products filtered by category and optionally subcategory (public)
func (s *ProductService) GetProductsByCategoryAndSubcategory(categoryID uuid.UUID, subcategoryID *uuid.UUID, page, pageSize int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.GetProductsByCategoryAndSubcategory(categoryID, subcategoryID, pageSize, offset)
}

// SearchProducts searches products by name or description (public)
func (s *ProductService) SearchProducts(query string, page, pageSize int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.SearchProducts(query, pageSize, offset)
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(ctx context.Context, userID, businessID, productID uuid.UUID, updates map[string]interface{}) (*models.Product, error) {
	// Check permission
	allowed, err := s.enforcer.Enforce(
		userID.String(),
		permissions.BusinessDomain(businessID.String()),
		permissions.ResourceProduct,
		permissions.ActionUpdate,
	)
	if err != nil {
		return nil, fmt.Errorf("permission check failed: %w", err)
	}
	if !allowed {
		return nil, errors.New("insufficient permissions to update product")
	}

	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Allowed fields for update
	allowedFields := map[string]bool{
		"name":           true,
		"description":    true,
		"price":          true,
		"image_url":      true,
		"stock":          true,
		"is_active":      true,
		"category_id":    true,
		"subcategory_id": true, // ✅ Add this
	}

	for key, value := range updates {
		if allowedFields[key] {
			switch key {
			case "name":
				product.Name = value.(string)
			case "description":
				product.Description = value.(string)
			case "price":
				product.Price = value.(float64)
			case "image_url":
				product.ImageURL = value.(string)
			case "stock":
				product.Stock = int(value.(float64))
			case "is_active":
				product.IsActive = value.(bool)
			case "category_id":
				product.CategoryID = value.(uuid.UUID)
			case "subcategory_id":
				if val, ok := value.(*uuid.UUID); ok {
					product.SubcategoryID = val
				} else if val, ok := value.(uuid.UUID); ok {
					product.SubcategoryID = &val
				}
			}
		}
	}

	if err := s.repo.UpdateProduct(product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, userID, businessID, productID uuid.UUID) error {
	// Check permission
	allowed, err := s.enforcer.Enforce(
		userID.String(),
		permissions.BusinessDomain(businessID.String()),
		permissions.ResourceProduct,
		permissions.ActionDelete,
	)
	if err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}
	if !allowed {
		return errors.New("insufficient permissions to delete product")
	}

	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	return s.repo.DeleteProduct(productID)
}

// GetProductsByBusinessIDWithPagination gets products for a business with pagination
func (s *ProductService) GetProductsByBusinessIDWithPagination(ctx context.Context, userID, businessID uuid.UUID, page, pageSize int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Check permission (read access)
	allowed, err := s.enforcer.Enforce(
		userID.String(),
		permissions.BusinessDomain(businessID.String()),
		permissions.ResourceProduct,
		permissions.ActionRead,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("permission check failed: %w", err)
	}
	if !allowed {
		return nil, 0, errors.New("insufficient permissions to view products")
	}

	offset := (page - 1) * pageSize
	return s.repo.GetProductsByBusinessIDWithPagination(businessID, pageSize, offset)
}

// GetProductsByIDs gets multiple products by IDs (for cart - public)
func (s *ProductService) GetProductsByIDs(ids []uuid.UUID) ([]models.Product, error) {
	return s.repo.GetProductsByIDs(ids)
}