package models

import (
	"time"

	"github.com/google/uuid"
)

type BusinessMember struct {
	BaseModel
	BusinessID   uuid.UUID  `gorm:"type:uuid;not null;index;uniqueIndex:idx_business_user" json:"business_id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index;uniqueIndex:idx_business_user" json:"user_id"`
	Role         string     `gorm:"not null;size:30" json:"role"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	JoinedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
	InvitedBy    *uuid.UUID `gorm:"type:uuid" json:"invited_by,omitempty"`

	Business     Business   `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Inviter      *User      `gorm:"foreignKey:InvitedBy" json:"inviter,omitempty"`
}