
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductServiceSubcategory represents subcategories for products and services
type ProductServiceSubcategory struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"type:varchar(50)" json:"icon"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship
	Category ProductServiceCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (ProductServiceSubcategory) TableName() string {
	return "product_service_subcategories"
}