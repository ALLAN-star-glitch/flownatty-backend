package models

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/google/uuid"
)

type Business struct {
    BaseModel
    // REMOVE: UserID field - now in BusinessMember
    Name            string    `gorm:"not null" json:"name"`
    Category        string    `gorm:"not null" json:"category"`
    Description     string    `json:"description"`
    Logo            string    `json:"logo"`
    Phone           string    `gorm:"not null" json:"phone"`
    Email           string    `json:"email"`
    Address         string    `json:"address"`
    Location        string    `json:"location"`
    Latitude        float64   `json:"latitude"`
    Longitude       float64   `json:"longitude"`
    IsVerified      bool      `gorm:"default:false" json:"is_verified"`
    IsActive        bool      `gorm:"default:true" json:"is_active"`
    
    // Relationships
    Members         []BusinessMember `gorm:"foreignKey:BusinessID" json:"members,omitempty"`
    Products        []Product        `gorm:"foreignKey:BusinessID" json:"products,omitempty"`
    Posts           []Post           `gorm:"foreignKey:BusinessID" json:"posts,omitempty"`
    Orders          []Order          `gorm:"foreignKey:BusinessID" json:"orders,omitempty"`
}

// Helper methods
func (b *Business) GetOwnerIDs() []uuid.UUID {
    var ids []uuid.UUID
    for _, member := range b.Members {
        if member.Role == permissions.RoleBusinessOwner.String() && member.IsActive {
            ids = append(ids, member.UserID)
        }
    }
    return ids
}

func (b *Business) GetStaffIDs() []uuid.UUID {
    var ids []uuid.UUID
    for _, member := range b.Members {
        if member.Role == permissions.RoleBusinessStaff.String() && member.IsActive {
            ids = append(ids, member.UserID)
        }
    }
    return ids
}

func (b *Business) IsUserMember(userID uuid.UUID) bool {
    for _, member := range b.Members {
        if member.UserID == userID && member.IsActive {
            return true
        }
    }
    return false
}

func (b *Business) GetUserRole(userID uuid.UUID) string {
    for _, member := range b.Members {
        if member.UserID == userID && member.IsActive {
            return member.Role
        }
    }
    return ""
}