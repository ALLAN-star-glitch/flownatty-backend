package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	repo   *AuthRepository
	config *config.Config
}

func NewAuthService(repo *AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo:   repo,
		config: cfg,
	}
}

// Generate OTP - Always returns random 6-digit OTP
func (s *AuthService) GenerateOTP() string {
	// Always generate random OTP for both development and production
	// rand is automatically seeded in Go 1.20+
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Generate JWT token
func (s *AuthService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
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