package models

import (
    "time"

    "github.com/ALLAN-star-glitch/flownatty-backend/internal/auth/permissions"
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
    IsEmailVerified   bool       `gorm:"default:false" json:"is_email_verified"`
    VerifiedAt        *time.Time `json:"verified_at"`
    EmailVerifiedAt   *time.Time `json:"email_verified_at"`
    LastLoginAt       *time.Time `json:"last_login_at"`
    LastActiveAt      *time.Time `json:"last_active_at"`
    BusinessID        *uuid.UUID `gorm:"index" json:"business_id,omitempty"`
    Business          *Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
    Permissions       []string   `gorm:"-" json:"-"`
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
    return u.Role == permissions.RoleBusinessOwner.String() || u.BusinessID != nil
}

func (u *User) IsAdmin() bool {
    return u.Role == permissions.RoleAdmin.String() || 
           u.Role == permissions.RoleSuperAdmin.String()
}

func (u *User) IsBusinessStaff() bool {
    return u.Role == permissions.RoleBusinessStaff.String()
}

func (u *User) IsSuperAdmin() bool {
    return u.Role == permissions.RoleSuperAdmin.String()
}

func (u *User) HasBusiness() bool {
    return u.BusinessID != nil
}