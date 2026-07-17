package models

import (
	"log"
	"strings"
	"time"

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
	TwoFactorEnabled  bool       `gorm:"default:true" json:"two_factor_enabled"`
	TwoFactorSecret   string     `json:"-"`

	// Relationships
	BusinessMembers   []BusinessMember `gorm:"foreignKey:UserID" json:"business_members,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.Password != "" {
        //  Check if already hashed (bcrypt hash starts with $2a$)
        if !strings.HasPrefix(u.Password, "$2a$") {
            hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
            if err != nil {
                return err
            }
            u.Password = string(hashed)
        }
    }
    return nil
}

func (u *User) ComparePassword(password string) error {
    //  Add debug logging
    log.Printf("🔍 Comparing password for user: %s", u.Email)
    log.Printf("   Provided password: %s", password)
    log.Printf("   Stored hash: %s...", u.Password[:20])
    
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    
    if err != nil {
        log.Printf("❌ Password comparison failed: %v", err)
    } else {
        log.Printf("✅ Password matched!")
    }
    
    return err
}