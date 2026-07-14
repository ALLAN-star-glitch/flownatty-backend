package models

import (
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	PhoneNumber       string     `gorm:"uniqueIndex;not null" json:"phone_number"`
	Email             string     `gorm:"uniqueIndex;not null" json:"email"`
	Password          string     `gorm:"column:password_hash;not null" json:"-"`
	Name              string     `gorm:"not null" json:"name"`
	Avatar            string     `json:"avatar"`
	Role              string     `gorm:"default:consumer" json:"role"`
	IsVerified        bool       `gorm:"default:true" json:"is_verified"`
	IsEmailVerified   bool       `gorm:"default:true" json:"is_email_verified"`
    IsActive          bool       `gorm:"default:true" json:"is_active"`
	VerifiedAt        *time.Time `json:"verified_at"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	LastActiveAt      *time.Time `json:"last_active_at"`
	
	// REMOVE: BusinessID and Business fields
	// Use BusinessMembers for all business relationships (supports multiple businesses per user)
	
	// 2FA fields
	TwoFactorEnabled  bool       `gorm:"default:true" json:"two_factor_enabled"`
	TwoFactorSecret   string     `json:"-"` // For TOTP (Google Authenticator) - future
	
	// Relationships
	BusinessMembers   []BusinessMember `gorm:"foreignKey:UserID" json:"business_members,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// Helper methods using permissions package constants
func (u *User) IsConsumer() bool {
	return u.Role == permissions.RoleConsumer.String()
}

func (u *User) IsBusinessOwner() bool {
	// Check if user has business_owner role in any business
	for _, member := range u.BusinessMembers {
		if member.Role == permissions.RoleBusinessOwner.String() && member.IsActive {
			return true
		}
	}
	return false
}

func (u *User) IsAdmin() bool {
	return u.Role == permissions.RoleAdmin.String() || 
	       u.Role == permissions.RoleSuperAdmin.String()
}

func (u *User) IsBusinessStaff() bool {
	// Check if user has business_staff role in any business
	for _, member := range u.BusinessMembers {
		if member.Role == permissions.RoleBusinessStaff.String() && member.IsActive {
			return true
		}
	}
	return false
}

func (u *User) IsSuperAdmin() bool {
	return u.Role == permissions.RoleSuperAdmin.String()
}

// HasBusiness checks if user belongs to any business
func (u *User) HasBusiness() bool {
	for _, member := range u.BusinessMembers {
		if member.IsActive {
			return true
		}
	}
	return false
}

// IsMemberOf checks if user is a member of a specific business
func (u *User) IsMemberOf(businessID uuid.UUID) bool {
	for _, member := range u.BusinessMembers {
		if member.BusinessID == businessID && member.IsActive {
			return true
		}
	}
	return false
}

// GetRoleInBusiness returns the user's role in a specific business
func (u *User) GetRoleInBusiness(businessID uuid.UUID) string {
	for _, member := range u.BusinessMembers {
		if member.BusinessID == businessID && member.IsActive {
			return member.Role
		}
	}
	return ""
}

// GetBusinesses returns all businesses the user belongs to
func (u *User) GetBusinesses() []Business {
	var businesses []Business
	for _, member := range u.BusinessMembers {
		if member.IsActive {
			businesses = append(businesses, member.Business)
		}
	}
	return businesses
}

// GetBusinessIDs returns all business IDs the user belongs to
func (u *User) GetBusinessIDs() []uuid.UUID {
	var ids []uuid.UUID
	for _, member := range u.BusinessMembers {
		if member.IsActive {
			ids = append(ids, member.BusinessID)
		}
	}
	return ids
}

// GetPrimaryBusinessID returns the first active business ID (for backward compatibility)
func (u *User) GetPrimaryBusinessID() *uuid.UUID {
	for _, member := range u.BusinessMembers {
		if member.IsActive {
			return &member.BusinessID
		}
	}
	return nil
}