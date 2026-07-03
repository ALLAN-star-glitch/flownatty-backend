package models

import (
    "time"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type User struct {
    BaseModel
    PhoneNumber       string     `gorm:"uniqueIndex;not null" json:"phone_number"`
    Email             string     `gorm:"uniqueIndex;not null" json:"email"`
    Password string `gorm:"column:password_hash;not null" json:"-"`
    Name              string     `gorm:"not null" json:"name"`
    Avatar            string     `json:"avatar"`
    Role              string     `gorm:"default:consumer" json:"role"`
    IsVerified        bool       `gorm:"default:true" json:"is_verified"`
    IsEmailVerified   bool       `gorm:"default:false" json:"is_email_verified"`
    VerifiedAt        *time.Time `json:"verified_at"`
    EmailVerifiedAt   *time.Time `json:"email_verified_at"`
    LastLoginAt       *time.Time `json:"last_login_at"`
    BusinessID        *uuid.UUID `gorm:"index" json:"business_id,omitempty"`
    Business          *Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
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