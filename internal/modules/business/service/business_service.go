package service

import (
    "errors"
    "fmt"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/repository"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type BusinessService struct {
    repo        *repository.BusinessRepository
    productRepo *repository.ProductRepository
    memberRepo  *repository.BusinessMemberRepository
    db          *gorm.DB
}

func NewBusinessService(
    repo *repository.BusinessRepository,
    productRepo *repository.ProductRepository,
    memberRepo *repository.BusinessMemberRepository,
    db *gorm.DB,
) *BusinessService {
    return &BusinessService{
        repo:        repo,
        productRepo: productRepo,
        memberRepo:  memberRepo,
        db:          db,
    }
}

// ================================================
// BUSINESS OPERATIONS
// ================================================

// GetBusinessByID gets a business by ID
func (s *BusinessService) GetBusinessByID(id uuid.UUID) (*models.Business, error) {
    return s.repo.GetBusinessByID(id)
}

// GetBusinessWithProducts gets a business with its products
func (s *BusinessService) GetBusinessWithProducts(id uuid.UUID) (*models.Business, error) {
    return s.repo.GetBusinessByIDWithProducts(id)
}

// UpdateBusiness updates a business
func (s *BusinessService) UpdateBusiness(id uuid.UUID, updates map[string]interface{}) (*models.Business, error) {
    business, err := s.repo.GetBusinessByID(id)
    if err != nil {
        return nil, err
    }
    if business == nil {
        return nil, errors.New("business not found")
    }

    // Allowed fields to update
    allowedFields := map[string]bool{
        "name": true, "category": true, "description": true,
        "logo": true, "phone": true, "email": true,
        "address": true, "location": true,
        "latitude": true, "longitude": true,
        "is_active": true,
    }

    for key, value := range updates {
        if allowedFields[key] {
            switch key {
            case "name":
                business.Name = value.(string)
            case "category":
                business.Category = value.(string)
            case "description":
                business.Description = value.(string)
            case "logo":
                business.Logo = value.(string)
            case "phone":
                business.Phone = value.(string)
            case "email":
                business.Email = value.(string)
            case "address":
                business.Address = value.(string)
            case "location":
                business.Location = value.(string)
            case "latitude":
                business.Latitude = value.(float64)
            case "longitude":
                business.Longitude = value.(float64)
            case "is_active":
                business.IsActive = value.(bool)
            }
        }
    }

    if err := s.repo.UpdateBusiness(business); err != nil {
        return nil, fmt.Errorf("failed to update business: %w", err)
    }

    return business, nil
}

// DeleteBusiness deletes a business
func (s *BusinessService) DeleteBusiness(id uuid.UUID) error {
    return s.repo.DeleteBusiness(id)
}

// SearchBusinesses searches for businesses
func (s *BusinessService) SearchBusinesses(query, category string, page, pageSize int) ([]models.Business, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    offset := (page - 1) * pageSize
    return s.repo.SearchBusinesses(query, category, pageSize, offset)
}

// GetBusinessesByCategory gets businesses by category
func (s *BusinessService) GetBusinessesByCategory(category string) ([]models.Business, error) {
    return s.repo.GetBusinessesByCategory(category)
}

// GetBusinessesWithProducts gets businesses that have products
func (s *BusinessService) GetBusinessesWithProducts(page, pageSize int) ([]models.Business, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    offset := (page - 1) * pageSize
    return s.repo.GetBusinessesWithProducts(pageSize, offset)
}

// GetAllBusinesses gets all active businesses (with pagination)
func (s *BusinessService) GetAllBusinesses(page, pageSize int) ([]models.Business, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    offset := (page - 1) * pageSize
    return s.repo.GetAllBusinesses(pageSize, offset)
}

// GetBusinessStats gets business statistics
func (s *BusinessService) GetBusinessStats(businessID uuid.UUID) (map[string]interface{}, error) {
    return s.repo.GetBusinessStats(businessID)
}

// VerifyBusiness verifies a business (admin only)
func (s *BusinessService) VerifyBusiness(id uuid.UUID) error {
    business, err := s.repo.GetBusinessByID(id)
    if err != nil {
        return err
    }
    if business == nil {
        return errors.New("business not found")
    }

    return s.repo.UpdateBusinessVerification(id, true)
}

// ================================================
// BUSINESS MEMBER OPERATIONS
// ================================================

// AddBusinessMember adds a user to a business
func (s *BusinessService) AddBusinessMember(businessID, userID uuid.UUID, role string) (*models.BusinessMember, error) {
    // Check if user is already a member
    existing, err := s.memberRepo.GetByUserAndBusiness(userID, businessID)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, errors.New("user is already a member of this business")
    }

    member := &models.BusinessMember{
        BusinessID: businessID,
        UserID:     userID,
        Role:       role,
        IsActive:   true,
    }

    if err := s.memberRepo.Create(member); err != nil {
        return nil, fmt.Errorf("failed to add business member: %w", err)
    }

    return member, nil
}

// GetBusinessMembers gets all members of a business
func (s *BusinessService) GetBusinessMembers(businessID uuid.UUID) ([]models.BusinessMember, error) {
    return s.memberRepo.GetMembersByBusiness(businessID)
}

// GetBusinessesByUserID gets all businesses a user belongs to
func (s *BusinessService) GetBusinessesByUserID(userID uuid.UUID) ([]models.Business, error) {
    members, err := s.memberRepo.GetMembersByUser(userID)
    if err != nil {
        return nil, err
    }

    var businesses []models.Business
    for _, member := range members {
        if member.IsActive {
            businesses = append(businesses, member.Business)
        }
    }
    return businesses, nil
}

// GetBusinessByUserID gets the first business a user belongs to
func (s *BusinessService) GetBusinessByUserID(userID uuid.UUID) (*models.Business, error) {
    members, err := s.memberRepo.GetMembersByUser(userID)
    if err != nil {
        return nil, err
    }

    for _, member := range members {
        if member.IsActive {
            return &member.Business, nil
        }
    }
    return nil, nil
}

// UpdateBusinessMemberRole updates a member's role
func (s *BusinessService) UpdateBusinessMemberRole(businessID, userID uuid.UUID, role string) error {
    member, err := s.memberRepo.GetByUserAndBusiness(userID, businessID)
    if err != nil {
        return err
    }
    if member == nil {
        return errors.New("user is not a member of this business")
    }

    member.Role = role
    return s.memberRepo.Update(member)
}

// RemoveBusinessMember removes a user from a business
func (s *BusinessService) RemoveBusinessMember(businessID, userID uuid.UUID) error {
    return s.memberRepo.Delete(userID, businessID)
}

// IsUserMemberOfBusiness checks if a user is a member of a business
func (s *BusinessService) IsUserMemberOfBusiness(userID, businessID uuid.UUID) (bool, error) {
    member, err := s.memberRepo.GetByUserAndBusiness(userID, businessID)
    if err != nil {
        return false, err
    }
    return member != nil && member.IsActive, nil
}

// GetUserRoleInBusiness gets a user's role in a business
func (s *BusinessService) GetUserRoleInBusiness(userID, businessID uuid.UUID) (string, error) {
    member, err := s.memberRepo.GetByUserAndBusiness(userID, businessID)
    if err != nil {
        return "", err
    }
    if member == nil {
        return "", errors.New("user is not a member of this business")
    }
    return member.Role, nil
}

// ================================================
// BUSINESS OWNER OPERATIONS
// ================================================

// GetBusinessOwners gets all owners of a business
func (s *BusinessService) GetBusinessOwners(businessID uuid.UUID) ([]models.BusinessMember, error) {
    return s.memberRepo.GetBusinessOwners(businessID)
}

// IsUserBusinessOwner checks if a user is an owner of a business
func (s *BusinessService) IsUserBusinessOwner(userID, businessID uuid.UUID) (bool, error) {
    member, err := s.memberRepo.GetByUserAndBusiness(userID, businessID)
    if err != nil {
        return false, err
    }
    return member != nil && member.IsActive && member.Role == "business_owner", nil
}