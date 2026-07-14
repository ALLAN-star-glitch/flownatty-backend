package authrepo

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
    err := r.db.Preload("BusinessMembers").Preload("BusinessMembers.Business").
        Where("email = ?", email).First(&user).Error
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
    err := r.db.Preload("BusinessMembers").Preload("BusinessMembers.Business").
        Where("phone_number = ?", phone).First(&user).Error
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
    err := r.db.Preload("BusinessMembers").Preload("BusinessMembers.Business").
        Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// GetUserByIDWithBusiness gets a user with their business memberships loaded
func (r *AuthRepository) GetUserByIDWithBusiness(id uuid.UUID) (*models.User, error) {
    var user models.User
    err := r.db.Preload("BusinessMembers").Preload("BusinessMembers.Business").
        Where("id = ?", id).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

// UpdateUser updates an existing user
func (r *AuthRepository) UpdateUser(user *models.User) error {
    return r.db.Save(user).Error
}

// ================================================
// BUSINESS OPERATIONS (Basic CRUD for Auth)
// ================================================

// CreateBusiness creates a new business
func (r *AuthRepository) CreateBusiness(business *models.Business) error {
    return r.db.Create(business).Error
}

// GetBusinessByID gets a business by ID
func (r *AuthRepository) GetBusinessByID(id uuid.UUID) (*models.Business, error) {
    var business models.Business
    err := r.db.Where("id = ? AND is_active = ?", id, true).First(&business).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &business, nil
}

// GetBusinessByUserID gets a business by user ID (first business user is member of)
// Deprecated: Use BusinessMemberRepository.GetMembersByUser instead
func (r *AuthRepository) GetBusinessByUserID(userID uuid.UUID) (*models.Business, error) {
    var member models.BusinessMember
    err := r.db.Preload("Business").
        Where("user_id = ? AND is_active = ?", userID, true).
        First(&member).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &member.Business, nil
}

// UpdateBusiness updates a business
func (r *AuthRepository) UpdateBusiness(business *models.Business) error {
    return r.db.Save(business).Error
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

// ================================================
// ADMIN OPERATIONS
// ================================================

// GetAllUsers gets all users with pagination
func (r *AuthRepository) GetAllUsers(limit, offset int, search string) ([]models.User, int64, error) {
    var users []models.User
    var total int64

    db := r.db.Model(&models.User{})

    // Apply search filter if provided
    if search != "" {
        db = db.Where("email ILIKE ? OR name ILIKE ? OR phone_number ILIKE ?", 
            "%"+search+"%", "%"+search+"%", "%"+search+"%")
    }

    // Count total
    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get paginated results
    err := db.Preload("BusinessMembers").Preload("BusinessMembers.Business").
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&users).Error

    return users, total, err
}

// GetUserStats gets user statistics
func (r *AuthRepository) GetUserStats() (map[string]interface{}, error) {
    var totalUsers int64
    var totalConsumers int64
    var totalBusinessOwners int64

    // Total users
    if err := r.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
        return nil, err
    }

    // Total consumers
    if err := r.db.Model(&models.User{}).Where("role = ?", "consumer").Count(&totalConsumers).Error; err != nil {
        return nil, err
    }

    // Total business owners
    if err := r.db.Model(&models.User{}).Where("role = ?", "business_owner").Count(&totalBusinessOwners).Error; err != nil {
        return nil, err
    }

    stats := map[string]interface{}{
        "total_users":         totalUsers,
        "total_consumers":     totalConsumers,
        "total_business_owners": totalBusinessOwners,
    }

    return stats, nil
}

// ================================================
// BUSINESS MEMBER OPERATIONS (Basic CRUD for Auth)
// ================================================

// CreateBusinessMember creates a new business member
func (r *AuthRepository) CreateBusinessMember(member *models.BusinessMember) error {
    return r.db.Create(member).Error
}

// GetBusinessMemberByUserAndBusiness gets a business member by user and business
func (r *AuthRepository) GetBusinessMemberByUserAndBusiness(userID, businessID uuid.UUID) (*models.BusinessMember, error) {
    var member models.BusinessMember
    err := r.db.Where("user_id = ? AND business_id = ? AND is_active = ?", userID, businessID, true).
        First(&member).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &member, nil
}

// DeleteBusinessMember deletes a business member
func (r *AuthRepository) DeleteBusinessMember(userID, businessID uuid.UUID) error {
    return r.db.Where("user_id = ? AND business_id = ?", userID, businessID).
        Delete(&models.BusinessMember{}).Error
}