package models

import (
    "github.com/google/uuid"
)

type Post struct {
    BaseModel
    BusinessID   uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
    Content      string    `gorm:"type:text;not null" json:"content"`
    ImageURL     string    `gorm:"size:255" json:"image_url"`
    Likes        int       `gorm:"default:0" json:"likes"`
    Comments     int       `gorm:"default:0" json:"comments"`
    IsPublished  bool      `gorm:"default:true" json:"is_published"`
    Business     Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
}