package auth

import (
	"context"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/auth/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/queue"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AuthModule struct {
	handler *AuthHandler
}

// NewAuthModule creates a new auth module with permissions and token services
func NewAuthModule(cfg *config.Config, permService *permissions.Service) *AuthModule {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.URL,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Redis connection failed: %v", err)
		log.Printf("Redis URL: %s", cfg.Redis.URL)
	} else {
		log.Println("Redis connected successfully")
	}

	// Initialize queue client
	queueClient := queue.NewClient(cfg.Redis.URL)

	// Initialize repository
	repo := NewAuthRepository(database.GetDB())

	// Initialize auth service
	authService := NewAuthService(repo, cfg, queueClient, redisClient)

	// Initialize token service
	tokenService := NewTokenService(repo, cfg)

	// Initialize handler with all dependencies
	handler := NewAuthHandler(authService, repo, cfg, permService, tokenService)

	return &AuthModule{
		handler: handler,
	}
}

// SetupRoutes registers all auth routes with the router
func (m *AuthModule) SetupRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		// ================================================
		// REGISTRATION FLOW (Public)
		// ================================================
		auth.POST("/register", m.handler.Register)
		auth.POST("/verify-otp", m.handler.VerifyOTP)
		auth.POST("/resend-otp", m.handler.ResendOTP)

		// ================================================
		// AUTHENTICATION FLOW (Public)
		// ================================================
		auth.POST("/login", m.handler.Login)

		// ================================================
		// TOKEN MANAGEMENT (Public - Cookie Based)
		// ================================================
		auth.POST("/refresh", m.handler.RefreshToken)

		// ================================================
		// SESSION MANAGEMENT (Cookie Based)
		// ================================================
		auth.POST("/logout", m.handler.Logout)

		// ================================================
		// PASSWORD RESET FLOW (Public - OTP Based)
		// ================================================
		// Step 1: User enters email + new password → OTP sent
		auth.POST("/forgot-password", m.handler.ForgotPassword)
		// Step 2: User enters OTP → Password reset
		auth.POST("/verify-reset-otp", m.handler.VerifyResetOTP)
	}
}

// GetHandler returns the auth handler
func (m *AuthModule) GetHandler() *AuthHandler {
	return m.handler
}