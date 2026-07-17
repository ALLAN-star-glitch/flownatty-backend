// internal/database/seeders/admin_seeder.go
package seeders

import (
	"log"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	adminID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	admin := &models.User{
		BaseModel: models.BaseModel{
			ID: adminID,
		},
		PhoneNumber:      "+254700000001",
		Email:            "allanmathenge22@gmail.com",
		Password:         string(hashedPassword),
		Name:             "Super Admin",
		Role:             "super_admin",
		IsVerified:       true,
		IsEmailVerified:  true,
		VerifiedAt:       &now,
		EmailVerifiedAt:  &now,
		LastActiveAt:     &now,
		TwoFactorEnabled: true,
	}

	// Use FirstOrCreate - find by email, create if not exists
	if err := db.Where(models.User{Email: admin.Email}).FirstOrCreate(admin).Error; err != nil {
		return err
	}

	// Assign admin role in Casbin (idempotent - only if not exists)
	var count int64
	db.Table("casbin_rule").
		Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", 
			"g", admin.ID.String(), "super_admin", "platform").
		Count(&count)
	
	if count == 0 {
		if err := db.Exec(`
			INSERT INTO casbin_rule (ptype, v0, v1, v2) 
			VALUES (?, ?, ?, ?)`,
			"g", admin.ID.String(), "super_admin", "platform",
		).Error; err != nil {
			return err
		}
		log.Printf("✅ Admin user created: %s (ID: %s)", admin.Email, admin.ID)
	} else {
		log.Printf("✅ Admin user already exists: %s", admin.Email)
	}

	return nil
}