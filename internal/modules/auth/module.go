package auth

import (
    "context"
    "log"

    "github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
    "github.com/ALLAN-star-glitch/flownatty-backend/pkg/queue"
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
)

type AuthModule struct {
    handler *AuthHandler
}

func NewAuthModule(cfg *config.Config) *AuthModule {
    // Initialize Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr: cfg.Redis.URL,
    })

    // ✅ Test Redis connection
    ctx := context.Background()
    if err := redisClient.Ping(ctx).Err(); err != nil {
        log.Printf("❌ Redis connection failed: %v", err)
        log.Printf("⚠️ Redis URL: %s", cfg.Redis.URL)
    } else {
        log.Println("✅ Redis connected successfully")
    }

    // Initialize queue client
    queueClient := queue.NewClient(cfg.Redis.URL)

    // Initialize repository
    repo := NewAuthRepository(database.GetDB())

    // Initialize service with Redis and queue
    service := NewAuthService(repo, cfg, queueClient, redisClient)

    // Initialize handler
    handler := NewAuthHandler(service, repo, cfg)

    return &AuthModule{
        handler: handler,
    }
}

func (m *AuthModule) RegisterRoutes(r *gin.RouterGroup) {
    auth := r.Group("/auth")
    {
        auth.POST("/register", m.handler.Register)
        auth.POST("/verify-otp", m.handler.VerifyOTP)
        auth.POST("/resend-otp", m.handler.ResendOTP)
    }
}