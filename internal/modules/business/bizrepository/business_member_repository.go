package bizrepository

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
    err := r.db.
        Preload("Business").
        Preload("Business.BusinessType").     
        Preload("Business.Sector").           
        Preload("Business.Subcategory").      
        Preload("Business.Subcategory.Sector").
        Preload("Business.EstablishmentType"). 
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

// ================================================
// ROLE-BASED LIST OPERATIONS (Updated with new roles)
// ================================================

// GetBusinessAdmins gets all business admins of a business
func (r *BusinessMemberRepository) GetBusinessAdmins(businessID uuid.UUID) ([]models.BusinessMember, error) {
	var members []models.BusinessMember
	err := r.db.Preload("User").
		Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleBusinessAdmin.String(), true).
		Find(&members).Error
	return members, err
}

// GetProductManagers gets all product managers of a business
func (r *BusinessMemberRepository) GetProductManagers(businessID uuid.UUID) ([]models.BusinessMember, error) {
	var members []models.BusinessMember
	err := r.db.Preload("User").
		Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleProductManager.String(), true).
		Find(&members).Error
	return members, err
}

// GetOrderManagers gets all order managers of a business
func (r *BusinessMemberRepository) GetOrderManagers(businessID uuid.UUID) ([]models.BusinessMember, error) {
	var members []models.BusinessMember
	err := r.db.Preload("User").
		Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleOrderManager.String(), true).
		Find(&members).Error
	return members, err
}

// GetContentManagers gets all content managers of a business
func (r *BusinessMemberRepository) GetContentManagers(businessID uuid.UUID) ([]models.BusinessMember, error) {
	var members []models.BusinessMember
	err := r.db.Preload("User").
		Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleContentManager.String(), true).
		Find(&members).Error
	return members, err
}

// GetServiceManagers gets all service managers of a business
func (r *BusinessMemberRepository) GetServiceManagers(businessID uuid.UUID) ([]models.BusinessMember, error) {
	var members []models.BusinessMember
	err := r.db.Preload("User").
		Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleServiceManager.String(), true).
		Find(&members).Error
	return members, err
}

// GetCustomerSupport gets all customer support members of a business
func (r *BusinessMemberRepository) GetCustomerSupport(businessID uuid.UUID) ([]models.BusinessMember, error) {
	var members []models.BusinessMember
	err := r.db.Preload("User").
		Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleCustomerSupport.String(), true).
		Find(&members).Error
	return members, err
}

// GetBusinessOwners gets all owners of a business
func (r *BusinessMemberRepository) GetBusinessOwners(businessID uuid.UUID) ([]models.BusinessMember, error) {
	var members []models.BusinessMember
	err := r.db.Preload("User").
		Where("business_id = ? AND role = ? AND is_active = ?", businessID, permissions.RoleBusinessAdmin.String(), true).
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

// HasRole checks if a user has a specific role in a business
func (r *BusinessMemberRepository) HasRole(userID, businessID uuid.UUID, role string) (bool, error) {
	var count int64
	err := r.db.Model(&models.BusinessMember{}).
		Where("user_id = ? AND business_id = ? AND role = ? AND is_active = ?",
			userID, businessID, role, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsBusinessAdmin checks if a user is a business admin
func (r *BusinessMemberRepository) IsBusinessAdmin(userID, businessID uuid.UUID) (bool, error) {
	return r.HasRole(userID, businessID, permissions.RoleBusinessAdmin.String())
}

// IsProductManager checks if a user is a product manager
func (r *BusinessMemberRepository) IsProductManager(userID, businessID uuid.UUID) (bool, error) {
	return r.HasRole(userID, businessID, permissions.RoleProductManager.String())
}

// IsOrderManager checks if a user is an order manager
func (r *BusinessMemberRepository) IsOrderManager(userID, businessID uuid.UUID) (bool, error) {
	return r.HasRole(userID, businessID, permissions.RoleOrderManager.String())
}

// IsContentManager checks if a user is a content manager
func (r *BusinessMemberRepository) IsContentManager(userID, businessID uuid.UUID) (bool, error) {
	return r.HasRole(userID, businessID, permissions.RoleContentManager.String())
}

// IsServiceManager checks if a user is a service manager
func (r *BusinessMemberRepository) IsServiceManager(userID, businessID uuid.UUID) (bool, error) {
	return r.HasRole(userID, businessID, permissions.RoleServiceManager.String())
}

// IsCustomerSupport checks if a user is customer support
func (r *BusinessMemberRepository) IsCustomerSupport(userID, businessID uuid.UUID) (bool, error) {
	return r.HasRole(userID, businessID, permissions.RoleCustomerSupport.String())
}

// IsAdminOrOwner checks if a user is a business admin or owner
func (r *BusinessMemberRepository) IsAdminOrOwner(userID, businessID uuid.UUID) (bool, error) {
	admin, err := r.IsBusinessAdmin(userID, businessID)
	if err != nil {
		return false, err
	}
	if admin {
		return true, nil
	}
	return r.IsBusinessOwner(userID, businessID)
}

// IsBusinessOwner checks if a user is an owner of a business
func (r *BusinessMemberRepository) IsBusinessOwner(userID, businessID uuid.UUID) (bool, error) {
	return r.HasRole(userID, businessID, permissions.RoleBusinessAdmin.String())
}

// IsStaff checks if a user is staff of a business (any business role)
func (r *BusinessMemberRepository) IsStaff(userID, businessID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.BusinessMember{}).
		Where("user_id = ? AND business_id = ? AND is_active = ?", userID, businessID, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasAnyBusinessRole checks if a user has any business role (excluding consumer)
func (r *BusinessMemberRepository) HasAnyBusinessRole(userID, businessID uuid.UUID) (bool, error) {
	var count int64
	roles := []string{
		permissions.RoleBusinessAdmin.String(),
		permissions.RoleProductManager.String(),
		permissions.RoleOrderManager.String(),
		permissions.RoleContentManager.String(),
		permissions.RoleServiceManager.String(),
		permissions.RoleCustomerSupport.String(),
	}
	err := r.db.Model(&models.BusinessMember{}).
		Where("user_id = ? AND business_id = ? AND role IN ? AND is_active = ?",
			userID, businessID, roles, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ================================================
// COUNT OPERATIONS
// ================================================

// CountMembers counts all members of a business
func (r *BusinessMemberRepository) CountMembers(businessID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.BusinessMember{}).
		Where("business_id = ? AND is_active = ?", businessID, true).
		Count(&count).Error
	return count, err
}

// CountByRole counts members of a business by role
func (r *BusinessMemberRepository) CountByRole(businessID uuid.UUID, role string) (int64, error) {
	var count int64
	err := r.db.Model(&models.BusinessMember{}).
		Where("business_id = ? AND role = ? AND is_active = ?", businessID, role, true).
		Count(&count).Error
	return count, err
}

// CountAdmins counts business admins
func (r *BusinessMemberRepository) CountAdmins(businessID uuid.UUID) (int64, error) {
	return r.CountByRole(businessID, permissions.RoleBusinessAdmin.String())
}

// CountManagers counts all managers (product, order, content, service)
func (r *BusinessMemberRepository) CountManagers(businessID uuid.UUID) (int64, error) {
	var count int64
	roles := []string{
		permissions.RoleProductManager.String(),
		permissions.RoleOrderManager.String(),
		permissions.RoleContentManager.String(),
		permissions.RoleServiceManager.String(),
	}
	err := r.db.Model(&models.BusinessMember{}).
		Where("business_id = ? AND role IN ? AND is_active = ?", businessID, roles, true).
		Count(&count).Error
	return count, err
}