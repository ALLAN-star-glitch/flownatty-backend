package repository

import (
    "errors"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ProductRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
    return &ProductRepository{db: db}
}

// CreateProduct creates a new product
func (r *ProductRepository) CreateProduct(product *models.Product) error {
    return r.db.Create(product).Error
}

// GetProductByID gets a product by ID
func (r *ProductRepository) GetProductByID(id uuid.UUID) (*models.Product, error) {
    var product models.Product
    err := r.db.Preload("Business").Preload("Category").
        Where("id = ? AND is_active = ?", id, true).
        First(&product).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &product, nil
}

// GetProductsByBusinessID gets all products for a business
func (r *ProductRepository) GetProductsByBusinessID(businessID uuid.UUID) ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("Category").
        Where("business_id = ? AND is_active = ?", businessID, true).
        Order("created_at DESC").
        Find(&products).Error
    return products, err
}

// GetProductsByCategoryID gets all products in a category
func (r *ProductRepository) GetProductsByCategoryID(categoryID uuid.UUID) ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("Business").Preload("Category").
        Where("category_id = ? AND is_active = ?", categoryID, true).
        Order("created_at DESC").
        Find(&products).Error
    return products, err
}

// SearchProducts searches products by name or description
func (r *ProductRepository) SearchProducts(query string, limit, offset int) ([]models.Product, int64, error) {
    var products []models.Product
    var total int64
    
    db := r.db.Model(&models.Product{}).
        Preload("Business").
        Preload("Category").
        Where("is_active = ?", true)
    
    if query != "" {
        db = db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
    }
    
    // Get total count
    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Get paginated results
    err := db.Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&products).Error
    
    return products, total, err
}

// UpdateProduct updates a product
func (r *ProductRepository) UpdateProduct(product *models.Product) error {
    return r.db.Save(product).Error
}

// DeleteProduct soft deletes a product
func (r *ProductRepository) DeleteProduct(id uuid.UUID) error {
    return r.db.Delete(&models.Product{}, id).Error
}

// GetProductsByBusinessIDWithPagination gets products for a business with pagination
func (r *ProductRepository) GetProductsByBusinessIDWithPagination(businessID uuid.UUID, limit, offset int) ([]models.Product, int64, error) {
    var products []models.Product
    var total int64
    
    db := r.db.Model(&models.Product{}).
        Preload("Category").
        Where("business_id = ?", businessID)
    
    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := db.Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&products).Error
    
    return products, total, err
}

// GetProductsByIDs gets multiple products by IDs (for cart)
func (r *ProductRepository) GetProductsByIDs(ids []uuid.UUID) ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("Business").
        Where("id IN ? AND is_active = ?", ids, true).
        Find(&products).Error
    return products, err
}