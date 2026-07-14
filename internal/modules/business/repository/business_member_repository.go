package repository

import (
	"errors"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessMemberRepository struct {
    db *gorm.DB
}

func NewBusinessMemberRepository(db *gorm.DB) *BusinessMemberRepository {
    return &BusinessMemberRepository{db: db}
}

// ================================================
// CRUD OPERATIONS
// ================================================

// Create creates a new business member
func (r *BusinessMemberRepository) Create(member *models.BusinessMember) error {
    return r.db.Create(member).Error
}

// GetByUserAndBusiness gets a business member by user and business
func (r *BusinessMemberRepository) GetByUserAndBusiness(userID, businessID uuid.UUID) (*models.BusinessMember, error) {
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

// GetByUserAndBusinessWithRole gets a business member with role check
func (r *BusinessMemberRepository) GetByUserAndBusinessWithRole(userID, businessID uuid.UUID, role string) (*models.BusinessMember, error) {
    var member models.BusinessMember
    err := r.db.Where("user_id = ? AND business_id = ? AND role = ? AND is_active = ?", 
        userID, businessID, role, true).
        First(&member).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &member, nil
}

// GetByInviteToken gets a business member by invite token
func (r *BusinessMemberRepository) GetByInviteToken(token string) (*models.BusinessMember, error) {
    var member models.BusinessMember
    err := r.db.Preload("Business").Preload("User").
        Where("invite_token = ? AND is_active = ?", token, false).
        First(&member).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &member, nil
}

// ================================================
// LIST OPERATIONS
// ================================================

// GetMembersByBusiness gets all members of a business
func (r *BusinessMemberRepository) GetMembersByBusiness(businessID uuid.UUID) ([]models.BusinessMember, error) {
    var members []models.BusinessMember
    err := r.db.Preload("User").
        Where("business_id = ? AND is_active = ?", businessID, true).
        Order("joined_at ASC").
        Find(&members).Error
    return members, err
}

// GetMembersByBusinessWithRole gets all members of a business with a specific role
func (r *BusinessMemberRepository) GetMembersByBusinessWithRole(businessID uuid.UUID, role string) ([]models.BusinessMember, error) {
    var members []models.BusinessMember
    err := r.db.Preload("User").
        Where("business_id = ? AND role = ? AND is_active = ?", businessID, role, true).
        Order("joined_at ASC").
        Find(&members).Error
    return members, err
}

// GetMembersByUser gets all businesses a user belongs to
func (r *BusinessMemberRepository) GetMembersByUser(userID uuid.UUID) ([]models.BusinessMember, error) {
    var members []models.BusinessMember
    err := r.db.Preload("Business").
        Where("user_id = ? AND is_active = ?", userID, true).
        Order("joined_at ASC").
        Find(&members).Error
    return members, err
}

// GetActiveMembersByBusiness gets all active members of a business
func (r *BusinessMemberRepository) GetActiveMembersByBusiness(businessID uuid.UUID) ([]models.BusinessMember, error) {
    var members []models.BusinessMember
    err := r.db.Preload("User").
        Where("business_id = ? AND is_active = ?", businessID, true).
        Find(&members).Error
    return members, err
}

// GetBusinessOwners gets all owners of a business
func (r *BusinessMemberRepository) GetBusinessOwners(businessID uuid.UUID) ([]models.BusinessMember, error) {
    var members []models.BusinessMember
    err := r.db.Preload("User").
        Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleBusinessOwner.String(), true).
        Find(&members).Error
    return members, err
}

// GetBusinessStaff gets all staff of a business
func (r *BusinessMemberRepository) GetBusinessStaff(businessID uuid.UUID) ([]models.BusinessMember, error) {
    var members []models.BusinessMember
    err := r.db.Preload("User").
        Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleBusinessStaff.String(), true).
        Find(&members).Error
    return members, err
}

// ================================================
// UPDATE OPERATIONS
// ================================================

// Update updates a business member
func (r *BusinessMemberRepository) Update(member *models.BusinessMember) error {
    return r.db.Save(member).Error
}

// UpdateRole updates a member's role
func (r *BusinessMemberRepository) UpdateRole(userID, businessID uuid.UUID, role string) error {
    return r.db.Model(&models.BusinessMember{}).
        Where("user_id = ? AND business_id = ?", userID, businessID).
        Update("role", role).Error
}

// ActivateMember activates a business member
func (r *BusinessMemberRepository) ActivateMember(userID, businessID uuid.UUID) error {
    return r.db.Model(&models.BusinessMember{}).
        Where("user_id = ? AND business_id = ?", userID, businessID).
        Updates(map[string]interface{}{
            "is_active": true,
        }).Error
}

// DeactivateMember deactivates a business member
func (r *BusinessMemberRepository) DeactivateMember(userID, businessID uuid.UUID) error {
    return r.db.Model(&models.BusinessMember{}).
        Where("user_id = ? AND business_id = ?", userID, businessID).
        Updates(map[string]interface{}{
            "is_active": false,
        }).Error
}

// UpdateInviteToken updates a member's invite token
func (r *BusinessMemberRepository) UpdateInviteToken(userID, businessID uuid.UUID, token string, expiry *time.Time) error {
    return r.db.Model(&models.BusinessMember{}).
        Where("user_id = ? AND business_id = ?", userID, businessID).
        Updates(map[string]interface{}{
            "invite_token":  token,
            "invite_expiry": expiry,
        }).Error
}

// ================================================
// DELETE OPERATIONS
// ================================================

// Delete soft deletes a business member
func (r *BusinessMemberRepository) Delete(userID, businessID uuid.UUID) error {
    return r.db.Where("user_id = ? AND business_id = ?", userID, businessID).
        Delete(&models.BusinessMember{}).Error
}

// DeleteByBusiness deletes all members of a business
func (r *BusinessMemberRepository) DeleteByBusiness(businessID uuid.UUID) error {
    return r.db.Where("business_id = ?", businessID).
        Delete(&models.BusinessMember{}).Error
}

// ================================================
// CHECK OPERATIONS
// ================================================

// Exists checks if a business member exists
func (r *BusinessMemberRepository) Exists(userID, businessID uuid.UUID) (bool, error) {
    var count int64
    err := r.db.Model(&models.BusinessMember{}).
        Where("user_id = ? AND business_id = ? AND is_active = ?", userID, businessID, true).
        Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// IsOwner checks if a user is an owner of a business
func (r *BusinessMemberRepository) IsOwner(userID, businessID uuid.UUID) (bool, error) {
    var count int64
    err := r.db.Model(&models.BusinessMember{}).
        Where("user_id = ? AND business_id = ? AND role = ? AND is_active = ?", 
            userID, businessID, permissions.RoleBusinessOwner.String(), true).
        Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// IsStaff checks if a user is staff of a business
func (r *BusinessMemberRepository) IsStaff(userID, businessID uuid.UUID) (bool, error) {
    var count int64
    err := r.db.Model(&models.BusinessMember{}).
        Where("user_id = ? AND business_id = ? AND role = ? AND is_active = ?", 
            userID, businessID, permissions.RoleBusinessStaff.String(), true).
        Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// CountMembers counts all members of a business
func (r *BusinessMemberRepository) CountMembers(businessID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.Model(&models.BusinessMember{}).
        Where("business_id = ? AND is_active = ?", businessID, true).
        Count(&count).Error
    return count, err
}

// CountOwners counts all owners of a business
func (r *BusinessMemberRepository) CountOwners(businessID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.Model(&models.BusinessMember{}).
        Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleBusinessOwner.String(), true).
        Count(&count).Error
    return count, err
}