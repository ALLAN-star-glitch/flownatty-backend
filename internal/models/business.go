// internal/models/business.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Business struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	
	// Business Classification
	BusinessTypeID *uuid.UUID `gorm:"type:uuid;index" json:"business_type_id,omitempty"`
	SectorID       *uuid.UUID `gorm:"type:uuid;index" json:"sector_id,omitempty"`
	SubcategoryID  *uuid.UUID `gorm:"type:uuid;index" json:"subcategory_id,omitempty"`
	
	// Establishment Type
	EstablishmentTypeID *uuid.UUID `gorm:"type:uuid;index" json:"establishment_type_id,omitempty"`
	
	// Business Details
	Description string         `gorm:"type:text" json:"description"`
	Logo        string         `gorm:"type:varchar(255)" json:"logo"`
	Phone       string         `gorm:"type:varchar(20);not null" json:"phone"`
	Email       string         `gorm:"type:varchar(255)" json:"email"`
	Address     string         `gorm:"type:text" json:"address"`
	Location    string         `gorm:"type:varchar(255)" json:"location"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	
	// For Markets/Complexes (Physical)
	MarketName  string `gorm:"type:varchar(255)" json:"market_name"`
	StallNumber string `gorm:"type:varchar(50)" json:"stall_number"`
	
	// For Digital/Remote Businesses
	Website     string `gorm:"type:varchar(255)" json:"website"`
	SocialMedia string `gorm:"type:varchar(255)" json:"social_media"` // JSON or comma-separated
	IsRemote    bool   `gorm:"default:false" json:"is_remote"`
	IsDelivery  bool   `gorm:"default:false" json:"is_delivery"`
	
	// Business Metadata
	EmployeeCount   int    `gorm:"default:0" json:"employee_count"`
	YearEstablished int    `gorm:"default:0" json:"year_established"`
	BusinessSize    string `gorm:"type:varchar(20);default:'micro'" json:"business_size"`
	
	// Status
	IsVerified  bool           `gorm:"default:false" json:"is_verified"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	BusinessType      BusinessType        `gorm:"foreignKey:BusinessTypeID" json:"business_type,omitempty"`
	Sector            BusinessSector      `gorm:"foreignKey:SectorID" json:"sector,omitempty"`
	Subcategory       BusinessSubcategory `gorm:"foreignKey:SubcategoryID" json:"subcategory,omitempty"`
	EstablishmentType EstablishmentType   `gorm:"foreignKey:EstablishmentTypeID" json:"establishment_type,omitempty"`
	Members           []BusinessMember    `gorm:"foreignKey:BusinessID" json:"members,omitempty"`
	Products          []Product           `gorm:"foreignKey:BusinessID" json:"products,omitempty"`
	Posts             []Post              `gorm:"foreignKey:BusinessID" json:"posts,omitempty"`
	Orders            []Order             `gorm:"foreignKey:BusinessID" json:"orders,omitempty"`
}

func (Business) TableName() string {
	return "businesses"
}