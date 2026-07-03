package auth

import (
    "errors"
    "time"

    "github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
    "gorm.io/gorm"
)

type AuthRepository struct {
    db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
    return &AuthRepository{db: db}
}

// User operations
func (r *AuthRepository) CreateUser(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *AuthRepository) GetUserByPhone(phone string) (*models.User, error) {
    var user models.User
    err := r.db.Where("phone_number = ?", phone).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *AuthRepository) GetUserByID(id string) (*models.User, error) {
    var user models.User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// OTP operations
func (r *AuthRepository) SaveOTP(otp *models.OTP) error {
    // Delete old unused OTPs for this email using raw SQL to avoid prepared statement issues
    err := r.db.Exec("DELETE FROM otps WHERE email = ? AND is_used = false", otp.Email).Error
    if err != nil {
        return err
    }
    return r.db.Create(otp).Error
}

func (r *AuthRepository) GetOTPByEmail(email, otpCode string) (*models.OTP, error) {
    var otp models.OTP
    err := r.db.Where(
        "email = ? AND otp_code = ? AND is_used = false AND expires_at > ?",
        email, otpCode, time.Now(),
    ).First(&otp).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("invalid or expired OTP")
        }
        return nil, err
    }
    return &otp, nil
}

func (r *AuthRepository) UpdateOTP(otp *models.OTP) error {
    return r.db.Save(otp).Error
}

func (r *AuthRepository) GetLatestOTP(email string) (*models.OTP, error) {
    var otp models.OTP
    err := r.db.Where("email = ? AND is_used = false", email).
        Order("created_at DESC").
        First(&otp).Error
    if err != nil {
        return nil, err
    }
    return &otp, nil
}