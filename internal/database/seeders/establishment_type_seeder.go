// internal/database/seeders/establishment_type_seeder.go

package seeders

import (
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/establishmentconstants"
	"gorm.io/gorm"
)

func SeedEstablishmentTypes(db *gorm.DB) error {
	// Use raw SQL with ON CONFLICT for PostgreSQL
	sortOrder := 1
	for _, info := range establishmentconstants.EstablishmentTypes {
		// Get the icon from the icons map
		icon := establishmentconstants.EstablishmentTypeIcons[info.Name]
		
		// Insert or update - let PostgreSQL generate the UUID
		query := `
			INSERT INTO establishment_types (id, name, display_name, description, icon, category, sort_order, is_active, created_at, updated_at)
			VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
			ON CONFLICT (name) DO UPDATE SET
				display_name = EXCLUDED.display_name,
				description = EXCLUDED.description,
				icon = EXCLUDED.icon,
				category = EXCLUDED.category,
				sort_order = EXCLUDED.sort_order,
				is_active = EXCLUDED.is_active,
				updated_at = NOW()
			WHERE establishment_types.display_name != EXCLUDED.display_name 
			   OR establishment_types.description != EXCLUDED.description
			   OR establishment_types.icon != EXCLUDED.icon
			   OR establishment_types.category != EXCLUDED.category
		`
		
		err := db.Exec(query,
			info.Name,
			info.DisplayName,
			info.Description,
			icon,
			info.Category,
			sortOrder,
			true,
		).Error
		
		if err != nil {
			return err
		}
		sortOrder++
	}

	log.Printf("✅ Establishment types seeded: %d types", len(establishmentconstants.EstablishmentTypes))
	return nil
}