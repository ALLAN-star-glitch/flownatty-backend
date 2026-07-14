package repository

import (
	"errors"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessRepository struct {
    db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) *BusinessRepository {
    return &BusinessRepository{db: db}
}

// ================================================
// BUSINESS CRUD OPERATIONS
// ================================================

// CreateBusiness creates a new business
func (r *BusinessRepository) CreateBusiness(business *models.Business) error {
    return r.db.Create(business).Error
}

// GetBusinessByID gets a business by ID
func (r *BusinessRepository) GetBusinessByID(id uuid.UUID) (*models.Business, error) {
    var business models.Business
    err := r.db.Preload("Members").Preload("Members.User").
        Where("id = ? AND is_active = ?", id, true).First(&business).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &business, nil
}

// GetBusinessByIDWithProducts gets a business with its products
func (r *BusinessRepository) GetBusinessByIDWithProducts(id uuid.UUID) (*models.Business, error) {
    var business models.Business
    err := r.db.Preload("Products", "is_active = ?", true).
        Preload("Members").Preload("Members.User").
        Where("id = ? AND is_active = ?", id, true).
        First(&business).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &business, nil
}

// GetBusinessByIDWithAll gets a business with all relationships loaded
func (r *BusinessRepository) GetBusinessByIDWithAll(id uuid.UUID) (*models.Business, error) {
    var business models.Business
    err := r.db.Preload("Members").Preload("Members.User").
        Preload("Products").
        Preload("Posts").
        Where("id = ? AND is_active = ?", id, true).
        First(&business).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &business, nil
}

// UpdateBusiness updates a business
func (r *BusinessRepository) UpdateBusiness(business *models.Business) error {
    return r.db.Save(business).Error
}

// DeleteBusiness soft deletes a business
func (r *BusinessRepository) DeleteBusiness(id uuid.UUID) error {
    return r.db.Delete(&models.Business{}, id).Error
}

// ================================================
// BUSINESS QUERY OPERATIONS
// ================================================

// SearchBusinesses searches businesses by name or category
func (r *BusinessRepository) SearchBusinesses(query string, category string, limit, offset int) ([]models.Business, int64, error) {
    var businesses []models.Business
    var total int64

    db := r.db.Model(&models.Business{}).Where("is_active = ?", true)

    if query != "" {
        db = db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
    }

    if category != "" && category != "all" {
        db = db.Where("category = ?", category)
    }

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := db.Preload("Members").Preload("Members.User").
        Order("name ASC").
        Limit(limit).
        Offset(offset).
        Find(&businesses).Error

    return businesses, total, err
}

// GetBusinessesByCategory gets all businesses in a category
func (r *BusinessRepository) GetBusinessesByCategory(category string) ([]models.Business, error) {
    var businesses []models.Business
    err := r.db.Preload("Members").Preload("Members.User").
        Where("category = ? AND is_active = ?", category, true).
        Order("name ASC").
        Find(&businesses).Error
    return businesses, err
}

// GetBusinessesWithProducts gets businesses that have at least one product
func (r *BusinessRepository) GetBusinessesWithProducts(limit, offset int) ([]models.Business, int64, error) {
    var businesses []models.Business
    var total int64

    // Subquery to get businesses with at least one product
    subQuery := r.db.Table("products").
        Select("DISTINCT business_id").
        Where("is_active = ?", true)

    db := r.db.Model(&models.Business{}).
        Where("id IN (?) AND is_active = ?", subQuery, true)

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := db.Preload("Products", "is_active = ?", true).
        Preload("Members").Preload("Members.User").
        Order("name ASC").
        Limit(limit).
        Offset(offset).
        Find(&businesses).Error

    return businesses, total, err
}

// GetAllBusinesses gets all active businesses (with pagination)
func (r *BusinessRepository) GetAllBusinesses(limit, offset int) ([]models.Business, int64, error) {
    var businesses []models.Business
    var total int64

    db := r.db.Model(&models.Business{}).Where("is_active = ?", true)

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := db.Preload("Members").Preload("Members.User").
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&businesses).Error

    return businesses, total, err
}

// ================================================
// BUSINESS VERIFICATION OPERATIONS
// ================================================

// UpdateBusinessVerification updates business verification status
func (r *BusinessRepository) UpdateBusinessVerification(id uuid.UUID, isVerified bool) error {
    return r.db.Model(&models.Business{}).
        Where("id = ?", id).
        Update("is_verified", isVerified).Error
}

// ================================================
// BUSINESS STATS OPERATIONS
// ================================================

// GetBusinessStats gets statistics for a business
func (r *BusinessRepository) GetBusinessStats(businessID uuid.UUID) (map[string]interface{}, error) {
    var productCount int64
    var orderCount int64
    var revenue float64
    var memberCount int64

    // Count products
    if err := r.db.Model(&models.Product{}).
        Where("business_id = ? AND is_active = ?", businessID, true).
        Count(&productCount).Error; err != nil {
        return nil, err
    }

    // Count members
    if err := r.db.Model(&models.BusinessMember{}).
        Where("business_id = ? AND is_active = ?", businessID, true).
        Count(&memberCount).Error; err != nil {
        return nil, err
    }

    // Count orders and revenue
    if err := r.db.Model(&models.Order{}).
        Where("business_id = ? AND status != ?", businessID, "cancelled").
        Select("COUNT(*) as count, COALESCE(SUM(total_amount), 0) as revenue").
        Row().Scan(&orderCount, &revenue); err != nil {
        return nil, err
    }

    stats := map[string]interface{}{
        "product_count": productCount,
        "member_count":  memberCount,
        "order_count":   orderCount,
        "revenue":       revenue,
    }

    return stats, nil
}

// GetBusinessStatsByDateRange gets statistics for a business within a date range
func (r *BusinessRepository) GetBusinessStatsByDateRange(businessID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
    var orderCount int64
    var revenue float64

    if err := r.db.Model(&models.Order{}).
        Where("business_id = ? AND status != ? AND created_at BETWEEN ? AND ?", businessID, "cancelled", startDate, endDate).
        Select("COUNT(*) as count, COALESCE(SUM(total_amount), 0) as revenue").
        Row().Scan(&orderCount, &revenue); err != nil {
        return nil, err
    }

    stats := map[string]interface{}{
        "order_count": orderCount,
        "revenue":     revenue,
        "start_date":  startDate,
        "end_date":    endDate,
    }

    return stats, nil
}