
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BusinessSector represents the main industry sector of a business


type BusinessSector struct {
    ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Name        string         `gorm:"type:varchar(100);unique;not null" json:"name"`        // Constant: "financial"
    DisplayName string         `gorm:"type:varchar(255)" json:"display_name"`                 // Display: "Financial Services"
    Description string         `gorm:"type:text" json:"description"`
    Icon        string         `gorm:"type:varchar(50)" json:"icon"`
    SortOrder   int            `gorm:"default:0" json:"sort_order"`
    IsActive    bool           `gorm:"default:true" json:"is_active"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BusinessSector) TableName() string {
	return "business_sectors"
}