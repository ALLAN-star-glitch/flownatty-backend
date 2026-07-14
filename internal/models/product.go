package models

import (
    "github.com/google/uuid"
)

type Product struct {
    BaseModel
    BusinessID   uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
    CategoryID   uuid.UUID `gorm:"type:uuid;not null;index" json:"category_id"`
    Name         string    `gorm:"not null;size:100" json:"name"`
    Description  string    `gorm:"type:text" json:"description"`
    Price        float64   `gorm:"not null;type:decimal(10,2)" json:"price"`
    ImageURL     string    `gorm:"size:255" json:"image_url"`
    Stock        int       `gorm:"default:0" json:"stock"`
    IsActive     bool      `gorm:"default:true" json:"is_active"`
    Business     Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
    Category     Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}