// internal/database/seeders/business_type_seeder.go

package seeders

import (
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/bizconstants"
	"gorm.io/gorm"
)

func SeedBusinessTypes(db *gorm.DB) error {
	// Loop through the slice directly
	for _, businessType := range bizconstants.AllBusinessTypes {
		// Get all data from constants
		name := businessType // The constant value (e.g., "private_company")
		displayName := bizconstants.BusinessTypeDisplayNames[businessType] // Display name (e.g., "Private Limited Company (Ltd)")
		description := bizconstants.BusinessTypeDescriptions[businessType]
		icon := bizconstants.BusinessTypeIcons[businessType]
		
		// Insert or update - let PostgreSQL generate the UUID
		query := `
			INSERT INTO business_types (id, name, display_name, description, icon, is_active, created_at, updated_at)
			VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, NOW(), NOW())
			ON CONFLICT (name) DO UPDATE SET
				display_name = EXCLUDED.display_name,
				description = EXCLUDED.description,
				icon = EXCLUDED.icon,
				is_active = EXCLUDED.is_active,
				updated_at = NOW()
			WHERE business_types.display_name != EXCLUDED.display_name 
			   OR business_types.description != EXCLUDED.description
			   OR business_types.icon != EXCLUDED.icon
		`
		
		err := db.Exec(query,
			name,        // ✅ "private_company"
			displayName, // ✅ "Private Limited Company (Ltd)"
			description,
			icon,
			true,
		).Error
		
		if err != nil {
			return err
		}
	}

	log.Printf("✅ Business types seeded: %d types", len(bizconstants.AllBusinessTypes))
	return nil
}