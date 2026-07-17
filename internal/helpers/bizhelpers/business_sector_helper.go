// internal/helpers/bizhelpers/business_sector_helper.go
package bizhelpers

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/bizconstants"
)

// GetBusinessSectorDisplayNames returns display names for business sectors
func GetBusinessSectorDisplayNames() map[string]string {
	return bizconstants.BusinessSectorDisplayNames
}


// IsValidBusinessSector checks if a business sector is valid
func IsValidBusinessSector(sector string) bool {
	for _, valid := range bizconstants.AllBusinessSectors {
		if valid == sector {
			return true
		}
	}
	return false
}

// GetBusinessSectorDisplayName returns display name for a business sector
func GetBusinessSectorDisplayName(sector string) string {
	return bizconstants.BusinessSectorDisplayNames[sector]
}