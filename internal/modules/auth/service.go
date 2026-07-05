package auth

import (
    "context"
    "errors"
    "fmt"
    "log"
    "math/rand"
    "time"

    "github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
    "github.com/ALLAN-star-glitch/flownatty-backend/pkg/queue"
    "github.com/golang-jwt/jwt/v5"
    "github.com/hibiken/asynq"
    "github.com/redis/go-redis/v9"
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

// Generate OTP - ALWAYS random
func (s *AuthService) GenerateOTP() string {
    return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Store OTP in Redis with TTL
func (s *AuthService) StoreOTP(email, otp string) error {
    ctx := context.Background()
    key := fmt.Sprintf("otp:%s", email)

    log.Printf("📝 Storing OTP in Redis: key=%s, value=%s", key, otp)

    err := s.redis.Set(ctx, key, otp, s.config.OTP.TTL).Err()
    if err != nil {
        log.Printf("❌ Failed to store OTP in Redis: %v", err)
        return err
    }

    // Verify it was stored
    val, _ := s.redis.Get(ctx, key).Result()
    log.Printf("✅ OTP stored and verified in Redis: %s -> %s", key, val)

    return nil
}

// Get OTP from Redis
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

// Delete OTP from Redis
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

// Store user registration data in Redis
func (s *AuthService) StoreUserData(email string, data map[string]interface{}) error {
    ctx := context.Background()
    key := fmt.Sprintf("user:data:%s", email)

    log.Printf("📝 Storing user data in Redis: key=%s", key)

    err := s.redis.HSet(ctx, key, data).Err()
    if err != nil {
        log.Printf("❌ Failed to store user data in Redis: %v", err)
        return err
    }

    // Set TTL to match OTP expiry (5 minutes)
    err = s.redis.Expire(ctx, key, s.config.OTP.TTL).Err()
    if err != nil {
        log.Printf("⚠️ Failed to set TTL for user data: %v", err)
    }

    log.Printf("✅ User data stored in Redis: %s", key)
    return nil
}

// Get user registration data from Redis
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

// Delete user registration data from Redis
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

// Enqueue OTP email via Asynq
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

// Enqueue Welcome Email via Asynq
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

// Generate JWT token
func (s *AuthService) GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role, 
		"exp":     time.Now().Add(s.config.JWT.Expiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}


// Validate JWT token
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