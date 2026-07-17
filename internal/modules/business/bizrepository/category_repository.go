package bizrepository

import (
    "errors"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type CategoryRepository struct {
    db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
    return &CategoryRepository{db: db}
}

// CreateCategory creates a new category
func (r *CategoryRepository) CreateCategory(category *models.ProductServiceCategory) error {
    return r.db.Create(category).Error
}

// GetCategoryByID gets a category by ID
func (r *CategoryRepository) GetCategoryByID(id uuid.UUID) (*models.ProductServiceCategory, error) {
    var category models.ProductServiceCategory
    err := r.db.Where("id = ? AND is_active = ?", id, true).First(&category).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &category, nil
}

// GetCategoryBySlug gets a category by slug
func (r *CategoryRepository) GetCategoryBySlug(slug string) (*models.ProductServiceCategory, error) {
    var category models.ProductServiceCategory
    err := r.db.Where("slug = ? AND is_active = ?", slug, true).First(&category).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &category, nil
}

// GetAllCategories gets all active categories
func (r *CategoryRepository) GetAllCategories() ([]models.ProductServiceCategory, error) {
    var categories []models.ProductServiceCategory
    err := r.db.Where("is_active = ?", true).
        Order("name ASC").
        Find(&categories).Error
    return categories, err
}

// UpdateCategory updates a category
func (r *CategoryRepository) UpdateCategory(category *models.ProductServiceCategory) error {
    return r.db.Save(category).Error
}

// DeleteCategory deletes a category
func (r *CategoryRepository) DeleteCategory(id uuid.UUID) error {
    return r.db.Delete(&models.ProductServiceCategory{}, id).Error
}