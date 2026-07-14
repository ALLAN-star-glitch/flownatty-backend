package repository

import (
	"errors"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAllUsers - Admin can see all users including soft-deleted
func (r *UserRepository) GetAllUsers(limit, offset int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Use Unscoped() to include soft-deleted users
	db := r.db.Unscoped().Model(&models.User{})

	if search != "" {
		db = db.Where("email ILIKE ? OR name ILIKE ? OR phone_number ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("BusinessMembers").
		Preload("BusinessMembers.Business").
		Preload("BusinessMembers.User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// GetUserByID gets a user by ID (includes soft-deleted)
func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Unscoped().
		Preload("BusinessMembers").
		Preload("BusinessMembers.Business").
		Preload("BusinessMembers.User").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetActiveUserByID gets only active users (excludes soft-deleted)
func (r *UserRepository) GetActiveUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("BusinessMembers").
		Preload("BusinessMembers.Business").
		Preload("BusinessMembers.User").
		Where("id = ? AND is_active = ?", id, true).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserRole updates a user's role
func (r *UserRepository) UpdateUserRole(id uuid.UUID, role string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("role", role).Error
}

// CountUsersByRole counts users with a specific role
func (r *UserRepository) CountUsersByRole(role string) (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}

// DeleteUser soft deletes a user (sets is_active = false and deleted_at)
func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"is_active":  false,
		}).Error
}

// HardDeleteUser permanently deletes a user and all related data
func (r *UserRepository) HardDeleteUser(id uuid.UUID) error {
	// Start a transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Delete from business_members
		if err := tx.Where("user_id = ?", id).Delete(&models.BusinessMember{}).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		// 2. Delete from refresh_tokens
		if err := tx.Where("user_id = ?", id).Delete(&models.RefreshToken{}).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		// 3. Delete from orders
		if err := tx.Where("user_id = ?", id).Delete(&models.Order{}).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		// 4. Delete from follows
		if err := tx.Where("user_id = ?", id).Delete(&models.Follow{}).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		// 5. Delete from casbin_rule (v0 column stores user_id)
		if err := tx.Where("v0 = ?", id.String()).Delete(&models.CasbinRule{}).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		// 6. NOTE: OTPs are stored in Redis, not in the database
		// No need to delete from otps table (it doesn't exist)

		// 7. Finally, delete the user (hard delete)
		if err := tx.Unscoped().Delete(&models.User{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetUserStats gets user statistics
func (r *UserRepository) GetUserStats() (map[string]interface{}, error) {
	var totalUsers int64
	var totalConsumers int64
	var totalBusinessOwners int64
	var totalAdmins int64
	var verifiedUsers int64
	var suspendedUsers int64
	var deletedUsers int64

	// Total users (including soft-deleted)
	if err := r.db.Unscoped().Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	// Total consumers (including soft-deleted)
	if err := r.db.Unscoped().Model(&models.User{}).Where("role = ?", "consumer").Count(&totalConsumers).Error; err != nil {
		return nil, err
	}

	// Total business owners (including soft-deleted)
	if err := r.db.Unscoped().Model(&models.User{}).Where("role = ?", "business_owner").Count(&totalBusinessOwners).Error; err != nil {
		return nil, err
	}

	// Total admins (including soft-deleted)
	if err := r.db.Unscoped().Model(&models.User{}).Where("role IN ?", []string{"admin", "super_admin"}).Count(&totalAdmins).Error; err != nil {
		return nil, err
	}

	// Verified users (including soft-deleted)
	if err := r.db.Unscoped().Model(&models.User{}).Where("is_verified = ?", true).Count(&verifiedUsers).Error; err != nil {
		return nil, err
	}

	// Suspended users (is_active = false, deleted_at IS NULL)
	if err := r.db.Model(&models.User{}).Where("is_active = ? AND deleted_at IS NULL", false).Count(&suspendedUsers).Error; err != nil {
		return nil, err
	}

	// Soft-deleted users (deleted_at IS NOT NULL)
	if err := r.db.Unscoped().Model(&models.User{}).Where("deleted_at IS NOT NULL").Count(&deletedUsers).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_users":           totalUsers,
		"total_consumers":       totalConsumers,
		"total_business_owners": totalBusinessOwners,
		"total_admins":          totalAdmins,
		"verified_users":        verifiedUsers,
		"suspended_users":       suspendedUsers,
		"deleted_users":         deletedUsers,
	}

	return stats, nil
}

// SuspendUser suspends a user (sets is_active to false)
func (r *UserRepository) SuspendUser(id uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ? AND is_active = ?", id, true).
		Update("is_active", false).Error
}

// ActivateUser activates a user (sets is_active to true and clears deleted_at)
// This handles both suspended users AND soft-deleted users
func (r *UserRepository) ActivateUser(id uuid.UUID) error {
	return r.db.Unscoped().
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  true,
			"deleted_at": nil,
		}).Error
}

// GetUserByEmail gets a user by email (includes soft-deleted)
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Unscoped().Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}