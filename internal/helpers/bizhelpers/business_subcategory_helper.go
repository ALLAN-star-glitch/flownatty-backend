// internal/helpers/bizhelpers/business_subcategory_helper.go
package bizhelpers

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/bizconstants"
)

// GetBusinessSubcategoryDisplayNames returns display names for business subcategories
func GetBusinessSubcategoryDisplayNames() map[string]string {
	return bizconstants.BusinessSubcategoryDisplayNames
}

// IsValidBusinessSubcategory checks if a business subcategory is valid
func IsValidBusinessSubcategory(subcategory string) bool {
	for _, valid := range bizconstants.AllBusinessSubcategories {
		if valid == subcategory {
			return true
		}
	}
	return false
}

// GetBusinessSubcategoryDisplayName returns display name for a business subcategory
func GetBusinessSubcategoryDisplayName(subcategory string) string {
	return bizconstants.BusinessSubcategoryDisplayNames[subcategory]
}