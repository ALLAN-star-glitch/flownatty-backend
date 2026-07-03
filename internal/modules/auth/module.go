package auth

import (
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
    "github.com/gin-gonic/gin"
)

type AuthModule struct {
    handler *AuthHandler
}

func NewAuthModule(cfg *config.Config) *AuthModule {
    // Initialize repository
    repo := NewAuthRepository(database.GetDB())

    // Initialize service
    service := NewAuthService(repo, cfg)

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