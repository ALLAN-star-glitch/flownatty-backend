// internal/models/establishment_type.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EstablishmentType represents the physical/digital type of business
type EstablishmentType struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(50);unique;not null" json:"name"`
	DisplayName string         `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"type:varchar(50)" json:"icon"`
	Category    string         `gorm:"type:varchar(20);default:'physical'" json:"category"` // physical, digital, hybrid, mobile
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Businesses []Business `gorm:"foreignKey:EstablishmentTypeID" json:"businesses,omitempty"`
}

func (EstablishmentType) TableName() string {
	return "establishment_types"
}

// Category constants
const (
	EstablishmentCategoryPhysical = "physical" // Physical storefront
	EstablishmentCategoryDigital  = "digital"  // Online/digital only
	EstablishmentCategoryHybrid   = "hybrid"   // Both physical and online
	EstablishmentCategoryMobile   = "mobile"   // Mobile/street vendors
	EstablishmentCategoryHome     = "home"     // Home-based business
)