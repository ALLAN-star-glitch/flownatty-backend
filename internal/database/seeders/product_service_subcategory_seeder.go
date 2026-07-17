// internal/database/seeders/product_service_subcategory_seeder.go

package seeders

import (
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/constants/productconstants"
	"gorm.io/gorm"
)

func SeedProductServiceSubcategories(db *gorm.DB) error {
	// Use raw SQL with ON CONFLICT for PostgreSQL
	for _, info := range productconstants.ProductSubcategories {
		// Get the category name from the category info
		categoryInfo, exists := productconstants.ProductCategories[info.Category]
		if !exists {
			log.Printf("⚠️ Category not found for subcategory: %s", info.DisplayName)
			continue
		}

		// Get the category_id from the database
		var categoryID string
		checkQuery := `SELECT id FROM product_service_categories WHERE name = ? AND type = ?`
		err := db.Raw(checkQuery, categoryInfo.DisplayName, categoryInfo.Type).Scan(&categoryID).Error
		if err != nil || categoryID == "" {
			log.Printf("⚠️ Category not found in DB: %s", categoryInfo.DisplayName)
			continue
		}

		// Insert or update - let PostgreSQL generate the UUID
		query := `
			INSERT INTO product_service_subcategories (id, category_id, name, description, icon, is_active, created_at, updated_at)
			VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, NOW(), NOW())
			ON CONFLICT (category_id, name) DO UPDATE SET
				description = EXCLUDED.description,
				icon = EXCLUDED.icon,
				is_active = EXCLUDED.is_active,
				updated_at = NOW()
			WHERE product_service_subcategories.description != EXCLUDED.description 
			   OR product_service_subcategories.icon != EXCLUDED.icon
		`

		err = db.Exec(query,
			categoryID,
			info.DisplayName,
			info.Description,
			info.Icon,
			true,
		).Error

		if err != nil {
			return err
		}
	}

	log.Printf("✅ Product/Service subcategories seeded: %d subcategories", len(productconstants.ProductSubcategories))
	return nil
}