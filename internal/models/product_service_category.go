// internal/models/product_service_category.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductServiceCategory represents categories for products and services
type ProductServiceCategory struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"type:varchar(50)" json:"icon"`
	Type        string         `gorm:"type:varchar(20);not null;index" json:"type"` // "product" or "service"
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subcategories []ProductServiceSubcategory `gorm:"foreignKey:CategoryID" json:"subcategories,omitempty"`
}

func (ProductServiceCategory) TableName() string {
	return "product_service_categories"
}

// Category type constants
const (
	ProductServiceTypeProduct = "product"
	ProductServiceTypeService = "service"
)