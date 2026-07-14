package models

import (
    "github.com/google/uuid"
)

type Follow struct {
    BaseModel
    UserID     uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_business" json:"user_id"`
    BusinessID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_business" json:"business_id"`
    User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Business   Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
}