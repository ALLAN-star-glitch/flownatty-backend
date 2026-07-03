package models

import (
    "time"
    "github.com/google/uuid"
)

type OTP struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    PhoneNumber  string    `gorm:"not null" json:"phone_number"`
    Email        string    `gorm:"not null" json:"email"`
    OTPCode      string    `gorm:"not null" json:"otp_code"`
    Name         string    `gorm:"not null" json:"name"`
    PasswordHash string    `gorm:"not null" json:"-"`
    Role         string    `gorm:"default:consumer" json:"role"`
    Purpose      string    `gorm:"default:signup" json:"purpose"`
    ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
    IsUsed       bool      `gorm:"default:false" json:"is_used"`
    CreatedAt    time.Time `json:"created_at"`
}