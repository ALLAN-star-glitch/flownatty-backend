// internal/helpers/bizhelpers/business_type_helper.go
package bizhelpers

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/bizconstants"
)

// GetBusinessTypeDisplayNames returns display names for business types
func GetBusinessTypeDisplayNames() map[string]string {
	return bizconstants.BusinessTypeDisplayNames
}

// GetBusinessTypeDescriptions returns descriptions for business types
func GetBusinessTypeDescriptions() map[string]string {
	return bizconstants.BusinessTypeDescriptions
}

// IsValidBusinessType checks if a business type is valid
func IsValidBusinessType(businessType string) bool {
	for _, valid := range bizconstants.AllBusinessTypes {
		if valid == businessType {
			return true
		}
	}
	return false
}

// GetBusinessTypeDisplayName returns display name for a business type
func GetBusinessTypeDisplayName(businessType string) string {
	return bizconstants.BusinessTypeDisplayNames[businessType]
}

// GetBusinessTypeDescription returns description for a business type
func GetBusinessTypeDescription(businessType string) string {
	return bizconstants.BusinessTypeDescriptions[businessType]
}