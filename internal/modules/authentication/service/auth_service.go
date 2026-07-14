package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authrepo"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/background"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/queue"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *authrepo.AuthRepository
	config *config.Config
	queue  *queue.Client
}

func NewAuthService(repo *authrepo.AuthRepository, cfg *config.Config, queueClient *queue.Client) *AuthService {
	return &AuthService{
		repo:   repo,
		config: cfg,
		queue:  queueClient,
	}
}

// ================================================
// OTP GENERATION & MANAGEMENT
// ================================================

// GenerateOTP - ALWAYS random
func (s *AuthService) GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// normalizeEmail converts email to lowercase for consistent keys
func (s *AuthService) normalizeEmail(email string) string {
	return email
}

// StoreOTP in Redis with TTL (for registration)
func (s *AuthService) StoreOTP(email, otp string) error {
	key := fmt.Sprintf("otp:%s", email)

	log.Printf("🔐 Storing registration OTP for %s: %s", email, otp)

	if err := redis.Set(key, otp, 5*time.Minute); err != nil {
		log.Printf("Failed to store OTP in Redis: %v", err)
		return err
	}

	_, exists, err := redis.Get(key)
	if err != nil {
		log.Printf("Failed to verify OTP storage: %v", err)
	} else if exists {
		log.Printf("✅ Registration OTP stored successfully for %s", email)
	}

	return nil
}

// GetOTP from Redis (for registration)
func (s *AuthService) GetOTP(email string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)

	log.Printf("🔍 Retrieving registration OTP for %s", email)

	otp, exists, err := redis.Get(key)
	if err != nil {
		log.Printf("❌ Error getting registration OTP for %s: %v", email, err)
		return "", err
	}
	if !exists {
		log.Printf("❌ Registration OTP not found for %s", email)
		return "", errors.New("OTP not found")
	}

	log.Printf("✅ Registration OTP retrieved for %s: %s", email, otp)
	return otp, nil
}

// DeleteOTP from Redis (for registration)
func (s *AuthService) DeleteOTP(email string) error {
	key := fmt.Sprintf("otp:%s", email)

	log.Printf("🗑️ Deleting registration OTP from Redis: key=%s", key)

	if err := redis.Delete(key); err != nil {
		log.Printf("❌ Failed to delete registration OTP from Redis: %v", err)
		return err
	}

	log.Printf("✅ Registration OTP deleted from Redis: %s", key)
	return nil
}

// ================================================
// 2FA OTP MANAGEMENT
// ================================================

// StoreTwoFactorOTP stores 2FA OTP in Redis
func (s *AuthService) StoreTwoFactorOTP(email, otp string) error {
	key := fmt.Sprintf("2fa:%s", email)

	log.Printf("🔐 Storing 2FA OTP for %s: %s", email, otp)

	if err := redis.Set(key, otp, 5*time.Minute); err != nil {
		log.Printf("❌ Failed to store 2FA OTP: %v", err)
		return err
	}

	_, exists, err := redis.Get(key)
	if err != nil {
		log.Printf("❌ Failed to verify 2FA OTP storage: %v", err)
	} else if exists {
		log.Printf("✅ 2FA OTP stored successfully for %s", email)
	}

	return nil
}

// GetTwoFactorOTP retrieves 2FA OTP from Redis
func (s *AuthService) GetTwoFactorOTP(email string) (string, error) {
	key := fmt.Sprintf("2fa:%s", email)

	log.Printf("🔍 Retrieving 2FA OTP for %s", email)

	otp, exists, err := redis.Get(key)
	if err != nil {
		log.Printf("❌ Error getting 2FA OTP for %s: %v", email, err)
		return "", err
	}
	if !exists {
		log.Printf("❌ 2FA OTP not found for %s", email)
		return "", errors.New("2FA OTP not found or expired")
	}

	log.Printf("✅ 2FA OTP retrieved for %s: %s", email, otp)
	return otp, nil
}

// DeleteTwoFactorOTP removes 2FA OTP from Redis
func (s *AuthService) DeleteTwoFactorOTP(email string) error {
	key := fmt.Sprintf("2fa:%s", email)

	log.Printf("🗑️ Deleting 2FA OTP for %s", email)

	if err := redis.Delete(key); err != nil {
		log.Printf("❌ Failed to delete 2FA OTP: %v", err)
		return err
	}

	log.Printf("✅ 2FA OTP deleted for %s", email)
	return nil
}

// ================================================
// BUSINESS OTP MANAGEMENT
// ================================================

// StoreBusinessOTP stores OTP for business email verification
func (s *AuthService) StoreBusinessOTP(email, otp string) error {
	key := fmt.Sprintf("business_otp:%s", email)

	log.Printf("🔐 Storing business OTP for %s: %s", email, otp)

	if err := redis.Set(key, otp, 10*time.Minute); err != nil {
		log.Printf("Failed to store business OTP in Redis: %v", err)
		return err
	}

	_, exists, err := redis.Get(key)
	if err != nil {
		log.Printf("Failed to verify business OTP storage: %v", err)
	} else if exists {
		log.Printf("✅ Business OTP stored successfully for %s", email)
	}

	return nil
}

// GetBusinessOTP retrieves business OTP from Redis
func (s *AuthService) GetBusinessOTP(email string) (string, error) {
	key := fmt.Sprintf("business_otp:%s", email)

	log.Printf("🔍 Retrieving business OTP for %s", email)

	otp, exists, err := redis.Get(key)
	if err != nil {
		log.Printf("❌ Error getting business OTP for %s: %v", email, err)
		return "", err
	}
	if !exists {
		log.Printf("❌ Business OTP not found for %s", email)
		return "", errors.New("business OTP not found or expired")
	}

	log.Printf("✅ Business OTP retrieved for %s: %s", email, otp)
	return otp, nil
}

// DeleteBusinessOTP deletes business OTP from Redis
func (s *AuthService) DeleteBusinessOTP(email string) error {
	key := fmt.Sprintf("business_otp:%s", email)

	log.Printf("🗑️ Deleting business OTP from Redis: key=%s", key)

	if err := redis.Delete(key); err != nil {
		log.Printf("❌ Failed to delete business OTP from Redis: %v", err)
		return err
	}

	log.Printf("✅ Business OTP deleted from Redis: %s", key)
	return nil
}

// ================================================
// USER REGISTRATION DATA MANAGEMENT
// ================================================

// StoreUserData stores user registration data in Redis
func (s *AuthService) StoreUserData(email string, data map[string]interface{}) error {
	key := fmt.Sprintf("user:data:%s", email)

	log.Printf("📝 Storing user data in Redis: key=%s", key)

	if err := redis.HSet(key, data); err != nil {
		log.Printf("❌ Failed to store user data in Redis: %v", err)
		return err
	}

	if err := redis.Expire(key, 5*time.Minute); err != nil {
		log.Printf("⚠️ Failed to set TTL for user data: %v", err)
	}

	log.Printf("✅ User data stored in Redis: %s", key)
	return nil
}

// GetUserData gets user registration data from Redis
func (s *AuthService) GetUserData(email string) (map[string]string, error) {
	key := fmt.Sprintf("user:data:%s", email)

	log.Printf("🔍 Looking up user data in Redis: key=%s", key)

	result, exists, err := redis.HGetAll(key)
	if err != nil {
		log.Printf("❌ Redis error: %v", err)
		return nil, err
	}
	if !exists || len(result) == 0 {
		log.Printf("⚠️ User data not found in Redis: %s", key)
		return nil, errors.New("user data not found")
	}

	log.Printf("✅ User data found in Redis: %s", key)
	return result, nil
}

// DeleteUserData deletes user registration data from Redis
func (s *AuthService) DeleteUserData(email string) error {
	key := fmt.Sprintf("user:data:%s", email)

	log.Printf("🗑️ Deleting user data from Redis: key=%s", key)

	if err := redis.Delete(key); err != nil {
		log.Printf("❌ Failed to delete user data from Redis: %v", err)
		return err
	}

	log.Printf("✅ User data deleted from Redis: %s", key)
	return nil
}

// ================================================
// BUSINESS DATA MANAGEMENT
// ================================================

// StoreBusinessData stores business registration data in Redis
func (s *AuthService) StoreBusinessData(userID string, data map[string]interface{}) error {
	key := fmt.Sprintf("business:data:%s", userID)

	log.Printf("📝 Storing business data in Redis: key=%s", key)

	if err := redis.HSet(key, data); err != nil {
		log.Printf("❌ Failed to store business data in Redis: %v", err)
		return err
	}

	// Store for 24 hours to allow time for business email verification
	if err := redis.Expire(key, 24*time.Hour); err != nil {
		log.Printf("⚠️ Failed to set TTL for business data: %v", err)
	}

	log.Printf("✅ Business data stored in Redis: %s", key)
	return nil
}

// GetBusinessData retrieves business registration data from Redis
func (s *AuthService) GetBusinessData(userID string) (map[string]string, error) {
	key := fmt.Sprintf("business:data:%s", userID)

	log.Printf("🔍 Looking up business data in Redis: key=%s", key)

	result, exists, err := redis.HGetAll(key)
	if err != nil {
		log.Printf("❌ Redis error: %v", err)
		return nil, err
	}
	if !exists || len(result) == 0 {
		log.Printf("⚠️ Business data not found in Redis: %s", key)
		return nil, errors.New("business data not found")
	}

	log.Printf("✅ Business data found in Redis: %s", key)
	return result, nil
}

// DeleteBusinessData deletes business registration data from Redis
func (s *AuthService) DeleteBusinessData(userID string) error {
	key := fmt.Sprintf("business:data:%s", userID)

	log.Printf("🗑️ Deleting business data from Redis: key=%s", key)

	if err := redis.Delete(key); err != nil {
		log.Printf("❌ Failed to delete business data from Redis: %v", err)
		return err
	}

	log.Printf("✅ Business data deleted from Redis: %s", key)
	return nil
}

// MarkBusinessEmailVerified marks business email as verified
func (s *AuthService) MarkBusinessEmailVerified(userID string) error {
	key := fmt.Sprintf("business:data:%s", userID)

	log.Printf("✅ Marking business email as verified for user: %s", userID)

	if err := redis.HSetField(key, "business_email_verified", "true"); err != nil {
		log.Printf("❌ Failed to mark business email verified: %v", err)
		return err
	}

	return nil
}

// IsBusinessEmailVerified checks if business email is verified
func (s *AuthService) IsBusinessEmailVerified(userID string) (bool, error) {
	key := fmt.Sprintf("business:data:%s", userID)

	result, exists, err := redis.HGet(key, "business_email_verified")
	if err != nil {
		log.Printf("❌ Error checking business email verification: %v", err)
		return false, err
	}
	if !exists {
		return false, nil
	}

	return result == "true", nil
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
// AUTHENTICATION METHODS
// ================================================

// AuthenticateUser authenticates a user for login
func (s *AuthService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := user.ComparePassword(password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repo.UpdateUser(user); err != nil {
		log.Printf("Failed to update last login: %v", err)
	}

	return user, nil
}

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

	// Get primary business ID from BusinessMembers (first active business)
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
// PASSWORD RESET METHODS
// ================================================

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

// GeneratePasswordResetToken creates a password reset token
func (s *AuthService) GeneratePasswordResetToken(email string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}
	token := base64.URLEncoding.EncodeToString(bytes)

	key := fmt.Sprintf("password_reset:%s", token)
	if err := redis.Set(key, email, 15*time.Minute); err != nil {
		return "", fmt.Errorf("failed to store reset token: %w", err)
	}

	return token, nil
}

// ValidatePasswordResetToken validates a password reset token
func (s *AuthService) ValidatePasswordResetToken(token string) (string, error) {
	key := fmt.Sprintf("password_reset:%s", token)

	email, exists, err := redis.Get(key)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("invalid or expired reset token")
	}

	return email, nil
}

// DeletePasswordResetToken deletes a used reset token
func (s *AuthService) DeletePasswordResetToken(token string) error {
	key := fmt.Sprintf("password_reset:%s", token)
	return redis.Delete(key)
}

// ResetPassword resets a user's password
func (s *AuthService) ResetPassword(email, newPassword string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	return s.repo.UpdateUser(user)
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

// GetUserByIDWithBusiness gets a user with their business memberships
func (s *AuthService) GetUserByIDWithBusiness(id uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByIDWithBusiness(id)
}

// ================================================
// ADMIN METHODS
// ================================================

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

// ================================================
// RESPONSE MODELS
// ================================================

// UserInfo struct
type UserInfo struct {
	ID          string  `json:"id"`
	PhoneNumber string  `json":"phone_number"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	Avatar      string  `json:"avatar"`
	Role        string  `json:"role"`
	BusinessID  *string `json:"business_id,omitempty"`
}