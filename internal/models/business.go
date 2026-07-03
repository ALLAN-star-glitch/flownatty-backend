package models

import (
    "github.com/google/uuid"
)

type Business struct {
    BaseModel
    UserID      uuid.UUID `gorm:"uniqueIndex;not null" json:"user_id"`
    Name        string    `gorm:"not null" json:"name"`
    Category    string    `gorm:"not null" json:"category"`
    Description string    `json:"description"`
    Logo        string    `json:"logo"`
    Phone       string    `gorm:"not null" json:"phone"`
    Email       string    `json:"email"`
    Address     string    `json:"address"`
    IsVerified  bool      `gorm:"default:false" json:"is_verified"`
    IsActive    bool      `gorm:"default:true" json:"is_active"`
}