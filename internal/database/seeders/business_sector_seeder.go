// internal/database/seeders/business_sector_seeder.go

package seeders

import (
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/bizconstants"
	"gorm.io/gorm"
)

func SeedBusinessSectors(db *gorm.DB) error {
	// Loop through AllBusinessSectors slice
	sortOrder := 1
	for _, sectorKey := range bizconstants.AllBusinessSectors {
		// Get all data from constants
		name := sectorKey // The constant value (e.g., "financial")
		displayName := bizconstants.BusinessSectorDisplayNames[sectorKey] // Display name (e.g., "Financial Services")
		description := bizconstants.BusinessSectorDescriptions[sectorKey]
		icon := bizconstants.BusinessSectorIcons[sectorKey]

		// Insert or update the sector - let PostgreSQL generate the UUID
		query := `
			INSERT INTO business_sectors (id, name, display_name, description, icon, sort_order, is_active, created_at, updated_at)
			VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, ?, NOW(), NOW())
			ON CONFLICT (name) DO UPDATE SET
				display_name = EXCLUDED.display_name,
				description = EXCLUDED.description,
				icon = EXCLUDED.icon,
				sort_order = EXCLUDED.sort_order,
				is_active = EXCLUDED.is_active,
				updated_at = NOW()
			WHERE business_sectors.display_name != EXCLUDED.display_name 
			   OR business_sectors.description != EXCLUDED.description
			   OR business_sectors.sort_order != EXCLUDED.sort_order
		`

		err := db.Exec(query,
			name,
			displayName,
			description,
			icon,
			sortOrder,
			true,
		).Error

		if err != nil {
			return err
		}
		sortOrder++
	}

	log.Printf("✅ Business sectors seeded: %d sectors", len(bizconstants.AllBusinessSectors))
	return nil
}