package repository

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAllUsers gets all users with pagination
func (r *UserRepository) GetAllUsers(limit, offset int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := r.db.Model(&models.User{})

	// Apply search filter if provided
	if search != "" {
		db = db.Where("email ILIKE ? OR name ILIKE ? OR phone_number ILIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := db.Preload("BusinessMembers").Preload("BusinessMembers.Business").Preload("BusinessMembers.User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// GetUserStats gets user statistics
func (r *UserRepository) GetUserStats() (map[string]interface{}, error) {
	var totalUsers int64
	var totalConsumers int64
	var totalBusinessOwners int64

	// Total users
	if err := r.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	// Total consumers
	if err := r.db.Model(&models.User{}).Where("role = ?", "consumer").Count(&totalConsumers).Error; err != nil {
		return nil, err
	}

	// Total business owners
	if err := r.db.Model(&models.User{}).Where("role = ?", "business_owner").Count(&totalBusinessOwners).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_users":           totalUsers,
		"total_consumers":       totalConsumers,
		"total_business_owners": totalBusinessOwners,
	}

	return stats, nil
}