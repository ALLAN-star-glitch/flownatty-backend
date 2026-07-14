package models

import (
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/google/uuid"
)

type BusinessMember struct {
    BaseModel
    BusinessID   uuid.UUID  `gorm:"type:uuid;not null;index;uniqueIndex:idx_business_user" json:"business_id"`
    UserID       uuid.UUID  `gorm:"type:uuid;not null;index;uniqueIndex:idx_business_user" json:"user_id"`
    Role         string     `gorm:"not null;size:20" json:"role"` // Uses permissions.Role values
    IsActive     bool       `gorm:"default:true" json:"is_active"`
    JoinedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
    InvitedBy    *uuid.UUID `gorm:"type:uuid" json:"invited_by,omitempty"`
    
    Business     Business   `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
    User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Inviter      *User      `gorm:"foreignKey:InvitedBy" json:"inviter,omitempty"`
}

// Helper methods using permissions roles
func (m *BusinessMember) IsOwner() bool {
    return m.Role == permissions.RoleBusinessOwner.String()
}

func (m *BusinessMember) IsStaff() bool {
    return m.Role == permissions.RoleBusinessStaff.String()
}

func (m *BusinessMember) IsBusinessRole() bool {
    return m.Role == permissions.RoleBusinessOwner.String() || 
           m.Role == permissions.RoleBusinessStaff.String()
}