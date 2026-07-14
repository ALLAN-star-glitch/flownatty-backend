package service

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetAllUsers gets all users with pagination
func (s *UserService) GetAllUsers(page, pageSize int, search string) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.GetAllUsers(pageSize, offset, search)
}

// GetUserStats gets user statistics
func (s *UserService) GetUserStats() (map[string]interface{}, error) {
	return s.repo.GetUserStats()
}