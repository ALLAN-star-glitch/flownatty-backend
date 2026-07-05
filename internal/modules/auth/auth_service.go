package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/queue"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *AuthRepository
	config *config.Config
	queue  *queue.Client
	redis  *redis.Client
}

func NewAuthService(repo *AuthRepository, cfg *config.Config, queueClient *queue.Client, redisClient *redis.Client) *AuthService {
	return &AuthService{
		repo:   repo,
		config: cfg,
		queue:  queueClient,
		redis:  redisClient,
	}
}

// ================================================
// OTP GENERATION & MANAGEMENT
// ================================================

// GenerateOTP - ALWAYS random
func (s *AuthService) GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// StoreOTP in Redis with TTL
func (s *AuthService) StoreOTP(email, otp string) error {
	ctx := context.Background()
	key := fmt.Sprintf("otp:%s", email)

	log.Printf("📝 Storing OTP in Redis: key=%s, value=%s", key, otp)

	err := s.redis.Set(ctx, key, otp, s.config.OTP.TTL).Err()
	if err != nil {
		log.Printf("❌ Failed to store OTP in Redis: %v", err)
		return err
	}

	val, _ := s.redis.Get(ctx, key).Result()
	log.Printf("✅ OTP stored and verified in Redis: %s -> %s", key, val)

	return nil
}

// GetOTP from Redis
func (s *AuthService) GetOTP(email string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("otp:%s", email)

	log.Printf("🔍 Looking up OTP in Redis: key=%s", key)

	otp, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("⚠️ OTP not found in Redis: %s", key)
			return "", errors.New("OTP expired or not found")
		}
		log.Printf("❌ Redis error: %v", err)
		return "", err
	}

	log.Printf("✅ OTP found in Redis: %s -> %s", key, otp)
	return otp, nil
}

// DeleteOTP from Redis
func (s *AuthService) DeleteOTP(email string) error {
	ctx := context.Background()
	key := fmt.Sprintf("otp:%s", email)

	log.Printf("🗑️ Deleting OTP from Redis: key=%s", key)

	err := s.redis.Del(ctx, key).Err()
	if err != nil {
		log.Printf("❌ Failed to delete OTP from Redis: %v", err)
		return err
	}

	log.Printf("✅ OTP deleted from Redis: %s", key)
	return nil
}

// ================================================
// USER REGISTRATION DATA MANAGEMENT
// ================================================

// StoreUserData stores user registration data in Redis
func (s *AuthService) StoreUserData(email string, data map[string]interface{}) error {
	ctx := context.Background()
	key := fmt.Sprintf("user:data:%s", email)

	log.Printf("📝 Storing user data in Redis: key=%s", key)

	err := s.redis.HSet(ctx, key, data).Err()
	if err != nil {
		log.Printf("❌ Failed to store user data in Redis: %v", err)
		return err
	}

	err = s.redis.Expire(ctx, key, s.config.OTP.TTL).Err()
	if err != nil {
		log.Printf("⚠️ Failed to set TTL for user data: %v", err)
	}

	log.Printf("✅ User data stored in Redis: %s", key)
	return nil
}

// GetUserData gets user registration data from Redis
func (s *AuthService) GetUserData(email string) (map[string]string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:data:%s", email)

	log.Printf("🔍 Looking up user data in Redis: key=%s", key)

	result, err := s.redis.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("❌ Redis error: %v", err)
		return nil, err
	}

	if len(result) == 0 {
		log.Printf("⚠️ User data not found in Redis: %s", key)
		return nil, errors.New("user data not found")
	}

	log.Printf("✅ User data found in Redis: %s", key)
	return result, nil
}

// DeleteUserData deletes user registration data from Redis
func (s *AuthService) DeleteUserData(email string) error {
	ctx := context.Background()
	key := fmt.Sprintf("user:data:%s", email)

	log.Printf("🗑️ Deleting user data from Redis: key=%s", key)

	err := s.redis.Del(ctx, key).Err()
	if err != nil {
		log.Printf("❌ Failed to delete user data from Redis: %v", err)
		return err
	}

	log.Printf("✅ User data deleted from Redis: %s", key)
	return nil
}

// ================================================
// EMAIL ENQUEUE METHODS
// ================================================

// EnqueueOTPEmail enqueues OTP email via Asynq
func (s *AuthService) EnqueueOTPEmail(to, name, otp string) error {
	task := OTPEmailTask{
		To:      to,
		Name:    name,
		OTP:     otp,
		Expires: "5 minutes",
	}

	payload, err := task.Payload()
	if err != nil {
		return err
	}

	return s.queue.Enqueue(TypeEmailOTP, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueueWelcomeEmail enqueues welcome email via Asynq
func (s *AuthService) EnqueueWelcomeEmail(to, name string) error {
	task := WelcomeEmailTask{
		To:   to,
		Name: name,
	}

	payload, err := task.Payload()
	if err != nil {
		return err
	}

	return s.queue.Enqueue(TypeEmailWelcome, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueuePasswordResetOTPEmail enqueues password reset OTP email via Asynq
func (s *AuthService) EnqueuePasswordResetOTPEmail(to, name, otp string) error {
	task := PasswordResetOTPTask{
		To:      to,
		Name:    name,
		OTP:     otp,
		Expires: "5 minutes",
	}

	payload, err := task.Payload()
	if err != nil {
		return err
	}

	return s.queue.Enqueue(TypeEmailPasswordResetOTP, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueueLoginNotificationEmail enqueues login notification email via Asynq
func (s *AuthService) EnqueueLoginNotificationEmail(to, name, ipAddress, userAgent string) error {
	task := LoginNotificationTask{
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

	return s.queue.Enqueue(TypeEmailLoginNotification, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
}

// EnqueuePasswordResetConfirmEmail enqueues password reset confirmation email via Asynq
func (s *AuthService) EnqueuePasswordResetConfirmEmail(to, name string) error {
	task := PasswordResetConfirmTask{
		To:   to,
		Name: name,
	}

	payload, err := task.Payload()
	if err != nil {
		return err
	}

	return s.queue.Enqueue(TypeEmailPasswordResetConfirm, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
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

	// Update last login
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

	if user.BusinessID != nil {
		businessID := user.BusinessID.String()
		info.BusinessID = &businessID
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

// GeneratePasswordResetToken creates a password reset token
func (s *AuthService) GeneratePasswordResetToken(email string) (string, error) {
	ctx := context.Background()

	// Check if user exists
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Generate a secure random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}
	token := base64.URLEncoding.EncodeToString(bytes)

	// Store in Redis with TTL (15 minutes)
	key := fmt.Sprintf("password_reset:%s", token)
	err = s.redis.Set(ctx, key, email, 15*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store reset token: %w", err)
	}

	return token, nil
}

// ValidatePasswordResetToken validates a password reset token
func (s *AuthService) ValidatePasswordResetToken(token string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("password_reset:%s", token)

	email, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("invalid or expired reset token")
		}
		return "", err
	}

	return email, nil
}

// DeletePasswordResetToken deletes a used reset token
func (s *AuthService) DeletePasswordResetToken(token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("password_reset:%s", token)
	return s.redis.Del(ctx, key).Err()
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

	// Hash new password
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

	// Update last login
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

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repo.UpdateUser(user); err != nil {
		log.Printf("Failed to update last login: %v", err)
	}

	return user, nil
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