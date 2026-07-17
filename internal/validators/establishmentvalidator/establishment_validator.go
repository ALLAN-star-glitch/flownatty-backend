// internal/validators/establishmentvalidator/establishment_validator.go
package establishmentvalidator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/establishmentconstants"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/helpers/establishmenthelpers"
)

// EstablishmentValidator handles validation for establishment data
type EstablishmentValidator struct{}

// NewEstablishmentValidator creates a new establishment validator
func NewEstablishmentValidator() *EstablishmentValidator {
	return &EstablishmentValidator{}
}

// ================================================
// ESTABLISHMENT TYPE VALIDATION (Uses Helpers & Constants)
// ================================================

// ValidateEstablishmentType validates establishment type
func (v *EstablishmentValidator) ValidateEstablishmentType(establishmentType string) error {
	if establishmentType == "" {
		return errors.New("establishment type is required")
	}

	if !establishmenthelpers.IsValidEstablishmentType(establishmentType) {
		return fmt.Errorf("invalid establishment type: '%s'. Allowed types: %v",
			establishmentType, establishmentconstants.AllEstablishmentTypes)
	}

	return nil
}

// AllEstablishmentTypes returns all valid establishment types
func (v *EstablishmentValidator) AllEstablishmentTypes() []string {
	return establishmentconstants.AllEstablishmentTypes
}

// AllEstablishmentTypeDisplayNames returns display names for establishment types
func (v *EstablishmentValidator) AllEstablishmentTypeDisplayNames() map[string]string {
	return establishmenthelpers.GetEstablishmentTypeDisplayNames()
}

// AllEstablishmentTypeCategories returns category groups for establishment types
func (v *EstablishmentValidator) AllEstablishmentTypeCategories() map[string][]string {
	return establishmentconstants.EstablishmentTypeCategories
}

// GetEstablishmentCategory returns the category of an establishment type
func (v *EstablishmentValidator) GetEstablishmentCategory(establishmentType string) string {
	return establishmenthelpers.GetEstablishmentCategory(establishmentType)
}

// ================================================
// MARKET & STALL VALIDATION
// ================================================

// ValidateMarketName validates market name
func (v *EstablishmentValidator) ValidateMarketName(marketName string) error {
	marketName = strings.TrimSpace(marketName)

	if marketName == "" {
		return nil // Market name is optional
	}

	if len(marketName) < 2 {
		return errors.New("market name must be at least 2 characters")
	}

	if len(marketName) > 255 {
		return errors.New("market name must be less than 255 characters")
	}

	return nil
}

// ValidateStallNumber validates stall/unit number
func (v *EstablishmentValidator) ValidateStallNumber(stallNumber string) error {
	stallNumber = strings.TrimSpace(stallNumber)

	if stallNumber == "" {
		return nil // Stall number is optional
	}

	if len(stallNumber) < 1 {
		return errors.New("stall number must be at least 1 character")
	}

	if len(stallNumber) > 50 {
		return errors.New("stall number must be less than 50 characters")
	}

	return nil
}

// ================================================
// WEBSITE & SOCIAL MEDIA VALIDATION
// ================================================

// ValidateWebsite validates website URL
func (v *EstablishmentValidator) ValidateWebsite(website string) error {
	website = strings.TrimSpace(website)

	if website == "" {
		return nil // Website is optional
	}

	// Simple URL validation
	if !strings.HasPrefix(website, "http://") && !strings.HasPrefix(website, "https://") {
		return errors.New("invalid website URL. Include http:// or https://")
	}

	if len(website) > 255 {
		return errors.New("website URL must be less than 255 characters")
	}

	return nil
}

// ValidateSocialMedia validates social media handle or URL
func (v *EstablishmentValidator) ValidateSocialMedia(socialMedia string) error {
	socialMedia = strings.TrimSpace(socialMedia)

	if socialMedia == "" {
		return nil // Social media is optional
	}

	if len(socialMedia) > 255 {
		return errors.New("social media URL must be less than 255 characters")
	}

	return nil
}

// ================================================
// BUSINESS METADATA VALIDATION
// ================================================

// ValidateEmployeeCount validates employee count
func (v *EstablishmentValidator) ValidateEmployeeCount(count int) error {
	if count < 0 {
		return errors.New("employee count cannot be negative")
	}

	if count > 10000 {
		return errors.New("employee count must be less than 10,000")
	}

	return nil
}

// ValidateYearEstablished validates year established
func (v *EstablishmentValidator) ValidateYearEstablished(year int) error {
	if year < 1800 {
		return errors.New("year established must be 1800 or later")
	}

	if year > 2100 {
		return errors.New("year established must be 2100 or earlier")
	}

	return nil
}

// ================================================
// COMPLETE ESTABLISHMENT DATA VALIDATION
// ================================================

// EstablishmentData represents establishment data
type EstablishmentData struct {
	EstablishmentType string `json:"establishment_type"`
	MarketName        string `json:"market_name"`
	StallNumber       string `json:"stall_number"`
	Website           string `json:"website"`
	SocialMedia       string `json:"social_media"`
	EmployeeCount     int    `json:"employee_count"`
	YearEstablished   int    `json:"year_established"`
}

// ValidateComplete validates all establishment data
func (v *EstablishmentValidator) ValidateComplete(data *EstablishmentData) error {
	// 1. Validate Establishment Type
	if err := v.ValidateEstablishmentType(data.EstablishmentType); err != nil {
		return err
	}

	// 2. Validate Market Name (optional)
	if err := v.ValidateMarketName(data.MarketName); err != nil {
		return err
	}

	// 3. Validate Stall Number (optional)
	if err := v.ValidateStallNumber(data.StallNumber); err != nil {
		return err
	}

	// 4. Validate Website (optional)
	if err := v.ValidateWebsite(data.Website); err != nil {
		return err
	}

	// 5. Validate Social Media (optional)
	if err := v.ValidateSocialMedia(data.SocialMedia); err != nil {
		return err
	}

	// 6. Validate Employee Count (optional)
	if err := v.ValidateEmployeeCount(data.EmployeeCount); err != nil {
		return err
	}

	// 7. Validate Year Established (optional)
	if err := v.ValidateYearEstablished(data.YearEstablished); err != nil {
		return err
	}

	return nil
}

// ================================================
// VALIDATION ERRORS HELPERS
// ================================================

// GetEstablishmentTypeError returns user-friendly error for establishment type
func (v *EstablishmentValidator) GetEstablishmentTypeError() string {
	return fmt.Sprintf("Establishment type must be one of: %v", establishmentconstants.AllEstablishmentTypes)
}

// GetMarketNameError returns user-friendly error for market name
func (v *EstablishmentValidator) GetMarketNameError() string {
	return "Market name must be at least 2 characters"
}

// GetStallNumberError returns user-friendly error for stall number
func (v *EstablishmentValidator) GetStallNumberError() string {
	return "Stall number must be at least 1 character"
}

// GetWebsiteError returns user-friendly error for website
func (v *EstablishmentValidator) GetWebsiteError() string {
	return "Invalid website URL. Include http:// or https://"
}