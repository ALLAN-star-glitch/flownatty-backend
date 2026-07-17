// internal/helpers/establishmenthelpers/establishment_helper.go
package establishmenthelpers

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/establishmentconstants"
)

// GetEstablishmentTypeDisplayNames returns display names for establishment types
func GetEstablishmentTypeDisplayNames() map[string]string {
	return establishmentconstants.EstablishmentTypeDisplayNames
}

// IsValidEstablishmentType checks if an establishment type is valid
func IsValidEstablishmentType(establishmentType string) bool {
	for _, valid := range establishmentconstants.AllEstablishmentTypes {
		if valid == establishmentType {
			return true
		}
	}
	return false
}

// GetEstablishmentTypeDisplayName returns display name for an establishment type
func GetEstablishmentTypeDisplayName(establishmentType string) string {
	return establishmentconstants.EstablishmentTypeDisplayNames[establishmentType]
}

// GetEstablishmentCategory returns the category for an establishment type
func GetEstablishmentCategory(establishmentType string) string {
	if info, exists := establishmentconstants.EstablishmentTypes[establishmentType]; exists {
		return info.Category
	}
	return ""
}