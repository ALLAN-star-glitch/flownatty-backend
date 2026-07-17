package userservice

import (
	"errors"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/usermanagement/repository"
	"github.com/google/uuid"
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

	// Validate page .. to be at least 1, default to 1 if invalid
	if page < 1 {
		page = 1
	}

	// Validate pageSize to be between 1 and 100, default to 20 if invalid
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.GetAllUsers(pageSize, offset, search)
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {


	return s.repo.GetUserByID(id)
}

// GetUserByEmail gets a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

// GetUserStats gets user statistics
func (s *UserService) GetUserStats() (map[string]interface{}, error) {
	return s.repo.GetUserStats()
}

// CountUsersByRole counts users with a specific role
func (s *UserService) CountUsersByRole(role string) (int64, error) {
	return s.repo.CountUsersByRole(role)
}

// UpdateUserRole updates a user's role
func (s *UserService) UpdateUserRole(id uuid.UUID, newRole string) error {
	validRoles := map[string]bool{
		"consumer":        true,
		"business_owner":  true,
		"admin":           true,
		"super_admin":     true,
	}
	if !validRoles[newRole] {
		return errors.New("invalid role")
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Prevent demoting the last super_admin
	if newRole != "super_admin" && user.Role == "super_admin" {
		count, err := s.repo.CountUsersByRole("super_admin")
		if err != nil {
			return err
		}
		if count <= 1 {
			return errors.New("cannot demote the only super_admin")
		}
	}

	return s.repo.UpdateUserRole(id, newRole)
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uuid.UUID) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Prevent soft-deleting an already soft-deleted user
	if !user.DeletedAt.Time.IsZero() {
		return errors.New("user is already soft-deleted")
	}

	// Prevent deleting the last super_admin
	if user.Role == "super_admin" {
		count, err := s.repo.CountUsersByRole("super_admin")
		if err != nil {
			return err
		}
		if count <= 1 {
			return errors.New("cannot delete the only super_admin")
		}
	}

	return s.repo.DeleteUser(id)
}

// HardDeleteUser permanently deletes a user
func (s *UserService) HardDeleteUser(id uuid.UUID) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Prevent deleting the last super_admin
	if user.Role == "super_admin" {
		count, err := s.repo.CountUsersByRole("super_admin")
		if err != nil {
			return err
		}
		if count <= 1 {
			return errors.New("cannot delete the only super_admin")
		}
	}

	return s.repo.HardDeleteUser(id)
}

// SuspendUser suspends a user
func (s *UserService) SuspendUser(id uuid.UUID) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Cannot suspend a soft-deleted user
	if !user.DeletedAt.Time.IsZero() {
		return errors.New("cannot suspend a soft-deleted user")
	}

	// Cannot suspend a suspended user
	if !user.IsActive {
		return errors.New("user is already suspended")
	}

	// Prevent suspending the last super_admin
	if user.Role == "super_admin" {
		count, err := s.repo.CountUsersByRole("super_admin")
		if err != nil {
			return err
		}
		if count <= 1 {
			return errors.New("cannot suspend the only super_admin")
		}
	}

	return s.repo.SuspendUser(id)
}

// ActivateUser activates a user (handles both suspended AND soft-deleted users)
func (s *UserService) ActivateUser(id uuid.UUID) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Check if user is already active
	if user.IsActive && user.DeletedAt.Time.IsZero() {
		return errors.New("user is already active")
	}

	// ActivateUser will clear both is_active and deleted_at
	return s.repo.ActivateUser(id)
}