// internal/modules/business/bizrepository/onboarding_repository.go

package bizrepository

import (
	"errors"
	"fmt"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OnboardingRepository struct {
	db *gorm.DB
}

func NewOnboardingRepository(db *gorm.DB) *OnboardingRepository {
	return &OnboardingRepository{db: db}
}

// OnboardingRequest represents the business data collected during registration
type OnboardingRequest struct {
	BusinessType    string     `json:"business_type"`
	BusinessName    string     `json:"business_name"`
	BusinessCategory string    `json:"business_category"`
	BusinessPhone   string     `json:"business_phone"`
	BusinessEmail   string     `json:"business_email"`
	BusinessAddress string     `json:"business_address"`
	BusinessTypeID  *uuid.UUID `json:"-"`
	SectorID        *uuid.UUID `json:"-"`
}

// GetSectorIDByName gets the sector ID from the sector name
func (r *OnboardingRepository) GetSectorIDByName(name string) (*uuid.UUID, error) {
	var sector models.BusinessSector
	err := r.db.Where("LOWER(name) = LOWER(?) AND is_active = ?", name, true).First(&sector).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sector.ID, nil
}

// CreateBusinessWithAdmin creates a business and adds the user as business admin - DURING REGISTRATION
func (r *OnboardingRepository) CreateBusinessWithAdminInit(
	business *models.Business,
	adminID uuid.UUID,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create the business
		if err := tx.Create(business).Error; err != nil {
			return fmt.Errorf("failed to create business: %w", err)
		}

		// 2. Add user as business admin
		member := &models.BusinessMember{
			BusinessID: business.ID,
			UserID:     adminID,
			Role:       permissions.RoleBusinessAdmin.String(),
			IsActive:   true,
			JoinedAt:   time.Now(),
		}

		if err := tx.Create(member).Error; err != nil {
			return fmt.Errorf("failed to add business admin: %w", err)
		}

		return nil
	})
}