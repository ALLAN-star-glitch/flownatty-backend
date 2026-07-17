// internal/database/seeders/business_subcategory_seeder.go

package seeders

import (
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/bizconstants"
	"gorm.io/gorm"
)

func SeedBusinessSubcategories(db *gorm.DB) error {
	// Build subcategories from constants and insert with ON CONFLICT
	for _, subcategoryName := range bizconstants.AllBusinessSubcategories {
		// Get all data from constants
		name := subcategoryName // The constant value (e.g., "supermarket")
		displayName := bizconstants.BusinessSubcategoryDisplayNames[subcategoryName] // Display name
		description := bizconstants.BusinessSubcategoryDescriptions[subcategoryName]
		sectorKey := bizconstants.BusinessSubcategorySectorMap[subcategoryName]
		icon := bizconstants.BusinessSubcategoryIcons[subcategoryName]

		// Fallback values if display name not found
		if displayName == "" {
			displayName = subcategoryName
		}
		if description == "" {
			description = displayName
		}

		// Get the sector ID for this subcategory
		var sectorID string
		err := db.Raw("SELECT id FROM business_sectors WHERE name = ?", sectorKey).Scan(&sectorID).Error
		if err != nil || sectorID == "" {
			log.Printf("⚠️ Sector not found for subcategory: %s (sector: %s)", displayName, sectorKey)
			continue
		}

		// Insert or update the subcategory - let PostgreSQL generate the UUID
		query := `
			INSERT INTO business_subcategories (id, sector_id, name, display_name, description, icon, is_active, created_at, updated_at)
			VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, ?, NOW(), NOW())
			ON CONFLICT (name, sector_id) DO UPDATE SET
				display_name = EXCLUDED.display_name,
				description = EXCLUDED.description,
				icon = EXCLUDED.icon,
				is_active = EXCLUDED.is_active,
				updated_at = NOW()
			WHERE business_subcategories.display_name != EXCLUDED.display_name 
			   OR business_subcategories.description != EXCLUDED.description
			   OR business_subcategories.icon != EXCLUDED.icon
		`

		err = db.Exec(query,
			sectorID,
			name,    
			displayName, 
			description,
			icon,
			true,
		).Error

		if err != nil {
			return err
		}
	}

	log.Printf("✅ Business subcategories seeded: %d subcategories", len(bizconstants.AllBusinessSubcategories))
	return nil
}