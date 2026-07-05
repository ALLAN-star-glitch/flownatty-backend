package auth

import (
    "errors"
    "time"

    "github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type AuthRepository struct {
    db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
    return &AuthRepository{db: db}
}

// ================================================
// USER OPERATIONS
// ================================================

// CreateUser creates a new user
func (r *AuthRepository) CreateUser(user *models.User) error {
    return r.db.Create(user).Error
}

// GetUserByEmail finds a user by email
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

// GetUserByPhone finds a user by phone number
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

// GetUserByID finds a user by ID
func (r *AuthRepository) GetUserByID(id string) (*models.User, error) {
    var user models.User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// UpdateUser updates an existing user
func (r *AuthRepository) UpdateUser(user *models.User) error {
    return r.db.Save(user).Error
}

// ================================================
// OTP OPERATIONS
// ================================================

// SaveOTP saves a new OTP and deletes old unused ones
func (r *AuthRepository) SaveOTP(otp *models.OTP) error {
    // Delete old unused OTPs for this email
    err := r.db.Exec("DELETE FROM otps WHERE email = ? AND is_used = false", otp.Email).Error
    if err != nil {
        return err
    }
    return r.db.Create(otp).Error
}

// GetOTPByEmail finds a valid OTP by email and code
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

// UpdateOTP updates an OTP record
func (r *AuthRepository) UpdateOTP(otp *models.OTP) error {
    return r.db.Save(otp).Error
}

// GetLatestOTP gets the most recent unused OTP for an email
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

// ================================================
// REFRESH TOKEN OPERATIONS
// ================================================

// CreateRefreshToken creates a new refresh token
func (r *AuthRepository) CreateRefreshToken(token *models.RefreshToken) error {
    return r.db.Create(token).Error
}

// GetRefreshTokenByToken finds a refresh token by its token string
func (r *AuthRepository) GetRefreshTokenByToken(token string) (*models.RefreshToken, error) {
    var refreshToken models.RefreshToken
    err := r.db.Where("token = ?", token).First(&refreshToken).Error
    if err != nil {
        return nil, err
    }
    return &refreshToken, nil
}

// GetRefreshTokensByUserID finds all refresh tokens for a user
func (r *AuthRepository) GetRefreshTokensByUserID(userID uuid.UUID) ([]models.RefreshToken, error) {
    var tokens []models.RefreshToken
    err := r.db.Where("user_id = ? AND revoked = ?", userID, false).
        Order("created_at DESC").
        Find(&tokens).Error
    if err != nil {
        return nil, err
    }
    return tokens, nil
}

// RevokeRefreshToken revokes a single refresh token
func (r *AuthRepository) RevokeRefreshToken(token string) error {
    return r.db.Model(&models.RefreshToken{}).
        Where("token = ?", token).
        Update("revoked", true).Error
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user
func (r *AuthRepository) RevokeAllUserRefreshTokens(userID uuid.UUID) error {
    return r.db.Model(&models.RefreshToken{}).
        Where("user_id = ?", userID).
        Update("revoked", true).Error
}

// DeleteExpiredRefreshTokens deletes all expired refresh tokens
func (r *AuthRepository) DeleteExpiredRefreshTokens() error {
    return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}

// CleanupRevokedRefreshTokens deletes revoked tokens older than a certain time
func (r *AuthRepository) CleanupRevokedRefreshTokens(olderThan time.Time) error {
    return r.db.Where("revoked = ? AND updated_at < ?", true, olderThan).
        Delete(&models.RefreshToken{}).Error
}

// ================================================
// BUSINESS OPERATIONS
// ================================================

// UpdateUserBusinessID updates a user's business ID
func (r *AuthRepository) UpdateUserBusinessID(userID uuid.UUID, businessID uuid.UUID) error {
    return r.db.Model(&models.User{}).
        Where("id = ?", userID).
        Update("business_id", businessID).Error
}

// GetUserWithBusiness gets a user with their business loaded
func (r *AuthRepository) GetUserWithBusiness(userID string) (*models.User, error) {
    var user models.User
    err := r.db.Preload("Business").Where("id = ?", userID).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ================================================
// UTILITY OPERATIONS
// ================================================

// UserExistsByEmail checks if a user exists by email
func (r *AuthRepository) UserExistsByEmail(email string) (bool, error) {
    var count int64
    err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// UserExistsByPhone checks if a user exists by phone
func (r *AuthRepository) UserExistsByPhone(phone string) (bool, error) {
    var count int64
    err := r.db.Model(&models.User{}).Where("phone_number = ?", phone).Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}