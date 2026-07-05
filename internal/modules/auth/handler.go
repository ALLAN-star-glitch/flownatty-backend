package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/auth/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/email"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/validation"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var validator = validation.New()

type AuthHandler struct {
	service      *AuthService
	repo         *AuthRepository
	config       *config.Config
	emailService *email.EmailService
	permService  *permissions.Service
	tokenService *TokenService
}

// NewAuthHandler creates a new auth handler with all dependencies
func NewAuthHandler(
	service *AuthService,
	repo *AuthRepository,
	cfg *config.Config,
	permService *permissions.Service,
	tokenService *TokenService,
) *AuthHandler {
	emailSvc := email.NewEmailService(cfg.Resend.ApiKey, cfg.Resend.From)

	return &AuthHandler{
		service:      service,
		repo:         repo,
		config:       cfg,
		emailService: emailSvc,
		permService:  permService,
		tokenService: tokenService,
	}
}

// ================================================
// REQUEST MODELS WITH EXAMPLES
// ================================================

// RegisterRequest represents registration request
type RegisterRequest struct {
	Phone    string `json:"phone_number" binding:"required" example:"+254712345678"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Password string `json:"password" binding:"required,min=8" example:"SecurePass123!"`
	Role     string `json:"role" binding:"omitempty,oneof=consumer business" example:"consumer"`
}

// VerifyOTPRequest represents OTP verification request
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// ResendOTPRequest represents resend OTP request
type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"omitempty,email" example:"john@example.com"`
	Phone    string `json:"phone_number" example:"+254712345678"`
	Password string `json:"password" binding:"required" example:"SecurePass123!"`
}

// RefreshTokenRequest represents refresh token request (for mobile)
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"abc123xyz789..."`
}

// LogoutRequest represents logout request (for mobile)
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"abc123xyz789..."`
}

// ForgotPasswordRequest represents forgot password request (Step 1)
type ForgotPasswordRequest struct {
	Email       string `json:"email" binding:"required,email" example:"john@example.com"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"SecurePass123!"`
}

// VerifyResetOTPRequest represents OTP verification for password reset (Step 2)
type VerifyResetOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// ================================================
// RESPONSE MODELS WITH EXAMPLES
// ================================================

// OTPResponse represents OTP response
type OTPResponse struct {
	Email     string    `json:"email" example:"john@example.com"`
	Phone     string    `json:"phone" example:"+254712345678"`
	Role      string    `json:"role" example:"consumer"`
	ExpiresAt time.Time `json:"expires_at" example:"2026-07-06T17:40:00Z"`
}

// UserResponse represents user response after registration
type UserResponse struct {
	TokenType string `json:"token_type" example:"Bearer"`
	ExpiresIn int64  `json:"expires_in" example:"86400"`
	User      struct {
		ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
		Phone string `json:"phone" example:"+254712345678"`
		Email string `json:"email" example:"john@example.com"`
		Name  string `json:"name" example:"John Doe"`
		Role  string `json:"role" example:"consumer"`
	} `json:"user"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string   `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string   `json:"refresh_token" example:"abc123xyz789..."`
	TokenType    string   `json:"token_type" example:"Bearer"`
	ExpiresIn    int64    `json:"expires_in" example:"86400"`
	User         UserInfo `json:"user"`
}

// RefreshResponse represents refresh response
type RefreshResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string `json:"refresh_token" example:"abc123xyz789..."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int64  `json:"expires_in" example:"86400"`
}

// ForgotPasswordResponse represents forgot password response
type ForgotPasswordResponse struct {
	Message string `json:"message" example:"OTP sent to your email"`
	Expires int    `json:"expires_in" example:"300"`
}

// VerifyResetOTPResponse represents OTP verification response
type VerifyResetOTPResponse struct {
	Message string `json:"message" example:"Password reset successfully"`
}

// ================================================
// COOKIE HELPERS
// ================================================

// setAccessTokenCookie sets the access token as an HTTP-only cookie
func (h *AuthHandler) setAccessTokenCookie(c *gin.Context, token string) {
	c.SetCookie(
		"access_token",
		token,
		int(h.config.JWT.Expiration.Seconds()),
		"/",
		"",
		h.config.Environment == "production",
		true,
	)
}

// setRefreshTokenCookie sets the refresh token as an HTTP-only cookie
func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetCookie(
		"refresh_token",
		token,
		int(h.config.JWT.RefreshExpiration.Seconds()),
		"/auth/refresh",
		"",
		h.config.Environment == "production",
		true,
	)
}

// clearAuthCookies clears both auth cookies
func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", h.config.Environment == "production", true)
	c.SetCookie("refresh_token", "", -1, "/auth/refresh", "", h.config.Environment == "production", true)
}

// getRefreshTokenFromCookie extracts refresh token from cookie
func (h *AuthHandler) getRefreshTokenFromCookie(c *gin.Context) (string, error) {
	token, err := c.Cookie("refresh_token")
	if err != nil {
		return "", errors.New("refresh token not found in cookie")
	}
	if token == "" {
		return "", errors.New("refresh token is empty")
	}
	return token, nil
}

// ================================================
// REGISTER HANDLER
// ================================================

// Register Handler
// @Summary Register a new user
// @Description Register a new user with phone, email, and password. Sends OTP via email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 200 {object} response.BaseResponse{data=OTPResponse}
// @Failure 400 {object} response.BaseResponse
// @Failure 409 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	//  Validate phone number
	if !validator.Phone.Validate(req.Phone) {
		response.BadRequest(c, "Invalid Kenyan phone number", gin.H{
			"field": "phone_number",
			"format": "Use format: 254XXXXXXXXX, +254XXXXXXXXX, 07XXXXXXXX, 01XXXXXXXX, or 02XXXXXXXX",
		})
		return
	}
	req.Phone = validator.Phone.Normalize(req.Phone)

	// Validate password strength
	if valid, msg := validator.Password.Validate(req.Password); !valid {
		response.BadRequest(c, "Weak password", gin.H{
			"field": "password",
			"error": msg,
			"strength": validator.Password.StrengthLabel(validator.Password.Score(req.Password)),
		})
		return
	}

	if req.Role == "" {
		req.Role = permissions.RoleConsumer.String()
	}

	if req.Role != permissions.RoleConsumer.String() && req.Role != permissions.RoleBusinessOwner.String() {
		response.BadRequest(c, "Invalid role", gin.H{
			"allowed_roles": []string{
				permissions.RoleConsumer.String(),
				permissions.RoleBusinessOwner.String(),
			},
			"provided": req.Role,
		})
		return
	}

	existingUser, _ := h.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		response.Conflict(c, "Email already registered", gin.H{
			"field": "email",
			"value": req.Email,
		})
		return
	}

	existingPhone, _ := h.repo.GetUserByPhone(req.Phone)
	if existingPhone != nil {
		response.Conflict(c, "Phone number already registered", gin.H{
			"field": "phone",
			"value": req.Phone,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "Failed to process password", gin.H{
			"error": err.Error(),
		})
		return
	}

	otp := h.service.GenerateOTP()
	log.Printf("Generated OTP for %s: %s", req.Email, otp)

	if err := h.service.StoreOTP(req.Email, otp); err != nil {
		log.Printf("Failed to store OTP in Redis: %v", err)
		response.InternalError(c, "Failed to process request", gin.H{
			"error": err.Error(),
		})
		return
	}

	userData := map[string]interface{}{
		"phone":    req.Phone,
		"email":    req.Email,
		"name":     req.Name,
		"password": string(hashedPassword),
		"role":     req.Role,
	}
	if err := h.service.StoreUserData(req.Email, userData); err != nil {
		log.Printf("Failed to store user data in Redis: %v", err)
		response.InternalError(c, "Failed to process request", gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.EnqueueOTPEmail(req.Email, req.Name, otp); err != nil {
		log.Printf("Failed to enqueue OTP email: %v", err)
		if err := h.emailService.SendOTP(email.OTPEmailData{
			To:      req.Email,
			Name:    req.Name,
			OTP:     otp,
			Expires: "5 minutes",
		}); err != nil {
			log.Printf("Failed to send OTP email: %v", err)
		}
	}

	respData := gin.H{
		"email":      req.Email,
		"phone":      req.Phone,
		"role":       req.Role,
		"expires_at": time.Now().Add(5 * time.Minute),
	}

	response.Success(c, "OTP sent successfully", respData)
}

// ================================================
// VERIFY OTP HANDLER
// ================================================

// VerifyOTP Handler
// @Summary Verify OTP and create account
// @Description Verify the OTP sent to the user's email and create the account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "OTP verification details"
// @Success 201 {object} response.BaseResponse{data=UserResponse}
// @Failure 400 {object} response.BaseResponse
// @Failure 409 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("Verifying OTP for %s", req.Email)

	storedOTP, err := h.service.GetOTP(req.Email)
	if err != nil {
		response.BadRequest(c, "Invalid or expired OTP", gin.H{
			"email": req.Email,
		})
		return
	}

	if req.OTP != storedOTP {
		response.BadRequest(c, "Invalid OTP", gin.H{
			"email":         req.Email,
			"provided_otp": req.OTP,
		})
		return
	}

	log.Printf("OTP verified for %s", req.Email)

	userData, err := h.service.GetUserData(req.Email)
	if err != nil {
		response.BadRequest(c, "Registration data not found. Please register again.", gin.H{
			"email": req.Email,
		})
		return
	}

	h.service.DeleteOTP(req.Email)
	h.service.DeleteUserData(req.Email)

	h.createUserFromData(c, req.Email, userData)
}

// Helper to create user from Redis data
func (h *AuthHandler) createUserFromData(c *gin.Context, userEmail string, userData map[string]string) {
	existingUser, _ := h.repo.GetUserByEmail(userEmail)
	if existingUser != nil {
		response.Conflict(c, "User already registered", gin.H{
			"field": "email",
			"value": userEmail,
		})
		return
	}

	existingPhone, _ := h.repo.GetUserByPhone(userData["phone"])
	if existingPhone != nil {
		response.Conflict(c, "Phone number already registered", gin.H{
			"field": "phone",
			"value": userData["phone"],
		})
		return
	}

	now := time.Now()
	user := &models.User{
		PhoneNumber:     userData["phone"],
		Email:           userData["email"],
		Password:        userData["password"],
		Name:            userData["name"],
		Role:            userData["role"],
		IsVerified:      true,
		IsEmailVerified: false,
		VerifiedAt:      &now,
		LastActiveAt:    &now,
	}

	if err := h.repo.CreateUser(user); err != nil {
		log.Printf("Failed to create user: %v", err)
		response.InternalError(c, "Failed to create user", gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := context.Background()
	switch user.Role {
	case permissions.RoleConsumer.String():
		if err := h.permService.AssignConsumerRole(ctx, user.ID.String()); err != nil {
			log.Printf("Failed to assign consumer role to user %s: %v", user.ID, err)
		} else {
			log.Printf("Assigned consumer role to user: %s", user.ID)
		}

	case permissions.RoleBusinessOwner.String():
		if err := h.permService.AssignConsumerRole(ctx, user.ID.String()); err != nil {
			log.Printf("Failed to assign consumer role to user %s: %v", user.ID, err)
		} else {
			log.Printf("Assigned consumer role to business user: %s", user.ID)
		}
		log.Printf("Business user created: %s - waiting for business creation to assign owner role", user.ID)

	default:
		log.Printf("Unknown role: %s, skipping permission assignment", user.Role)
	}

	accessToken, err := h.tokenService.GenerateAccessToken(user)
	if err != nil {
		response.InternalError(c, "Failed to generate token", gin.H{
			"error": err.Error(),
		})
		return
	}

	h.setAccessTokenCookie(c, accessToken)

	if err := h.service.EnqueueWelcomeEmail(user.Email, user.Name); err != nil {
		log.Printf("Failed to enqueue welcome email: %v", err)
		if err := h.emailService.SendWelcome(email.WelcomeEmailData{
			To:   user.Email,
			Name: user.Name,
		}); err != nil {
			log.Printf("Failed to send welcome email: %v", err)
		}
	}

	response.Created(c, "Account created successfully", gin.H{
		"token_type": "Bearer",
		"expires_in": int64(h.config.JWT.Expiration.Seconds()),
		"user": gin.H{
			"id":    user.ID,
			"phone": user.PhoneNumber,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

// ================================================
// RESEND OTP HANDLER
// ================================================

// ResendOTP Handler
// @Summary Resend OTP
// @Description Resend OTP to the user's email for registration
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ResendOTPRequest true "Email to resend OTP"
// @Success 200 {object} response.BaseResponse{data=OTPResponse}
// @Failure 400 {object} response.BaseResponse
// @Failure 409 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req ResendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	existingUser, _ := h.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		response.Conflict(c, "Email already registered", gin.H{
			"field": "email",
			"value": req.Email,
		})
		return
	}

	newOTP := h.service.GenerateOTP()

	if err := h.service.StoreOTP(req.Email, newOTP); err != nil {
		log.Printf("Failed to store OTP in Redis: %v", err)
		response.InternalError(c, "Failed to resend OTP", gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.StoreUserData(req.Email, map[string]interface{}{
		"refresh": time.Now().Unix(),
	}); err != nil {
		log.Printf("Failed to refresh user data TTL: %v", err)
	}

	log.Printf("New OTP generated for %s: %s", req.Email, newOTP)

	if err := h.emailService.SendOTP(email.OTPEmailData{
		To:      req.Email,
		Name:    "User",
		OTP:     newOTP,
		Expires: "5 minutes",
	}); err != nil {
		log.Printf("Failed to resend OTP email: %v", err)
	}

	respData := gin.H{
		"email":      req.Email,
		"expires_at": time.Now().Add(5 * time.Minute),
	}

	response.Success(c, "OTP resent successfully", respData)
}

// ================================================
// LOGIN HANDLER (Supports both Web + Mobile)
// ================================================

// Login Handler
// @Summary Login user
// @Description Authenticate user by email or phone. For web: tokens set in cookies. For mobile: tokens returned in response body.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials (email or phone)"
// @Success 200 {object} response.BaseResponse{data=LoginResponse}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	// ✅ Authenticate by email or phone
	var user *models.User
	var err error

	if req.Email != "" {
		user, err = h.service.AuthenticateUserByEmail(req.Email, req.Password)
	} else if req.Phone != "" {
		// ✅ Validate and normalize phone
		if !validator.Phone.Validate(req.Phone) {
			response.BadRequest(c, "Invalid Kenyan phone number", gin.H{
				"field": "phone_number",
				"format": "Use format: 254XXXXXXXXX, +254XXXXXXXXX, 07XXXXXXXX, 01XXXXXXXX, or 02XXXXXXXX",
			})
			return
		}
		req.Phone = validator.Phone.Normalize(req.Phone)
		user, err = h.service.AuthenticateUserByPhone(req.Phone, req.Password)
	} else {
		response.BadRequest(c, "Email or phone number required", nil)
		return
	}

	if err != nil {
		response.Unauthorized(c, "Invalid credentials", gin.H{
			"error": err.Error(),
		})
		return
	}

	accessToken, err := h.tokenService.GenerateAccessToken(user)
	if err != nil {
		response.InternalError(c, "Failed to generate access token", gin.H{
			"error": err.Error(),
		})
		return
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(
		user.ID,
		c.GetHeader("User-Agent"),
		c.ClientIP(),
	)
	if err != nil {
		response.InternalError(c, "Failed to generate refresh token", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set cookies for web
	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, refreshToken)

	userInfo := h.service.BuildUserInfo(user)

	// Enqueue login notification email
	if err := h.service.EnqueueLoginNotificationEmail(
		user.Email,
		user.Name,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	); err != nil {
		log.Printf("Failed to enqueue login notification email: %v", err)
		// Fallback: send synchronously
		if err := h.emailService.SendLoginNotification(email.LoginNotificationData{
			To:        user.Email,
			Name:      user.Name,
			Time:      time.Now().Format("2006-01-02 15:04:05 UTC"),
			IPAddress: c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		}); err != nil {
			log.Printf("Failed to send login notification email: %v", err)
		}
	}

	// Return tokens in body for mobile
	response.Success(c, "Login successful", LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.config.JWT.Expiration.Seconds()),
		User:         userInfo,
	})
}

// ================================================
// REFRESH TOKEN HANDLER (Supports both Web + Mobile)
// ================================================

// RefreshToken Handler
// @Summary Refresh access token
// @Description Get a new access token. For web: uses refresh token from cookie. For mobile: send refresh_token in body.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token (required for mobile, optional for web)"
// @Success 200 {object} response.BaseResponse{data=RefreshResponse}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var refreshToken string

	// Priority 1: Get from request body (mobile)
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		refreshToken = req.RefreshToken
	}

	// Priority 2: Get from cookie (web)
	if refreshToken == "" {
		cookieToken, err := h.getRefreshTokenFromCookie(c)
		if err == nil && cookieToken != "" {
			refreshToken = cookieToken
		}
	}

	if refreshToken == "" {
		response.Unauthorized(c, "Refresh token required", gin.H{
			"reason": "provide refresh_token in body or cookie",
		})
		return
	}

	newAccessToken, newRefreshToken, err := h.tokenService.RotateRefreshToken(
		refreshToken,
		c.GetHeader("User-Agent"),
		c.ClientIP(),
	)
	if err != nil {
		response.Unauthorized(c, "Invalid refresh token", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set new cookies for web
	h.setAccessTokenCookie(c, newAccessToken)
	h.setRefreshTokenCookie(c, newRefreshToken)

	// Return tokens in body for mobile
	response.Success(c, "Token refreshed successfully", RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.config.JWT.Expiration.Seconds()),
	})
}

// ================================================
// LOGOUT HANDLER (Supports both Web + Mobile)
// ================================================

// Logout Handler
// @Summary Logout user
// @Description Revoke refresh token. For web: uses cookie. For mobile: send refresh_token in body.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Refresh token (required for mobile, optional for web)"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var refreshToken string

	// Priority 1: Get from request body (mobile)
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		refreshToken = req.RefreshToken
	}

	// Priority 2: Get from cookie (web)
	if refreshToken == "" {
		cookieToken, err := h.getRefreshTokenFromCookie(c)
		if err == nil && cookieToken != "" {
			refreshToken = cookieToken
		}
	}

	if refreshToken != "" {
		if err := h.tokenService.RevokeRefreshToken(refreshToken); err != nil {
			log.Printf("Failed to revoke refresh token: %v", err)
		}
	}

	// Clear cookies for web
	h.clearAuthCookies(c)

	response.Success(c, "Logged out successfully", nil)
}

// ================================================
// FORGOT PASSWORD HANDLER (Step 1)
// ================================================

// ForgotPassword Handler
// @Summary Forgot password - Step 1
// @Description Send OTP to user's email for password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email and new password"
// @Success 200 {object} response.BaseResponse{data=ForgotPasswordResponse}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate new password strength
	if valid, msg := validator.Password.Validate(req.NewPassword); !valid {
		response.BadRequest(c, "Weak password", gin.H{
			"field": "new_password",
			"error": msg,
			"strength": validator.Password.StrengthLabel(validator.Password.Score(req.NewPassword)),
		})
		return
	}

	// Check if user exists first (without revealing existence)
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		response.InternalError(c, "Failed to process request", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Security: Don't reveal if user exists
	if user == nil {
		response.Success(c, "OTP sent to your email", ForgotPasswordResponse{
			Message: "If your email is registered, you will receive an OTP",
			Expires: 300,
		})
		return
	}

	// Generate OTP
	otp := h.service.GenerateOTP()

	// Store OTP and new password in Redis with TTL (5 minutes)
	ctx := context.Background()
	key := fmt.Sprintf("reset:%s", req.Email)
	data := map[string]interface{}{
		"otp":          otp,
		"new_password": req.NewPassword,
	}
	if err := h.service.redis.HSet(ctx, key, data).Err(); err != nil {
		log.Printf("Failed to store reset data: %v", err)
		response.InternalError(c, "Failed to process request", gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := h.service.redis.Expire(ctx, key, 5*time.Minute).Err(); err != nil {
		log.Printf("Failed to set TTL: %v", err)
	}

	log.Printf("Password reset OTP for %s: %s", req.Email, otp)

	// Enqueue password reset OTP email via Asynq
	if err := h.service.EnqueuePasswordResetOTPEmail(req.Email, user.Name, otp); err != nil {
		log.Printf("Failed to enqueue password reset OTP email: %v", err)
		// Fallback: send synchronously
		if err := h.emailService.SendPasswordResetOTP(email.PasswordResetOTPData{
			To:      req.Email,
			Name:    user.Name,
			OTP:     otp,
			Expires: "5 minutes",
		}); err != nil {
			log.Printf("Failed to send password reset OTP email: %v", err)
		}
	}

	response.Success(c, "OTP sent to your email", ForgotPasswordResponse{
		Message: "Check your email for the OTP",
		Expires: 300,
	})
}

// ================================================
// VERIFY RESET OTP HANDLER (Step 2)
// ================================================

// VerifyResetOTP Handler
// @Summary Verify reset OTP - Step 2
// @Description Verify OTP and reset password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyResetOTPRequest true "OTP verification details"
// @Success 200 {object} response.BaseResponse{data=VerifyResetOTPResponse}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /auth/verify-reset-otp [post]
func (h *AuthHandler) VerifyResetOTP(c *gin.Context) {
	var req VerifyResetOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := context.Background()
	key := fmt.Sprintf("reset:%s", req.Email)

	// Get stored data from Redis
	data, err := h.service.redis.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("Failed to get reset data: %v", err)
		response.InternalError(c, "Failed to verify OTP", gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(data) == 0 {
		response.Unauthorized(c, "Invalid or expired OTP", gin.H{
			"error": "OTP not found or expired",
		})
		return
	}

	// Verify OTP
	storedOTP, ok := data["otp"]
	if !ok || req.OTP != storedOTP {
		response.Unauthorized(c, "Invalid OTP", gin.H{
			"error": "Invalid OTP code",
		})
		return
	}

	// Get new password
	newPassword, ok := data["new_password"]
	if !ok {
		response.InternalError(c, "Invalid reset data", gin.H{
			"error": "Reset data corrupted",
		})
		return
	}

	// Reset password
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		response.NotFound(c, "User not found", gin.H{
			"error": "User not found",
		})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "Failed to process password", gin.H{
			"error": err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if err := h.repo.UpdateUser(user); err != nil {
		response.InternalError(c, "Failed to update password", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete reset data from Redis
	h.service.redis.Del(ctx, key)

	// Enqueue password reset confirmation email via Asynq
	if err := h.service.EnqueuePasswordResetConfirmEmail(req.Email, user.Name); err != nil {
		log.Printf("Failed to enqueue password reset confirmation email: %v", err)
		// Fallback: send synchronously
		if err := h.emailService.SendPasswordResetConfirm(email.PasswordResetConfirmData{
			To:   req.Email,
			Name: user.Name,
		}); err != nil {
			log.Printf("Failed to send password reset confirmation email: %v", err)
		}
	}

	response.Success(c, "Password reset successfully", VerifyResetOTPResponse{
		Message: "Password reset successfully",
	})
}