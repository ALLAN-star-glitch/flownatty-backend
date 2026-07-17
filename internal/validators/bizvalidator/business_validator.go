// internal/validators/bizvalidator/business_validator.go
package bizvalidator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/bizconstants"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/establishmentconstants"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/helpers/bizhelpers"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/helpers/establishmenthelpers"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/validation"
)

// BusinessValidator handles validation for business data
type BusinessValidator struct {
	validator *validation.Validator
}

// NewBusinessValidator creates a new business validator
func NewBusinessValidator() *BusinessValidator {
	return &BusinessValidator{
		validator: validation.New(),
	}
}

// ================================================
// BUSINESS TYPE VALIDATION (Uses Helpers)
// ================================================

// ValidateBusinessType validates business type using bizhelpers
func (v *BusinessValidator) ValidateBusinessType(businessType string) error {
	if businessType == "" {
		return errors.New("business type is required")
	}

	if !bizhelpers.IsValidBusinessType(businessType) {
		return fmt.Errorf("invalid business type: '%s'. Allowed types: %v",
			businessType, bizconstants.AllBusinessTypes)
	}

	return nil
}

// AllBusinessTypes returns all valid business types
func (v *BusinessValidator) AllBusinessTypes() []string {
	return bizconstants.AllBusinessTypes
}

// BusinessTypeDisplayNames returns display names for business types
func (v *BusinessValidator) BusinessTypeDisplayNames() map[string]string {
	return bizhelpers.GetBusinessTypeDisplayNames()
}

// ================================================
// BUSINESS CATEGORY VALIDATION (Uses Helpers)
// ================================================

// ValidateBusinessCategory validates business category
func (v *BusinessValidator) ValidateBusinessCategory(category string) error {
	if category == "" {
		return errors.New("business category is required")
	}

	for _, valid := range bizconstants.AllBusinessSectors {
		if valid == category {
			return nil
		}
	}

	return fmt.Errorf("invalid business category: '%s'. Allowed categories: %v",
		category, bizconstants.AllBusinessSectors)
}

// AllBusinessCategories returns all valid business categories
func (v *BusinessValidator) AllBusinessCategories() []string {
	return bizconstants.AllBusinessSectors
}

// ================================================
// BUSINESS SECTOR VALIDATION (Uses Helpers)
// ================================================

// ValidateBusinessSector validates business sector using bizhelpers
func (v *BusinessValidator) ValidateBusinessSector(sector string) error {
	if sector == "" {
		return nil
	}

	if !bizhelpers.IsValidBusinessSector(sector) {
		return fmt.Errorf("invalid business sector: '%s'. Allowed sectors: %v",
			sector, bizconstants.AllBusinessSectors)
	}

	return nil
}

// AllBusinessSectors returns all valid business sectors
func (v *BusinessValidator) AllBusinessSectors() []string {
	return bizconstants.AllBusinessSectors
}

// BusinessSectorDisplayNames returns display names for business sectors
func (v *BusinessValidator) BusinessSectorDisplayNames() map[string]string {
	return bizhelpers.GetBusinessSectorDisplayNames()
}

// ================================================
// BUSINESS SUBCATEGORY VALIDATION (Uses Helpers)
// ================================================

// ValidateBusinessSubcategory validates business subcategory using bizhelpers
func (v *BusinessValidator) ValidateBusinessSubcategory(subcategory string) error {
	if subcategory == "" {
		return nil
	}

	if !bizhelpers.IsValidBusinessSubcategory(subcategory) {
		return fmt.Errorf("invalid business subcategory: '%s'. Allowed subcategories: %v",
			subcategory, bizconstants.AllBusinessSubcategories)
	}

	return nil
}

// AllBusinessSubcategories returns all valid business subcategories
func (v *BusinessValidator) AllBusinessSubcategories() []string {
	return bizconstants.AllBusinessSubcategories
}

// BusinessSubcategoryDisplayNames returns display names for business subcategories
func (v *BusinessValidator) BusinessSubcategoryDisplayNames() map[string]string {
	return bizhelpers.GetBusinessSubcategoryDisplayNames()
}

// ================================================
// ESTABLISHMENT TYPE VALIDATION (Uses Helpers)
// ================================================

// ValidateEstablishmentType validates establishment type using establishmenthelpers
func (v *BusinessValidator) ValidateEstablishmentType(establishmentType string) error {
	if establishmentType == "" {
		return nil
	}

	if !establishmenthelpers.IsValidEstablishmentType(establishmentType) {
		return fmt.Errorf("invalid establishment type: '%s'. Allowed types: %v",
			establishmentType, establishmentconstants.AllEstablishmentTypes)
	}

	return nil
}

// AllEstablishmentTypes returns all valid establishment types
func (v *BusinessValidator) AllEstablishmentTypes() []string {
	return establishmentconstants.AllEstablishmentTypes
}

// EstablishmentTypeDisplayNames returns display names for establishment types
func (v *BusinessValidator) EstablishmentTypeDisplayNames() map[string]string {
	return establishmenthelpers.GetEstablishmentTypeDisplayNames()
}

// ================================================
// BUSINESS PHONE VALIDATION
// ================================================

// ValidateBusinessPhone validates Kenyan phone numbers
func (v *BusinessValidator) ValidateBusinessPhone(phone string) error {
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return errors.New("business phone number is required")
	}

	if !v.validator.Phone.Validate(phone) {
		return fmt.Errorf("invalid business phone number: '%s'", phone)
	}

	return nil
}

// NormalizeBusinessPhone converts phone to standard format +254XXXXXXXXX
func (v *BusinessValidator) NormalizeBusinessPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	return v.validator.Phone.Normalize(phone)
}

// ================================================
// BUSINESS EMAIL VALIDATION
// ================================================

// ValidateBusinessEmail validates business email format
func (v *BusinessValidator) ValidateBusinessEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil
	}

	if !v.validator.Email.Validate(email) {
		return fmt.Errorf("invalid business email: '%s'. %s",
			email, v.validator.Email.ErrorMessage())
	}

	return nil
}

// ValidateBusinessEmailRequired validates business email (required field)
func (v *BusinessValidator) ValidateBusinessEmailRequired(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return errors.New("business email is required")
	}

	return v.ValidateBusinessEmail(email)
}

// ================================================
// BUSINESS NAME VALIDATION
// ================================================

// ValidateBusinessName validates business name
func (v *BusinessValidator) ValidateBusinessName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("business name is required")
	}

	if !v.validator.String.ValidateLength(name, 2, 255) {
		return errors.New("business name must be between 2 and 255 characters")
	}

	if !v.validator.String.ValidateNoSpecialChars(name) {
		return errors.New("business name contains invalid characters")
	}

	return nil
}

// ================================================
// BUSINESS ADDRESS VALIDATION
// ================================================

// ValidateBusinessAddress validates business address
func (v *BusinessValidator) ValidateBusinessAddress(address string) error {
	address = strings.TrimSpace(address)

	if address == "" {
		return errors.New("business address is required")
	}

	if !v.validator.String.ValidateLength(address, 3, 500) {
		return errors.New("business address must be between 3 and 500 characters")
	}

	return nil
}

// ================================================
// BUSINESS DESCRIPTION VALIDATION
// ================================================

// ValidateBusinessDescription validates business description
func (v *BusinessValidator) ValidateBusinessDescription(description string) error {
	description = strings.TrimSpace(description)

	if description == "" {
		return nil
	}

	if !v.validator.String.ValidateLength(description, 10, 2000) {
		return errors.New("business description must be between 10 and 2000 characters")
	}

	return nil
}

// ================================================
// UUID VALIDATION HELPERS
// ================================================

// ValidateUUID validates a UUID format
func (v *BusinessValidator) ValidateUUID(id string) error {
	if id == "" {
		return nil
	}

	if !v.validator.UUID.Validate(id) {
		return errors.New("invalid UUID format")
	}

	return nil
}

// ================================================
// COMPLETE REGISTRATION VALIDATION
// ================================================

// BusinessRegistrationData represents business registration data
type BusinessRegistrationData struct {
	BusinessType        string `json:"business_type"`
	BusinessName        string `json:"business_name"`
	BusinessCategory    string `json:"business_category"`
	BusinessSector      string `json:"business_sector"`
	BusinessSubcategory string `json:"business_subcategory"`
	BusinessPhone       string `json:"business_phone"`
	BusinessEmail       string `json:"business_email"`
	BusinessAddress     string `json:"business_address"`
	BusinessDesc        string `json:"business_description"`
	EstablishmentType   string `json:"establishment_type"`
}

// ValidateRegistration validates all business registration data
func (v *BusinessValidator) ValidateRegistration(data *BusinessRegistrationData) error {
	// 1. Validate Business Type
	if err := v.ValidateBusinessType(data.BusinessType); err != nil {
		return err
	}

	// 2. Validate Business Name
	if err := v.ValidateBusinessName(data.BusinessName); err != nil {
		return err
	}

	// 3. Validate Business Category
	if err := v.ValidateBusinessCategory(data.BusinessCategory); err != nil {
		return err
	}

	// 4. Validate Business Sector (optional)
	if err := v.ValidateBusinessSector(data.BusinessSector); err != nil {
		return err
	}

	// 5. Validate Business Subcategory (optional)
	if err := v.ValidateBusinessSubcategory(data.BusinessSubcategory); err != nil {
		return err
	}

	// 6. Validate Business Phone
	if err := v.ValidateBusinessPhone(data.BusinessPhone); err != nil {
		return err
	}

	// 7. Validate Business Email (optional)
	if err := v.ValidateBusinessEmail(data.BusinessEmail); err != nil {
		return err
	}

	// 8. Validate Business Address
	if err := v.ValidateBusinessAddress(data.BusinessAddress); err != nil {
		return err
	}

	// 9. Validate Business Description (optional)
	if err := v.ValidateBusinessDescription(data.BusinessDesc); err != nil {
		return err
	}

	// 10. Validate Establishment Type (optional)
	if err := v.ValidateEstablishmentType(data.EstablishmentType); err != nil {
		return err
	}

	return nil
}

// ================================================
// BUSINESS UPDATE VALIDATION
// ================================================

// BusinessUpdateData represents business update data
type BusinessUpdateData struct {
	BusinessName        string `json:"business_name,omitempty"`
	BusinessCategory    string `json:"business_category,omitempty"`
	BusinessPhone       string `json:"business_phone,omitempty"`
	BusinessEmail       string `json:"business_email,omitempty"`
	BusinessAddress     string `json:"business_address,omitempty"`
	BusinessDesc        string `json:"business_description,omitempty"`
	SectorID            string `json:"sector_id,omitempty"`
	SubcategoryID       string `json:"subcategory_id,omitempty"`
	EstablishmentTypeID string `json:"establishment_type_id,omitempty"`
	IsActive            *bool  `json:"is_active,omitempty"`
}

// ValidateUpdate validates business update data
func (v *BusinessValidator) ValidateUpdate(data *BusinessUpdateData) error {
	if data.BusinessName != "" {
		if err := v.ValidateBusinessName(data.BusinessName); err != nil {
			return err
		}
	}

	if data.BusinessCategory != "" {
		if err := v.ValidateBusinessCategory(data.BusinessCategory); err != nil {
			return err
		}
	}

	if data.BusinessPhone != "" {
		if err := v.ValidateBusinessPhone(data.BusinessPhone); err != nil {
			return err
		}
	}

	if data.BusinessEmail != "" {
		if err := v.ValidateBusinessEmail(data.BusinessEmail); err != nil {
			return err
		}
	}

	if data.BusinessAddress != "" {
		if err := v.ValidateBusinessAddress(data.BusinessAddress); err != nil {
			return err
		}
	}

	if data.BusinessDesc != "" {
		if err := v.ValidateBusinessDescription(data.BusinessDesc); err != nil {
			return err
		}
	}

	if data.SectorID != "" {
		if err := v.ValidateUUID(data.SectorID); err != nil {
			return err
		}
	}

	if data.SubcategoryID != "" {
		if err := v.ValidateUUID(data.SubcategoryID); err != nil {
			return err
		}
	}

	if data.EstablishmentTypeID != "" {
		if err := v.ValidateUUID(data.EstablishmentTypeID); err != nil {
			return err
		}
	}

	return nil
}

// ================================================
// VALIDATION ERRORS HELPERS
// ================================================

// GetBusinessTypeError returns user-friendly error for business type
func (v *BusinessValidator) GetBusinessTypeError() string {
	return fmt.Sprintf("Business type must be one of: %v", bizconstants.AllBusinessTypes)
}

// GetBusinessCategoryError returns user-friendly error for business category
func (v *BusinessValidator) GetBusinessCategoryError() string {
	return fmt.Sprintf("Business category must be one of: %v", bizconstants.AllBusinessSectors)
}

// GetBusinessSectorError returns user-friendly error for business sector
func (v *BusinessValidator) GetBusinessSectorError() string {
	return fmt.Sprintf("Business sector must be one of: %v", bizconstants.AllBusinessSectors)
}

// GetBusinessSubcategoryError returns user-friendly error for business subcategory
func (v *BusinessValidator) GetBusinessSubcategoryError() string {
	return fmt.Sprintf("Business subcategory must be one of: %v", bizconstants.AllBusinessSubcategories)
}

// GetBusinessPhoneError returns user-friendly error for business phone
func (v *BusinessValidator) GetBusinessPhoneError() string {
	return "Invalid phone number. Use format: 254XXXXXXXXX, +254XXXXXXXXX, 07XXXXXXXX, or 01XXXXXXXX"
}

// GetEstablishmentTypeError returns user-friendly error for establishment type
func (v *BusinessValidator) GetEstablishmentTypeError() string {
	return fmt.Sprintf("Establishment type must be one of: %v", establishmentconstants.AllEstablishmentTypes)
}