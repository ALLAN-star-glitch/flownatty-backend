package bizservice

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizrepository"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/validators/bizvalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessService struct {
	repo           *bizrepository.BusinessRepository
	productRepo    *bizrepository.ProductRepository
	onboardingRepo *bizrepository.OnboardingRepository
	memberRepo     *bizrepository.BusinessMemberRepository
	enforcer       *permissions.Enforcer
	permService    *permissions.Service
	validator      *bizvalidator.BusinessValidator
	db             *gorm.DB
}

func NewBusinessService(
	repo *bizrepository.BusinessRepository,
	productRepo *bizrepository.ProductRepository,
	onboardingRepo *bizrepository.OnboardingRepository,
	memberRepo *bizrepository.BusinessMemberRepository,
	enforcer *permissions.Enforcer,
	permService *permissions.Service,
	db *gorm.DB,
) *BusinessService {
	return &BusinessService{
		repo:           repo,
		productRepo:    productRepo,
		onboardingRepo: onboardingRepo,
		memberRepo:     memberRepo,
		enforcer:       enforcer,
		permService:    permService,
		validator:      bizvalidator.NewBusinessValidator(),
		db:             db,
	}
}

// ================================================
// BUSINESS OPERATIONS
// ================================================

// GetBusinessByID gets a business by ID (public - no permission check)
func (s *BusinessService) GetBusinessByID(id uuid.UUID) (*models.Business, error) {
	return s.repo.GetBusinessByID(id)
}

// UpdateBusiness updates a business
func (s *BusinessService) UpdateBusiness(ctx context.Context, userID, businessID uuid.UUID, updates map[string]interface{}) (*models.Business, error) {
	//  Use enforcer directly
	allowed, err := s.enforcer.Enforce(
		userID.String(),
		permissions.BusinessDomain(businessID.String()),
		permissions.ResourceBusiness,
		permissions.ActionUpdate,
	)
	if err != nil {
		return nil, fmt.Errorf("permission check failed: %w", err)
	}
	if !allowed {
		return nil, errors.New("insufficient permissions to update business")
	}

	business, err := s.repo.GetBusinessByID(businessID)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, errors.New("business not found")
	}

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
func (s *BusinessService) DeleteBusiness(ctx context.Context, userID, businessID uuid.UUID) error {
	//  Use enforcer directly
	allowed, err := s.enforcer.Enforce(
		userID.String(),
		permissions.BusinessDomain(businessID.String()),
		permissions.ResourceBusiness,
		permissions.ActionDelete,
	)
	if err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}
	if !allowed {
		return errors.New("insufficient permissions to delete business")
	}

	return s.repo.DeleteBusiness(businessID)
}

// GetBusinessStats gets business statistics
func (s *BusinessService) GetBusinessStats(ctx context.Context, userID, businessID uuid.UUID) (map[string]interface{}, error) {
	//  Use enforcer directly
	isAdmin := s.enforcer.IsBusinessAdmin(userID.String(), businessID.String())
	isManager := s.enforcer.IsProductManager(userID.String(), businessID.String()) ||
		s.enforcer.IsOrderManager(userID.String(), businessID.String()) ||
		s.enforcer.IsContentManager(userID.String(), businessID.String()) ||
		s.enforcer.IsServiceManager(userID.String(), businessID.String())

	if !isAdmin && !isManager {
		return nil, errors.New("insufficient permissions to view business stats")
	}

	return s.repo.GetBusinessStats(businessID)
}

// VerifyBusiness verifies a business (platform admin only)
func (s *BusinessService) VerifyBusiness(ctx context.Context, userID, businessID uuid.UUID) error {
	//  Use enforcer directly for platform admin
	if !s.enforcer.IsAdmin(userID.String()) {
		return errors.New("only platform admins can verify businesses")
	}

	business, err := s.repo.GetBusinessByID(businessID)
	if err != nil {
		return err
	}
	if business == nil {
		return errors.New("business not found")
	}

	return s.repo.UpdateBusinessVerification(businessID, true)
}

// AddBusinessMember adds a user to a business
func (s *BusinessService) AddBusinessMember(ctx context.Context, userID, businessID, targetUserID uuid.UUID, role string) (*models.BusinessMember, error) {
	//  Use enforcer directly
	if !s.enforcer.IsBusinessAdmin(userID.String(), businessID.String()) {
		return nil, errors.New("only business admins can add members")
	}

	existing, err := s.memberRepo.GetByUserAndBusiness(targetUserID, businessID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user is already a member of this business")
	}

	member := &models.BusinessMember{
		BusinessID: businessID,
		UserID:     targetUserID,
		Role:       role,
		IsActive:   true,
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, fmt.Errorf("failed to add business member: %w", err)
	}

	if _, err := s.enforcer.AddBusinessRole(targetUserID.String(), businessID.String(), permissions.Role(role)); err != nil {
		fmt.Printf("Failed to assign Casbin role: %v", err)
	}

	return member, nil
}

// GetBusinessMembers gets all members of a business
func (s *BusinessService) GetBusinessMembers(ctx context.Context, userID, businessID uuid.UUID) ([]models.BusinessMember, error) {
	//  Use enforcer directly
	roles := s.enforcer.GetRolesForUserInDomain(userID.String(), permissions.BusinessDomain(businessID.String()))
	if len(roles) == 0 {
		return nil, errors.New("insufficient permissions to view members")
	}

	return s.memberRepo.GetMembersByBusiness(businessID)
}

// UpdateBusinessMemberRole updates a member's role
func (s *BusinessService) UpdateBusinessMemberRole(ctx context.Context, userID, businessID, targetUserID uuid.UUID, role string) error {
	//  Use enforcer directly
	if !s.enforcer.IsBusinessAdmin(userID.String(), businessID.String()) {
		return errors.New("only business admins can update roles")
	}

	member, err := s.memberRepo.GetByUserAndBusiness(targetUserID, businessID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("user is not a member of this business")
	}

	if _, err := s.enforcer.RemoveBusinessRole(targetUserID.String(), businessID.String(), permissions.Role(member.Role)); err != nil {
		fmt.Printf("Failed to remove old Casbin role: %v", err)
	}

	member.Role = role
	if err := s.memberRepo.Update(member); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	if _, err := s.enforcer.AddBusinessRole(targetUserID.String(), businessID.String(), permissions.Role(role)); err != nil {
		fmt.Printf("Failed to assign new Casbin role: %v", err)
	}

	return nil
}

// RemoveBusinessMember removes a user from a business
func (s *BusinessService) RemoveBusinessMember(ctx context.Context, userID, businessID, targetUserID uuid.UUID) error {
	//  Use enforcer directly
	if !s.enforcer.IsBusinessAdmin(userID.String(), businessID.String()) {
		return errors.New("only business admins can remove members")
	}

	if userID == targetUserID {
		return errors.New("cannot remove yourself")
	}

	target, err := s.memberRepo.GetByUserAndBusiness(targetUserID, businessID)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("target user is not a member of this business")
	}

	if _, err := s.enforcer.RemoveBusinessRole(targetUserID.String(), businessID.String(), permissions.Role(target.Role)); err != nil {
		fmt.Printf("Failed to remove Casbin role: %v", err)
	}

	return s.memberRepo.Delete(targetUserID, businessID)
}

// ================================================
// PUBLIC METHODS (No permission checks)
// ================================================

// SearchBusinesses searches for businesses (public)
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

// GetBusinessesByCategory gets businesses by category (public)
func (s *BusinessService) GetBusinessesByCategory(category string) ([]models.Business, error) {
	return s.repo.GetBusinessesByCategory(category)
}

// GetBusinessesWithProducts gets businesses that have products (public)
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

// GetAllBusinesses gets all active businesses (public)
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

// GetBusinessesByUserID gets all businesses a user belongs to (self)
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

// GetBusinessByUserID gets the first business a user belongs to (self)
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


// GetBusinessOwners gets all owners of a business (display only)
func (s *BusinessService) GetBusinessOwners(businessID uuid.UUID) ([]models.BusinessMember, error) {
	return s.memberRepo.GetBusinessOwners(businessID)
}

// IsUserBusinessOwner checks if a user is an owner (display only)
func (s *BusinessService) IsUserBusinessOwner(ctx context.Context, userID, businessID uuid.UUID) (bool, error) {
	member, err := s.memberRepo.GetByUserAndBusiness(userID, businessID)
	if err != nil {
		return false, err
	}
	return member != nil && member.IsActive && member.Role == permissions.RoleBusinessAdmin.String(), nil
}

// ================================================
// REFERENCE DATA METHODS
// ================================================

// GetAllBusinessTypes gets all business types
func (s *BusinessService) GetAllBusinessTypes() ([]models.BusinessType, error) {
	return s.repo.GetAllBusinessTypes()
}

// GetBusinessTypeByID gets a business type by ID
func (s *BusinessService) GetBusinessTypeByID(id uuid.UUID) (*models.BusinessType, error) {
	return s.repo.GetBusinessTypeByID(id)
}

// GetAllBusinessSectors gets all business sectors
func (s *BusinessService) GetAllBusinessSectors() ([]models.BusinessSector, error) {
	return s.repo.GetAllBusinessSectors()
}

// GetAllBusinessSubcategories gets all business subcategories
func (s *BusinessService) GetAllBusinessSubcategories() ([]models.BusinessSubcategory, error) {
	return s.repo.GetAllBusinessSubcategories()
}

// GetBusinessSubcategoriesBySector gets subcategories by sector ID
func (s *BusinessService) GetBusinessSubcategoriesBySector(sectorID uuid.UUID) ([]models.BusinessSubcategory, error) {
	return s.repo.GetBusinessSubcategoriesBySector(sectorID)
}

// GetAllEstablishmentTypes gets all establishment types
func (s *BusinessService) GetAllEstablishmentTypes() ([]models.EstablishmentType, error) {
	return s.repo.GetAllEstablishmentTypes()
}

// GetBusinessTypeByName gets a business type by name
func (s *BusinessService) GetBusinessTypeByName(name string) (*models.BusinessType, error) {
	return s.repo.GetBusinessTypeByName(name)
}

// GetSectorByName gets a sector by name
func (s *BusinessService) GetSectorByName(name string) (*models.BusinessSector, error) {
	return s.repo.GetSectorByName(name)
}

// OnboardBusiness creates a new business during user registration

func (s *BusinessService) OnboardBusinessInit(
	ctx context.Context,
	adminID uuid.UUID,
	req *bizrepository.OnboardingRequest,
) (*models.Business, error) {
	
	// 1. Validate using existing validator
	if err := s.validator.ValidateBusinessType(req.BusinessType); err != nil {
		return nil, err
	}
	
	if err := s.validator.ValidateBusinessName(req.BusinessName); err != nil {
		return nil, err
	}
	
	if err := s.validator.ValidateBusinessPhone(req.BusinessPhone); err != nil {
		return nil, err
	}
	
	if err := s.validator.ValidateBusinessEmailRequired(req.BusinessEmail); err != nil {
		return nil, err
	}
	
	if req.BusinessAddress != "" {
		if err := s.validator.ValidateBusinessAddress(req.BusinessAddress); err != nil {
			return nil, err
		}
	}

	// 2. Normalize phone
	normalizedPhone := s.validator.NormalizeBusinessPhone(req.BusinessPhone)

	// 3. Create business with REAL IDs from the request
	business := &models.Business{
		Name:           req.BusinessName,
		Email:          req.BusinessEmail,
		Phone:          normalizedPhone,
		Address:        req.BusinessAddress,
		BusinessTypeID: req.BusinessTypeID, // REAL UUID from DB
		SectorID:       req.SectorID,       // REAL UUID from DB
		IsActive:       true,
		IsVerified:     false,
	}

	// 4. Create business with admin
	if err := s.onboardingRepo.CreateBusinessWithAdminInit(business, adminID); err != nil {
		return nil, fmt.Errorf("failed to create business: %w", err)
	}

	// 5. Assign business admin role
	if err := s.permService.AssignBusinessAdminRole(ctx, adminID.String(), business.ID.String()); err != nil {
		fmt.Printf("Warning: Failed to assign business admin role: %v\n", err)
	}

	log.Printf("Business onboarded: %s, TypeID: %v, SectorID: %v", 
		business.ID, business.BusinessTypeID, business.SectorID)

	return business, nil
}

// ================================================
// MEMBER METHODS (For Auth Module)
// ================================================

// GetMemberByUserAndBusiness gets a member by user and business ID
func (s *BusinessService) GetMemberByUserAndBusiness(ctx context.Context, userID, businessID uuid.UUID) (*models.BusinessMember, error) {
	return s.memberRepo.GetByUserAndBusiness(userID, businessID)
}

// GetMembersByUser gets all memberships for a user
func (s *BusinessService) GetMembersByUser(ctx context.Context, userID uuid.UUID) ([]models.BusinessMember, error) {
	return s.memberRepo.GetMembersByUser(userID)
}

// GetMembersByBusiness gets all members of a business
func (s *BusinessService) GetMembersByBusiness(ctx context.Context, businessID uuid.UUID) ([]models.BusinessMember, error) {
	return s.memberRepo.GetMembersByBusiness(businessID)
}

// GetMembersByBusinessWithRole gets all members of a business with a specific role
func (s *BusinessService) GetMembersByBusinessWithRole(ctx context.Context, businessID uuid.UUID, role string) ([]models.BusinessMember, error) {
	return s.memberRepo.GetMembersByBusinessWithRole(businessID, role)
}

// IsUserMemberOfBusiness checks if a user is a member of a business (using context)
func (s *BusinessService) IsUserMemberOfBusiness(ctx context.Context, userID, businessID uuid.UUID) (bool, error) {
	member, err := s.memberRepo.GetByUserAndBusiness(userID, businessID)
	if err != nil {
		return false, err
	}
	return member != nil && member.IsActive, nil
}

// IsUserBusinessAdmin checks if a user is a business admin
func (s *BusinessService) IsUserBusinessAdmin(ctx context.Context, userID, businessID uuid.UUID) (bool, error) {
	member, err := s.memberRepo.GetByUserAndBusiness(userID, businessID)
	if err != nil {
		return false, err
	}
	return member != nil && member.IsActive && member.Role == permissions.RoleBusinessAdmin.String(), nil
}

// GetUserRoleInBusiness gets a user's role in a business
func (s *BusinessService) GetUserRoleInBusiness(ctx context.Context, userID, businessID uuid.UUID) (string, error) {
	member, err := s.memberRepo.GetByUserAndBusiness(userID, businessID)
	if err != nil {
		return "", err
	}
	if member == nil {
		return "", errors.New("user is not a member of this business")
	}
	return member.Role, nil
}

func (s *BusinessService) GetBusinessesByType(businessType string, page, pageSize int) ([]models.Business, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    offset := (page - 1) * pageSize
    return s.repo.GetBusinessesByType(businessType, pageSize, offset)
}

// GetBusinessesBySector gets businesses by sector
func (s *BusinessService) GetBusinessesBySector(sector string, page, pageSize int) ([]models.Business, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    offset := (page - 1) * pageSize
    return s.repo.GetBusinessesBySector(sector, pageSize, offset)
}
