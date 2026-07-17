// internal/models/business.go

package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BusinessSubcategory represents a specific type of business within a sector
type BusinessSubcategory struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SectorID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"sector_id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	DisplayName string         `gorm:"type:varchar(255)" json:"display_name"` // ✅ ADD THIS
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"type:varchar(50)" json:"icon"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship
	Sector BusinessSector `gorm:"foreignKey:SectorID" json:"sector,omitempty"`
}

func (BusinessSubcategory) TableName() string {
	return "business_subcategories"
}