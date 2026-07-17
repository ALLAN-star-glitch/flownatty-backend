package authservice

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authrepo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TokenService handles all token generation and management
type TokenService struct {
	repo   *authrepo.AuthRepository
	config *config.Config
}

// NewTokenService creates a new token service instance
func NewTokenService(repo *authrepo.AuthRepository, cfg *config.Config) *TokenService {
	return &TokenService{
		repo:   repo,
		config: cfg,
	}
}

// GenerateAccessToken creates a new access token (JWT)
// This is the ONLY place access tokens are generated
func (s *TokenService) GenerateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(s.config.JWT.Expiration).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

// GenerateRefreshToken creates a new refresh token
// Stores the token in the database for validation and revocation
func (s *TokenService) GenerateRefreshToken(userID uuid.UUID, userAgent, ipAddress string) (string, error) {
    // ✅ Generate 32 random bytes
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate refresh token: %w", err)
    }
    // ✅ Encode to Base64 URL-safe
    token := base64.URLEncoding.EncodeToString(bytes)  // ← This creates "zQA-gcU8f3e7iM3BjTZRDwYJ1Cdb5czB0Wb0cGGDzE8="

    // ✅ Store in database
    refreshToken := &models.RefreshToken{
        UserID:    userID,
        Token:     token,
        ExpiresAt: time.Now().Add(s.config.JWT.RefreshExpiration),
        Revoked:   false,
        UserAgent: userAgent,
        IPAddress: ipAddress,
    }

    if err := s.repo.CreateRefreshToken(refreshToken); err != nil {
        return "", fmt.Errorf("failed to store refresh token: %w", err)
    }

    return token, nil
}

// ValidateRefreshToken validates a refresh token
// Checks: exists, not revoked, not expired
func (s *TokenService) ValidateRefreshToken(token string) (*models.RefreshToken, error) {
	refreshToken, err := s.repo.GetRefreshTokenByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid refresh token")
		}
		return nil, err
	}

	if refreshToken.Revoked {
		return nil, errors.New("refresh token has been revoked")
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, errors.New("refresh token has expired")
	}

	return refreshToken, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (s *TokenService) RefreshAccessToken(refreshTokenStr string) (string, error) {
	// Validate the refresh token
	refreshToken, err := s.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return "", err
	}

	// Get the user
	user, err := s.repo.GetUserByID(refreshToken.UserID.String())
	if err != nil {
		return "", errors.New("user not found")
	}

	// Generate a new access token
	return s.GenerateAccessToken(user)
}

// RevokeRefreshToken revokes a single refresh token (logout)
func (s *TokenService) RevokeRefreshToken(token string) error {
	return s.repo.RevokeRefreshToken(token)
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user
// Useful for: password change, account compromise, security audit
func (s *TokenService) RevokeAllUserRefreshTokens(userID uuid.UUID) error {
	return s.repo.RevokeAllUserRefreshTokens(userID)
}

// RotateRefreshToken invalidates the old token and generates a new one
// This is more secure for each refresh operation
func (s *TokenService) RotateRefreshToken(oldToken string, userAgent, ipAddress string) (string, string, error) {
	// Validate the old token
	refreshToken, err := s.ValidateRefreshToken(oldToken)
	if err != nil {
		return "", "", err
	}

	// Get the user
	user, err := s.repo.GetUserByID(refreshToken.UserID.String())
	if err != nil {
		return "", "", errors.New("user not found")
	}

	// Revoke the old token
	if err := s.repo.RevokeRefreshToken(oldToken); err != nil {
		return "", "", fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Generate a new refresh token
	newRefreshToken, err := s.GenerateRefreshToken(user.ID, userAgent, ipAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// Generate a new access token
	newAccessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}