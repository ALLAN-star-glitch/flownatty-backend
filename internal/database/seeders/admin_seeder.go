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
	// Check if admin already exists
	var count int64
	db.Model(&models.User{}).Where("email = ?", "allanmathenge22@gmail.com").Count(&count)
	if count > 0 {
		log.Println("Admin user already exists")
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create admin user with known UUID
	adminID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	now := time.Now()

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

	if err := db.Create(admin).Error; err != nil {
		return err
	}

	// Assign admin role in Casbin
	if err := db.Exec(`
		INSERT INTO casbin_rule (ptype, v0, v1, v2) 
		VALUES (?, ?, ?, ?)`,
		"g", admin.ID.String(), "super_admin", "platform",
	).Error; err != nil {
		return err
	}

	log.Printf("Admin user created: %s (ID: %s)", admin.Email, admin.ID)
	return nil
}