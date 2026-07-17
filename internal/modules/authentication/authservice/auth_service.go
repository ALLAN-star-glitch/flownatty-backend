// internal/modules/authentication/service/auth_service.go

package authservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authrepo"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/background"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizrepository"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/queue"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo            *authrepo.AuthRepository
	config          *config.Config
	queue           *queue.Client
	permService     *permissions.Service
	businessService *bizservice.BusinessService
	tokenService    *TokenService
}

func NewAuthService(
	repo *authrepo.AuthRepository,
	cfg *config.Config,
	queueClient *queue.Client,
	permService *permissions.Service,
	businessService *bizservice.BusinessService,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		repo:            repo,
		config:          cfg,
		queue:           queueClient,
		permService:     permService,
		businessService: businessService,
		tokenService:    tokenService,
	}
}

// ================================================
// OTP GENERATION & MANAGEMENT
// ================================================

// GenerateOTP - ALWAYS random
func (s *AuthService) GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// StoreOTP in Redis with TTL (for registration)
func (s *AuthService) StoreOTP(email, otp string) error {
	key := fmt.Sprintf("otp:%s", email)
	if err := redis.Set(key, otp, 5*time.Minute); err != nil {
		return err
	}
	return nil
}

// GetOTP from Redis (for registration)
func (s *AuthService) GetOTP(email string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)
	otp, exists, err := redis.Get(key)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("OTP not found")
	}
	return otp, nil
}

// DeleteOTP from Redis (for registration)
func (s *AuthService) DeleteOTP(email string) error {
	key := fmt.Sprintf("otp:%s", email)
	return redis.Delete(key)
}

// StoreTwoFactorOTP stores 2FA OTP in Redis
func (s *AuthService) StoreTwoFactorOTP(email, otp string) error {
	key := fmt.Sprintf("2fa:%s", email)
	return redis.Set(key, otp, 5*time.Minute)
}

// GetTwoFactorOTP retrieves 2FA OTP from Redis
func (s *AuthService) GetTwoFactorOTP(email string) (string, error) {
	key := fmt.Sprintf("2fa:%s", email)
	otp, exists, err := redis.Get(key)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("2FA OTP not found or expired")
	}
	return otp, nil
}

// DeleteTwoFactorOTP removes 2FA OTP from Redis
func (s *AuthService) DeleteTwoFactorOTP(email string) error {
	key := fmt.Sprintf("2fa:%s", email)
	return redis.Delete(key)
}

// StoreBusinessOTP stores OTP for business email verification
func (s *AuthService) StoreBusinessOTP(email, otp string) error {
	key := fmt.Sprintf("business_otp:%s", email)
	return redis.Set(key, otp, 10*time.Minute)
}

// GetBusinessOTP retrieves business OTP from Redis
func (s *AuthService) GetBusinessOTP(email string) (string, error) {
	key := fmt.Sprintf("business_otp:%s", email)
	otp, exists, err := redis.Get(key)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("business OTP not found or expired")
	}
	return otp, nil
}

// DeleteBusinessOTP deletes business OTP from Redis
func (s *AuthService) DeleteBusinessOTP(email string) error {
	key := fmt.Sprintf("business_otp:%s", email)
	return redis.Delete(key)
}

// ================================================
// USER REGISTRATION DATA MANAGEMENT
// ================================================

// StoreUserData stores user registration data in Redis
func (s *AuthService) StoreUserData(email string, data map[string]interface{}) error {
	key := fmt.Sprintf("user:data:%s", email)
	if err := redis.HSet(key, data); err != nil {
		return err
	}
	return redis.Expire(key, 5*time.Minute)
}

// GetUserData gets user registration data from Redis
func (s *AuthService) GetUserData(email string) (map[string]string, error) {
	key := fmt.Sprintf("user:data:%s", email)
	result, exists, err := redis.HGetAll(key)
	if err != nil {
		return nil, err
	}
	if !exists || len(result) == 0 {
		return nil, errors.New("user data not found")
	}
	return result, nil
}

// DeleteUserData deletes user registration data from Redis
func (s *AuthService) DeleteUserData(email string) error {
	key := fmt.Sprintf("user:data:%s", email)
	return redis.Delete(key)
}

// ================================================
// BUSINESS DATA MANAGEMENT
// ================================================

// StoreBusinessData stores business registration data in Redis
func (s *AuthService) StoreBusinessData(userID string, data map[string]interface{}) error {
	key := fmt.Sprintf("business:data:%s", userID)
	if err := redis.HSet(key, data); err != nil {
		return err
	}
	return redis.Expire(key, 24*time.Hour)
}

// GetBusinessData retrieves business registration data from Redis
func (s *AuthService) GetBusinessData(userID string) (map[string]string, error) {
	key := fmt.Sprintf("business:data:%s", userID)
	result, exists, err := redis.HGetAll(key)
	if err != nil {
		return nil, err
	}
	if !exists || len(result) == 0 {
		return nil, errors.New("business data not found")
	}
	return result, nil
}

// DeleteBusinessData deletes business registration data from Redis
func (s *AuthService) DeleteBusinessData(userID string) error {
	key := fmt.Sprintf("business:data:%s", userID)
	return redis.Delete(key)
}

// MarkBusinessEmailVerified marks business email as verified
func (s *AuthService) MarkBusinessEmailVerified(userID string) error {
	key := fmt.Sprintf("business:data:%s", userID)
	return redis.HSetField(key, "business_email_verified", "true")
}

// IsBusinessEmailVerified checks if business email is verified
func (s *AuthService) IsBusinessEmailVerified(userID string) (bool, error) {
	key := fmt.Sprintf("business:data:%s", userID)
	result, exists, err := redis.HGet(key, "business_email_verified")
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return result == "true", nil
}

// ================================================
// USER CREATION & BUSINESS ONBOARDING (ALL BUSINESS LOGIC HERE)
// ================================================

// RegisterUser handles the complete user registration flow
func (s *AuthService) RegisterUser(ctx context.Context, userData map[string]string) error {
	// Check if user exists
	existingUser, _ := s.repo.GetUserByEmail(userData["email"])
	if existingUser != nil {
		return errors.New("email already registered")
	}

	existingPhone, _ := s.repo.GetUserByPhone(userData["phone"])
	if existingPhone != nil {
		return errors.New("phone number already registered")
	}

	// Generate and store OTP
	otp := s.GenerateOTP()
	if err := s.StoreOTP(userData["email"], otp); err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	// Convert map[string]string to map[string]interface{} for Redis storage
	userDataInterface := make(map[string]interface{})
	for key, value := range userData {
		userDataInterface[key] = value
	}

	// Store user data
	if err := s.StoreUserData(userData["email"], userDataInterface); err != nil {
		return fmt.Errorf("failed to store user data: %w", err)
	}

	// Send OTP email
	if err := s.EnqueueOTPEmail(userData["email"], userData["name"], otp); err != nil {
		log.Printf("Failed to enqueue OTP email: %v", err)
	}

	return nil
}

// VerifyOTPAndCreateUser handles OTP verification and user creation
func (s *AuthService) VerifyOTPAndCreateUser(ctx context.Context, email, otp string) (*models.User, *models.Business, map[string]interface{}, error) {
	// Verify OTP
	storedOTP, err := s.GetOTP(email)
	if err != nil {
		return nil, nil, nil, errors.New("invalid or expired OTP")
	}

	if otp != storedOTP {
		return nil, nil, nil, errors.New("invalid OTP")
	}

	// Get user data - returns map[string]string from Redis
	userDataMap, err := s.GetUserData(email)
	if err != nil {
		return nil, nil, nil, errors.New("registration data not found")
	}

	// Clean up OTP and user data
	s.DeleteOTP(email)
	s.DeleteUserData(email)

	// Check if this is a business user
	isBusiness := userDataMap["role"] == permissions.RoleBusinessAdmin.String()
	businessType := userDataMap["business_type"]

	// For non-sole_proprietor businesses, create user and store business data for later verification
	if isBusiness && businessType != "sole_proprietor" {
		user, businessData, err := s.CreateUserAndStoreBusinessData(ctx, userDataMap)
		if err != nil {
			return nil, nil, nil, err
		}
		// Generate tokens
		accessToken, refreshToken, err := s.GenerateTokens(ctx, user)
		if err != nil {
			return nil, nil, nil, err
		}
		
		// Send business verification email
		businessOTP := s.GenerateOTP()
		s.StoreBusinessOTP(userDataMap["business_email"], businessOTP)
		s.EnqueueBusinessOTPEmail(userDataMap["business_email"], userDataMap["business_name"], businessOTP)

		result := map[string]interface{}{
			"access_token":                   accessToken,
			"refresh_token":                  refreshToken,
			"business_verification_required": true,
			"business_email":                 userDataMap["business_email"],
			"business_name":                  userDataMap["business_name"],
			"next_step":                      "verify_business_email",
			"user":                           user,
			"business_data":                  businessData,
		}
		return user, nil, result, nil
	}

	// For consumers and sole proprietors, create user and business immediately
	return s.CreateUserWithBusiness(ctx, userDataMap)
}

// CreateUserOnly creates a user without creating a business
func (s *AuthService) CreateUserOnly(ctx context.Context, userData map[string]string) (*models.User, error) {
	now := time.Now()
	user := &models.User{
		PhoneNumber:      userData["phone"],
		Email:            userData["email"],
		Password:         userData["password"],
		Name:             userData["name"],
		Role:             userData["role"],
		IsVerified:       true,
		IsEmailVerified:  true,
		VerifiedAt:       &now,
		EmailVerifiedAt:  &now,
		LastActiveAt:     &now,
		TwoFactorEnabled: true,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign consumer role (all users get this)
	if err := s.permService.AssignConsumerRole(ctx, user.ID.String()); err != nil {
		log.Printf("Failed to assign consumer role: %v", err)
	}

	return user, nil
}



// CreateUserWithBusiness creates a user and immediately creates their business
func (s *AuthService) CreateUserWithBusiness(ctx context.Context, userData map[string]string) (*models.User, *models.Business, map[string]interface{}, error) {
	// 1. Create user
	user, err := s.CreateUserOnly(ctx, userData)
	if err != nil {
		return nil, nil, nil, err
	}

	// 2. Look up REAL business type ID from database
	var businessTypeID *uuid.UUID
	if userData["business_type"] != "" {
		businessType, err := s.businessService.GetBusinessTypeByName(userData["business_type"])
		if err != nil {
			log.Printf("Error looking up business type '%s': %v", userData["business_type"], err)
		} else if businessType != nil {
			businessTypeID = &businessType.ID
			log.Printf("Found business type: %s -> %s", userData["business_type"], businessType.ID)
		} else {
			log.Printf("Business type not found: %s", userData["business_type"])
		}
	}

	// 3. Look up REAL sector ID from database
	var sectorID *uuid.UUID
	if userData["business_category"] != "" {
		sector, err := s.businessService.GetSectorByName(userData["business_category"])
		if err != nil {
			log.Printf("Error looking up sector '%s': %v", userData["business_category"], err)
		} else if sector != nil {
			sectorID = &sector.ID
			log.Printf("Found sector: %s -> %s", userData["business_category"], sector.ID)
		} else {
			log.Printf("Sector not found: %s", userData["business_category"])
		}
	}

	// 4. Prepare business data with REAL IDs
	businessData := &bizrepository.OnboardingRequest{
		BusinessType:    userData["business_type"],
		BusinessName:    userData["business_name"],
		BusinessPhone:   userData["business_phone"],
		BusinessEmail:   userData["business_email"],
		BusinessAddress: userData["business_address"],
		BusinessTypeID:  businessTypeID, // REAL ID from DB
		SectorID:        sectorID,       // REAL ID from DB
	}

	// 5. Create business using business service
	business, err := s.businessService.OnboardBusinessInit(ctx, user.ID, businessData)
	if err != nil {
		log.Printf("Failed to create business for user %s: %v", user.ID, err)
		return user, nil, nil, fmt.Errorf("user created but business creation failed: %w", err)
	}

	// 6. Generate tokens
	accessToken, refreshToken, err := s.GenerateTokens(ctx, user)
	if err != nil {
		return user, business, nil, err
	}

	// 7. Send welcome email
	if err := s.EnqueueWelcomeEmail(user.Email, user.Name); err != nil {
		log.Printf("Failed to enqueue welcome email: %v", err)
	}

	result := map[string]interface{}{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"onboarding_required": true,
		"onboarding_url":     "/api/v1/businesses/onboarding",
		"business_name":      userData["business_name"],
		"user":               user,
		"business":           business,
	}

	log.Printf("Business created with ID: %s, TypeID: %v, SectorID: %v", 
		business.ID, business.BusinessTypeID, business.SectorID)

	return user, business, result, nil
}

// CreateUserAndStoreBusinessData creates a user and stores business data for later verification
func (s *AuthService) CreateUserAndStoreBusinessData(ctx context.Context, userData map[string]string) (*models.User, map[string]interface{}, error) {
	// 1. Create user
	user, err := s.CreateUserOnly(ctx, userData)
	if err != nil {
		return nil, nil, err
	}

	// 2. Prepare business data for Redis
	businessData := map[string]interface{}{
		"user_id":            user.ID.String(),
		"business_type":      userData["business_type"],
		"business_name":      userData["business_name"],
		"business_category":  userData["business_category"],
		"business_phone":     userData["business_phone"],
		"business_email":     userData["business_email"],
		"business_address":   userData["business_address"],
		"personal_email":     userData["email"],
		"personal_name":      userData["name"],
	}

	// 3. Store business data in Redis
	if err := s.StoreBusinessData(user.ID.String(), businessData); err != nil {
		return nil, nil, fmt.Errorf("failed to store business data: %w", err)
	}

	return user, businessData, nil
}


// CompleteBusinessVerification completes business verification and creates the business
func (s *AuthService) CompleteBusinessVerification(ctx context.Context, userID uuid.UUID, businessEmail, otp string) (*models.Business, map[string]interface{}, error) {
	// 1. Verify OTP
	storedOTP, err := s.GetBusinessOTP(businessEmail)
	if err != nil {
		return nil, nil, errors.New("invalid or expired OTP")
	}

	if otp != storedOTP {
		return nil, nil, errors.New("invalid OTP")
	}

	// 2. Get business data from Redis
	businessData, err := s.GetBusinessData(userID.String())
	if err != nil {
		return nil, nil, errors.New("business data not found")
	}

	// 3. Mark business email as verified
	if err := s.MarkBusinessEmailVerified(userID.String()); err != nil {
		log.Printf("Failed to mark business email verified: %v", err)
	}

	// 4. Delete OTP
	if err := s.DeleteBusinessOTP(businessEmail); err != nil {
		log.Printf("Failed to delete business OTP: %v", err)
	}

	// 5. Look up REAL business type ID from database
	var businessTypeID *uuid.UUID
	if businessData["business_type"] != "" {
		businessType, err := s.businessService.GetBusinessTypeByName(businessData["business_type"])
		if err != nil {
			log.Printf("Error looking up business type: %v", err)
		} else if businessType != nil {
			businessTypeID = &businessType.ID
			log.Printf("Found business type: %s -> %s", businessData["business_type"], businessType.ID)
		}
	}

	// 6. Look up REAL sector ID from database
	var sectorID *uuid.UUID
	if businessData["business_category"] != "" {
		sector, err := s.businessService.GetSectorByName(businessData["business_category"])
		if err != nil {
			log.Printf("Error looking up sector: %v", err)
		} else if sector != nil {
			sectorID = &sector.ID
			log.Printf("Found sector: %s -> %s", businessData["business_category"], sector.ID)
		}
	}

	// 7. Prepare onboarding request with REAL IDs
	req := &bizrepository.OnboardingRequest{
		BusinessType:    businessData["business_type"],
		BusinessName:    businessData["business_name"],
		BusinessPhone:   businessData["business_phone"],
		BusinessEmail:   businessEmail,
		BusinessAddress: businessData["business_address"],
		BusinessTypeID:  businessTypeID, // REAL ID from DB
		SectorID:        sectorID,       // REAL ID from DB
	}

	// 8. Create business
	business, err := s.businessService.OnboardBusinessInit(ctx, userID, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create business: %w", err)
	}

	// 9. Delete business data from Redis
	if err := s.DeleteBusinessData(userID.String()); err != nil {
		log.Printf("Failed to delete business data: %v", err)
	}

	// 10. Send business welcome email
	if err := s.EnqueueBusinessWelcomeEmail(businessEmail, businessData["business_name"], businessData["personal_name"]); err != nil {
		log.Printf("Failed to enqueue business welcome email: %v", err)
	}

	result := map[string]interface{}{
		"business_name":  businessData["business_name"],
		"business_email": businessEmail,
		"status":         "verified",
	}

	log.Printf("Business verified with ID: %s, TypeID: %v, SectorID: %v", 
		business.ID, business.BusinessTypeID, business.SectorID)

	return business, result, nil
}

// ================================================
// TOKEN MANAGEMENT
// ================================================

// GenerateTokens generates access and refresh tokens for a user
func (s *AuthService) GenerateTokens(ctx context.Context, user *models.User) (string, string, error) {
	accessToken, err := s.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(
		user.ID,
		"", // User-Agent will be set by handler
		"", // IP will be set by handler
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// RefreshTokens refreshes access token using refresh token
func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken, userAgent, ip string) (string, string, error) {
	return s.tokenService.RotateRefreshToken(refreshToken, userAgent, ip)
}

// RevokeToken revokes a refresh token
func (s *AuthService) RevokeToken(ctx context.Context, refreshToken string) error {
	return s.tokenService.RevokeRefreshToken(refreshToken)
}

// ================================================
// AUTHENTICATION METHODS
// ================================================

// LoginUser authenticates a user and initiates 2FA
func (s *AuthService) LoginUser(ctx context.Context, email, password string) (*models.User, string, error) {
	user, err := s.AuthenticateUserByEmail(email, password)
	if err != nil {
		return nil, "", err
	}

	// Generate 2FA OTP
	otp := s.GenerateOTP()
	if err := s.StoreTwoFactorOTP(user.Email, otp); err != nil {
		return nil, "", fmt.Errorf("failed to store 2FA OTP: %w", err)
	}

	// Send 2FA OTP email
	if err := s.EnqueueTwoFactorOTPEmail(user.Email, user.Name, otp); err != nil {
		log.Printf("Failed to enqueue 2FA OTP email: %v", err)
	}

	return user, otp, nil
}

// VerifyTwoFactorAndLogin verifies 2FA OTP and completes login
func (s *AuthService) VerifyTwoFactorAndLogin(ctx context.Context, email, otp string) (*models.User, string, string, error) {
	// Verify OTP
	storedOTP, err := s.GetTwoFactorOTP(email)
	if err != nil {
		return nil, "", "", errors.New("invalid or expired OTP")
	}

	if otp != storedOTP {
		return nil, "", "", errors.New("invalid OTP")
	}

	// Get user
	user, err := s.repo.GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, "", "", errors.New("user not found")
	}

	// Delete OTP
	s.DeleteTwoFactorOTP(email)

	// Generate tokens
	accessToken, refreshToken, err := s.GenerateTokens(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	// Send login notification
	if err := s.EnqueueLoginNotificationEmail(user.Email, user.Name, "", ""); err != nil {
		log.Printf("Failed to enqueue login notification: %v", err)
	}

	return user, accessToken, refreshToken, nil
}

// AuthenticateUserByEmail authenticates a user by email
func (s *AuthService) AuthenticateUserByEmail(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := user.ComparePassword(password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repo.UpdateUser(user); err != nil {
		log.Printf("Failed to update last login: %v", err)
	}

	return user, nil
}

// AuthenticateUserByPhone authenticates a user by phone number
func (s *AuthService) AuthenticateUserByPhone(phone, password string) (*models.User, error) {
	user, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := user.ComparePassword(password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repo.UpdateUser(user); err != nil {
		log.Printf("Failed to update last login: %v", err)
	}

	return user, nil
}

// ================================================
// PASSWORD RESET METHODS
// ================================================

// InitiatePasswordReset initiates password reset flow
func (s *AuthService) InitiatePasswordReset(ctx context.Context, email, newPassword string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		// Don't reveal if user exists
		return nil
	}

	otp := s.GenerateOTP()
	if err := s.StoreResetData(email, otp, newPassword); err != nil {
		return fmt.Errorf("failed to store reset data: %w", err)
	}

	if err := s.EnqueuePasswordResetOTPEmail(email, user.Name, otp); err != nil {
		log.Printf("Failed to enqueue password reset OTP: %v", err)
	}

	return nil
}

// VerifyResetOTPAndResetPassword verifies OTP and resets password
func (s *AuthService) VerifyResetOTPAndResetPassword(ctx context.Context, email, otp string) error {
	data, err := s.GetResetData(email)
	if err != nil {
		return errors.New("invalid or expired OTP")
	}

	if len(data) == 0 {
		return errors.New("invalid or expired OTP")
	}

	storedOTP, ok := data["otp"]
	if !ok || otp != storedOTP {
		return errors.New("invalid OTP")
	}

	newPassword, ok := data["new_password"]
	if !ok {
		return errors.New("invalid reset data")
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	if err := s.repo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.DeleteResetData(email)

	if err := s.EnqueuePasswordResetConfirmEmail(email, user.Name); err != nil {
		log.Printf("Failed to enqueue password reset confirmation: %v", err)
	}

	return nil
}

// StoreResetData stores password reset data in Redis
func (s *AuthService) StoreResetData(email, otp, newPassword string) error {
	key := fmt.Sprintf("reset:%s", email)
	data := map[string]interface{}{
		"otp":          otp,
		"new_password": newPassword,
	}
	if err := redis.HSet(key, data); err != nil {
		return err
	}
	return redis.Expire(key, 5*time.Minute)
}

// GetResetData retrieves password reset data from Redis
func (s *AuthService) GetResetData(email string) (map[string]string, error) {
	key := fmt.Sprintf("reset:%s", email)
	result, exists, err := redis.HGetAll(key)
	if err != nil {
		return nil, err
	}
	if !exists || len(result) == 0 {
		return nil, errors.New("reset data not found")
	}
	return result, nil
}

// DeleteResetData deletes password reset data from Redis
func (s *AuthService) DeleteResetData(email string) error {
	key := fmt.Sprintf("reset:%s", email)
	return redis.Delete(key)
}

// ================================================
// EMAIL ENQUEUE METHODS
// ================================================

// EnqueueOTPEmail enqueues OTP email via Asynq
func (s *AuthService) EnqueueOTPEmail(to, name, otp string) error {
	task := background.OTPEmailTask{
		To:      to,
		Name:    name,
		OTP:     otp,
		Expires: "5 minutes",
	}
	payload, err := task.Payload()
	if err != nil {
		return err
	}
	return s.queue.Enqueue(background.TypeEmailOTP, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueueBusinessOTPEmail enqueues business OTP email via Asynq
func (s *AuthService) EnqueueBusinessOTPEmail(to, businessName, otp string) error {
	task := background.BusinessOTPEmailTask{
		To:           to,
		BusinessName: businessName,
		OTP:          otp,
		Expires:      "10 minutes",
	}
	payload, err := task.Payload()
	if err != nil {
		return err
	}
	return s.queue.Enqueue(background.TypeEmailBusinessOTP, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueueWelcomeEmail enqueues welcome email via Asynq
func (s *AuthService) EnqueueWelcomeEmail(to, name string) error {
	task := background.WelcomeEmailTask{
		To:   to,
		Name: name,
	}
	payload, err := task.Payload()
	if err != nil {
		return err
	}
	return s.queue.Enqueue(background.TypeEmailWelcome, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueueBusinessWelcomeEmail enqueues business welcome email via Asynq
func (s *AuthService) EnqueueBusinessWelcomeEmail(to, businessName, ownerName string) error {
	task := background.BusinessWelcomeEmailTask{
		To:           to,
		BusinessName: businessName,
		OwnerName:    ownerName,
	}
	payload, err := task.Payload()
	if err != nil {
		return err
	}
	return s.queue.Enqueue(background.TypeEmailBusinessWelcome, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueuePasswordResetOTPEmail enqueues password reset OTP email via Asynq
func (s *AuthService) EnqueuePasswordResetOTPEmail(to, name, otp string) error {
	task := background.PasswordResetOTPTask{
		To:      to,
		Name:    name,
		OTP:     otp,
		Expires: "5 minutes",
	}
	payload, err := task.Payload()
	if err != nil {
		return err
	}
	return s.queue.Enqueue(background.TypeEmailPasswordResetOTP, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueueLoginNotificationEmail enqueues login notification email via Asynq
func (s *AuthService) EnqueueLoginNotificationEmail(to, name, ipAddress, userAgent string) error {
	task := background.LoginNotificationTask{
		To:        to,
		Name:      name,
		Time:      time.Now().Format("2006-01-02 15:04:05 UTC"),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	payload, err := task.Payload()
	if err != nil {
		return err
	}
	return s.queue.Enqueue(background.TypeEmailLoginNotification, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueuePasswordResetConfirmEmail enqueues password reset confirmation email via Asynq
func (s *AuthService) EnqueuePasswordResetConfirmEmail(to, name string) error {
	task := background.PasswordResetConfirmTask{
		To:   to,
		Name: name,
	}
	payload, err := task.Payload()
	if err != nil {
		return err
	}
	return s.queue.Enqueue(background.TypeEmailPasswordResetConfirm, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueueTwoFactorOTPEmail enqueues 2FA OTP email via Asynq
func (s *AuthService) EnqueueTwoFactorOTPEmail(to, name, otp string) error {
	task := background.TwoFactorOTPTask{
		To:      to,
		Name:    name,
		OTP:     otp,
		Expires: "5 minutes",
	}
	payload, err := task.Payload()
	if err != nil {
		return err
	}
	return s.queue.Enqueue(background.TypeEmailTwoFactorOTP, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// ================================================
// USER INFO & ADMIN METHODS
// ================================================

// BuildUserInfo builds user info for response
func (s *AuthService) BuildUserInfo(user *models.User) UserInfo {
	info := UserInfo{
		ID:          user.ID.String(),
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Role:        user.Role,
	}

	if len(user.BusinessMembers) > 0 {
		for _, member := range user.BusinessMembers {
			if member.IsActive {
				businessID := member.BusinessID.String()
				info.BusinessID = &businessID
				break
			}
		}
	}

	return info
}

// GetUserByIDWithBusiness gets a user with their business memberships
func (s *AuthService) GetUserByIDWithBusiness(id uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByIDWithBusiness(id)
}

// GetAllUsers gets all users with pagination
func (s *AuthService) GetAllUsers(page, pageSize int, search string) ([]models.User, int64, error) {
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
func (s *AuthService) GetUserStats() (map[string]interface{}, error) {
	return s.repo.GetUserStats()
}

// ValidateToken validates JWT token
func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		return userID, nil
	}
	return "", errors.New("invalid token")
}

// ================================================
// RESPONSE MODELS
// ================================================

// UserInfo struct
type UserInfo struct {
	ID          string  `json:"id"`
	PhoneNumber string  `json:"phone_number"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	Avatar      string  `json:"avatar"`
	Role        string  `json:"role"`
	BusinessID  *string `json:"business_id,omitempty"`
}