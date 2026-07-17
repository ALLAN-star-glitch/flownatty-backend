// internal/database/seeders/product_service_category_seeder.go

package seeders

import (
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/productconstants"
	"gorm.io/gorm"
)

func SeedProductServiceCategories(db *gorm.DB) error {
	// Use raw SQL with ON CONFLICT for PostgreSQL
	sortOrder := 1
	for _, info := range productconstants.ProductCategories {
		// Insert or update - let PostgreSQL generate the UUID
		query := `
			INSERT INTO product_service_categories (id, name, description, icon, type, sort_order, is_active, created_at, updated_at)
			VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, ?, NOW(), NOW())
			ON CONFLICT (name, type) DO UPDATE SET
				description = EXCLUDED.description,
				icon = EXCLUDED.icon,
				sort_order = EXCLUDED.sort_order,
				is_active = EXCLUDED.is_active,
				updated_at = NOW()
			WHERE product_service_categories.description != EXCLUDED.description 
			   OR product_service_categories.icon != EXCLUDED.icon
			   OR product_service_categories.sort_order != EXCLUDED.sort_order
		`
		
		err := db.Exec(query,
			info.DisplayName,
			info.Description,
			info.Icon,
			info.Type,
			sortOrder,
			true,
		).Error
		
		if err != nil {
			return err
		}
		sortOrder++
	}

	log.Printf("✅ Product/Service categories seeded: %d categories", len(productconstants.ProductCategories))
	return nil
}