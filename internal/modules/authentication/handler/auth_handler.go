package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authrepo"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/service"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/repository"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/email"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var validator = validation.New()

type AuthHandler struct {
	service      *service.AuthService
	repo         *authrepo.AuthRepository
	memberRepo   *repository.BusinessMemberRepository
	config       *config.Config
	emailService *email.EmailService
	permService  *permissions.Service
	tokenService *service.TokenService
}

// NewAuthHandler creates a new auth handler with all dependencies
func NewAuthHandler(
	service *service.AuthService,
	repo *authrepo.AuthRepository,
	memberRepo *repository.BusinessMemberRepository,
	cfg *config.Config,
	permService *permissions.Service,
	tokenService *service.TokenService,
) *AuthHandler {
	emailSvc := email.NewEmailService(cfg.Resend.ApiKey, cfg.Resend.From)

	return &AuthHandler{
		service:      service,
		repo:         repo,
		memberRepo:   memberRepo,
		config:       cfg,
		emailService: emailSvc,
		permService:  permService,
		tokenService: tokenService,
	}
}

// ================================================
// REQUEST MODELS
// ================================================

// RegisterRequest represents registration request
type RegisterRequest struct {
	Phone    string `json:"phone_number" binding:"required" example:"+254712345678"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Password string `json:"password" binding:"required,min=8" example:"SecurePass123!"`
	Role     string `json:"role" binding:"omitempty,oneof=consumer business_owner" example:"consumer"`

	// Business fields (required if role is business_owner)
	BusinessType     string `json:"business_type" binding:"omitempty,oneof=sole_proprietor partnership company cooperative"`
	BusinessName     string `json:"business_name"`
	BusinessCategory string `json:"business_category"`
	BusinessPhone    string `json:"business_phone"`
	BusinessEmail    string `json:"business_email"`
	BusinessAddress  string `json:"business_address"`
	BusinessDesc     string `json:"business_description"`
}

// VerifyOTPRequest represents OTP verification request
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// VerifyBusinessEmailRequest represents business email verification request
type VerifyBusinessEmailRequest struct {
	BusinessEmail string `json:"business_email" binding:"required,email" example:"info@business.com"`
	OTP           string `json:"otp" binding:"required,len=6" example:"123456"`
	UserID        string `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// ResendBusinessOTPRequest represents resend business OTP request
type ResendBusinessOTPRequest struct {
	BusinessEmail string `json:"business_email" binding:"required,email" example:"info@business.com"`
	UserID        string `json:"user_id" binding:"required"`
}

// ResendOTPRequest represents resend OTP request
type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePass123!"`
}

// RefreshTokenRequest represents refresh token request (for mobile)
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"abc123xyz789..."`
}

type VerifyTwoFactorOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// LogoutRequest represents logout request (for mobile)
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"abc123xyz789..."`
}

type TwoFactorAuthResponse struct {
	Requires2FA bool   `json:"requires_2fa" example:"true"`
	Email       string `json:"email" example:"john@example.com"`
	ExpiresIn   int    `json:"expires_in" example:"300"`
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
// RESPONSE MODELS
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
	AccessToken  string   `json:"access_token,omitempty"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int64    `json:"expires_in"`
	User        service.UserInfo `json:"user"`
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
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 400 {object} response.BaseResponse
// @Failure 409 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate phone number
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

	// Set default role
	if req.Role == "" {
		req.Role = permissions.RoleConsumer.String()
	}

	// Validate role
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

	// Validate business fields if role is business
	if req.Role == permissions.RoleBusinessOwner.String() {
		if req.BusinessType == "" {
			response.BadRequest(c, "Business type is required for business registration", gin.H{
				"field": "business_type",
				"allowed_types": []string{"sole_proprietor", "partnership", "company", "cooperative"},
			})
			return
		}
		if req.BusinessName == "" {
			response.BadRequest(c, "Business name is required for business registration", gin.H{
				"field": "business_name",
			})
			return
		}
		if req.BusinessCategory == "" {
			response.BadRequest(c, "Business category is required for business registration", gin.H{
				"field": "business_category",
				"allowed_categories": []string{"retail", "fashion", "beauty", "electronics", "food", "health"},
			})
			return
		}
		if req.BusinessPhone == "" {
			response.BadRequest(c, "Business phone is required for business registration", gin.H{
				"field": "business_phone",
			})
			return
		}
		// Validate business phone
		if !validator.Phone.Validate(req.BusinessPhone) {
			response.BadRequest(c, "Invalid business phone number", gin.H{
				"field": "business_phone",
			})
			return
		}
		req.BusinessPhone = validator.Phone.Normalize(req.BusinessPhone)

		// For non-sole_proprietor, business email is required
		if req.BusinessType != "sole_proprietor" && req.BusinessEmail == "" {
			response.BadRequest(c, "Business email is required for partnerships, companies, and cooperatives", gin.H{
				"field": "business_email",
			})
			return
		}
	}

	// Check if user exists
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

	// Store business data if role is business
	if req.Role == permissions.RoleBusinessOwner.String() {
		userData["business_type"] = req.BusinessType
		userData["business_name"] = req.BusinessName
		userData["business_category"] = req.BusinessCategory
		userData["business_phone"] = req.BusinessPhone
		userData["business_email"] = req.BusinessEmail
		userData["business_address"] = req.BusinessAddress
		userData["business_description"] = req.BusinessDesc
	}

	if err := h.service.StoreUserData(req.Email, userData); err != nil {
		log.Printf("Failed to store user data in Redis: %v", err)
		response.InternalError(c, "Failed to process request", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Enqueue OTP email via Asynq
	if err := h.service.EnqueueOTPEmail(req.Email, req.Name, otp); err != nil {
		log.Printf("Failed to enqueue OTP email: %v", err)
		if err := h.emailService.SendSignupOTP(req.Email, req.Name, otp, "5 minutes"); err != nil {
			log.Printf("Failed to send OTP email: %v", err)
		}
	}

	respData := gin.H{
		"email":      req.Email,
		"phone":      req.Phone,
		"role":       req.Role,
		"expires_at": time.Now().Add(5 * time.Minute),
	}

	if req.Role == permissions.RoleBusinessOwner.String() {
		respData["business_type"] = req.BusinessType
		respData["business_name"] = req.BusinessName
		respData["message"] = "OTP sent to your email. After verification, complete your business setup."
		respData["next_step"] = "verify_personal_email"

		if req.BusinessType != "sole_proprietor" && req.BusinessEmail != "" {
			respData["business_email_required"] = true
			respData["business_email"] = req.BusinessEmail
			respData["message"] = "OTP sent to your email. After personal verification, you'll need to verify your business email."
		}
	} else {
		respData["message"] = "OTP sent successfully. Verify to complete registration."
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
// @Router /api/v1/auth/verify-otp [post]
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

	// Check if this is a business user
	isBusiness := userData["role"] == permissions.RoleBusinessOwner.String()
	businessType := userData["business_type"]

	// For non-sole_proprietor businesses, create user and store business data for later verification
	if isBusiness && businessType != "sole_proprietor" {
		h.createUserAndStoreBusinessData(c, req.Email, userData)
		return
	}

	// For consumers and sole proprietors, create user and business immediately
	h.createUserFromData(c, req.Email, userData)
}

// ================================================
// CREATE USER WITH BUSINESS DATA
// ================================================

// createUserAndStoreBusinessData creates a user and stores business data for later verification
func (h *AuthHandler) createUserAndStoreBusinessData(c *gin.Context, userEmail string, userData map[string]string) {
	// Create user first
	user, err := h.createUserOnly(c, userEmail, userData)
	if err != nil {
		return
	}

	// Store business data in Redis for onboarding
	businessData := map[string]interface{}{
		"user_id":              user.ID.String(),
		"business_type":        userData["business_type"],
		"business_name":        userData["business_name"],
		"business_category":    userData["business_category"],
		"business_phone":       userData["business_phone"],
		"business_email":       userData["business_email"],
		"business_address":     userData["business_address"],
		"business_description": userData["business_description"],
		"personal_email":       userData["email"],
		"personal_name":        userData["name"],
	}

	if err := h.service.StoreBusinessData(user.ID.String(), businessData); err != nil {
		log.Printf("Failed to store business data: %v", err)
		response.InternalError(c, "Failed to process request", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Generate tokens
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

	// Set cookies
	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, refreshToken)

	// Send business verification email
	businessOTP := h.service.GenerateOTP()
	if err := h.service.StoreBusinessOTP(userData["business_email"], businessOTP); err != nil {
		log.Printf("Failed to store business OTP: %v", err)
	}

	if err := h.service.EnqueueBusinessOTPEmail(userData["business_email"], userData["business_name"], businessOTP); err != nil {
		log.Printf("Failed to send business OTP: %v", err)
		if err := h.emailService.SendBusinessOTP(userData["business_email"], userData["business_name"], businessOTP, "10 minutes"); err != nil {
			log.Printf("Failed to send business OTP email: %v", err)
		}
	}

	response.Created(c, "Account created. Please verify your business email.", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    int64(h.config.JWT.Expiration.Seconds()),
		"user": gin.H{
			"id":    user.ID,
			"phone": user.PhoneNumber,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
		"business_verification_required": true,
		"business_email":                 userData["business_email"],
		"business_name":                  userData["business_name"],
		"message":                        "Please verify your business email to complete registration.",
		"next_step":                      "verify_business_email",
	})
}

// ================================================
// CREATE USER ONLY (HELPER)
// ================================================

// createUserOnly creates a user without creating business
func (h *AuthHandler) createUserOnly(c *gin.Context, userEmail string, userData map[string]string) (*models.User, error) {
	existingUser, _ := h.repo.GetUserByEmail(userEmail)
	if existingUser != nil {
		response.Conflict(c, "User already registered", gin.H{
			"field": "email",
			"value": userEmail,
		})
		return nil, errors.New("user already exists")
	}

	existingPhone, _ := h.repo.GetUserByPhone(userData["phone"])
	if existingPhone != nil {
		response.Conflict(c, "Phone number already registered", gin.H{
			"field": "phone",
			"value": userData["phone"],
		})
		return nil, errors.New("phone already exists")
	}

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

	if err := h.repo.CreateUser(user); err != nil {
		log.Printf("Failed to create user: %v", err)
		response.InternalError(c, "Failed to create user", gin.H{
			"error": err.Error(),
		})
		return nil, err
	}

	ctx := context.Background()

	// Assign consumer role (base role for all users)
	if err := h.permService.AssignConsumerRole(ctx, user.ID.String()); err != nil {
		log.Printf("Failed to assign consumer role to user %s: %v", user.ID, err)
	} else {
		log.Printf("Assigned consumer role to user: %s", user.ID)
	}

	log.Printf("User created: %s - waiting for business verification", user.ID)

	return user, nil
}

// ================================================
// CREATE USER FROM DATA
// ================================================

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
		// Assign consumer role (base role)
		if err := h.permService.AssignConsumerRole(ctx, user.ID.String()); err != nil {
			log.Printf("Failed to assign consumer role to user %s: %v", user.ID, err)
		} else {
			log.Printf("Assigned consumer role to business user: %s", user.ID)
		}

		// For sole proprietors, create business immediately
		businessType := userData["business_type"]
		if businessType == "sole_proprietor" {
			// Store business data
			businessData := map[string]interface{}{
				"user_id":              user.ID.String(),
				"business_type":        businessType,
				"business_name":        userData["business_name"],
				"business_category":    userData["business_category"],
				"business_phone":       userData["business_phone"],
				"business_email":       userData["email"], // Use personal email as business email
				"business_address":     userData["business_address"],
				"business_description": userData["business_description"],
			}

			if err := h.service.StoreBusinessData(user.ID.String(), businessData); err != nil {
				log.Printf("Failed to store business data: %v", err)
			}

			// Mark business email as verified for sole proprietors
			if err := h.service.MarkBusinessEmailVerified(user.ID.String()); err != nil {
				log.Printf("Failed to mark business email verified: %v", err)
			}

			// Create business
			if err := h.createBusinessFromData(c, user.ID, businessData); err != nil {
				log.Printf("Failed to create business: %v", err)
				response.InternalError(c, "Failed to create business", gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		log.Printf("Business user created: %s - waiting for business creation", user.ID)

	default:
		log.Printf("Unknown role: %s, skipping permission assignment", user.Role)
	}

	// Generate tokens
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

	// Set cookies
	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, refreshToken)

	// Send welcome email
	if err := h.service.EnqueueWelcomeEmail(user.Email, user.Name); err != nil {
		log.Printf("Failed to enqueue welcome email: %v", err)
		if err := h.emailService.SendWelcome(email.WelcomeEmailData{
			To:   user.Email,
			Name: user.Name,
		}); err != nil {
			log.Printf("Failed to send welcome email: %v", err)
		}
	}

	// Return response
	responseData := gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    int64(h.config.JWT.Expiration.Seconds()),
		"user": gin.H{
			"id":    user.ID,
			"phone": user.PhoneNumber,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	}

	if user.Role == permissions.RoleBusinessOwner.String() {
		businessType := userData["business_type"]

		if businessType == "sole_proprietor" {
			responseData["onboarding_required"] = true
			responseData["onboarding_url"] = "/api/v1/businesses/onboarding"
			responseData["business_type"] = businessType
			responseData["business_name"] = userData["business_name"]
			responseData["message"] = "Account created. Please complete your business profile."
			responseData["business_email_verified"] = true
		} else {
			responseData["business_verification_required"] = true
			responseData["business_email"] = userData["business_email"]
			responseData["business_name"] = userData["business_name"]
			responseData["business_type"] = businessType
			responseData["message"] = "Account created. Please verify your business email to continue."
			responseData["next_step"] = "verify_business_email"
		}
	} else {
		responseData["message"] = "Account created successfully. Start shopping!"
	}

	response.Created(c, "Account created successfully", responseData)
}

// ================================================
// CREATE BUSINESS FROM DATA
// ================================================

// createBusinessFromData creates a business from stored data
func (h *AuthHandler) createBusinessFromData(c *gin.Context, userID uuid.UUID, businessData map[string]interface{}) error {
	business := &models.Business{
		Name:        businessData["business_name"].(string),
		Category:    businessData["business_category"].(string),
		Phone:       businessData["business_phone"].(string),
		Email:       businessData["business_email"].(string),
		Address:     businessData["business_address"].(string),
		Description: businessData["business_description"].(string),
		IsVerified:  false,
		IsActive:    true,
	}

	if err := h.repo.CreateBusiness(business); err != nil {
		return fmt.Errorf("failed to create business: %w", err)
	}

	// Add user as owner using BusinessMemberRepository
	member := &models.BusinessMember{
		BusinessID: business.ID,
		UserID:     userID,
		Role:       permissions.RoleBusinessOwner.String(),
		IsActive:   true,
		JoinedAt:   time.Now(),
	}

	if err := h.memberRepo.Create(member); err != nil {
		return fmt.Errorf("failed to add business owner: %w", err)
	}

	// Assign business owner role in permissions
	ctx := context.Background()
	if err := h.permService.AssignBusinessOwnerRole(ctx, userID.String(), business.ID.String()); err != nil {
		return fmt.Errorf("failed to assign permissions: %w", err)
	}

	// Delete business data from Redis
	if err := h.service.DeleteBusinessData(userID.String()); err != nil {
		log.Printf("Failed to delete business data: %v", err)
	}

	log.Printf("Business created: %s (ID: %s) for user: %s", business.Name, business.ID, userID)

	return nil
}

// ================================================
// VERIFY BUSINESS EMAIL HANDLER
// ================================================

// VerifyBusinessEmail godoc
// @Summary Verify business email
// @Description Verify business email with OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyBusinessEmailRequest true "Business email verification details"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/auth/verify-business-email [post]
func (h *AuthHandler) VerifyBusinessEmail(c *gin.Context) {
	var req VerifyBusinessEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("Verifying business email: %s for user: %s", req.BusinessEmail, req.UserID)

	// Get stored OTP
	storedOTP, err := h.service.GetBusinessOTP(req.BusinessEmail)
	if err != nil {
		response.BadRequest(c, "Invalid or expired OTP", gin.H{
			"business_email": req.BusinessEmail,
		})
		return
	}

	if req.OTP != storedOTP {
		response.BadRequest(c, "Invalid OTP", gin.H{
			"business_email": req.BusinessEmail,
			"provided_otp":   req.OTP,
		})
		return
	}

	// Get business data from Redis
	businessData, err := h.service.GetBusinessData(req.UserID)
	if err != nil {
		response.BadRequest(c, "Business data not found", gin.H{
			"user_id": req.UserID,
		})
		return
	}

	// Mark business email as verified
	if err := h.service.MarkBusinessEmailVerified(req.UserID); err != nil {
		log.Printf("Failed to mark business email verified: %v", err)
		response.InternalError(c, "Failed to verify business email", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete OTP
	if err := h.service.DeleteBusinessOTP(req.BusinessEmail); err != nil {
		log.Printf("Failed to delete business OTP: %v", err)
	}

	// Create business
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	if err := h.createBusinessFromData(c, userID, map[string]interface{}{
		"business_name":        businessData["business_name"],
		"business_category":    businessData["business_category"],
		"business_phone":       businessData["business_phone"],
		"business_email":       req.BusinessEmail,
		"business_address":     businessData["business_address"],
		"business_description": businessData["business_description"],
	}); err != nil {
		response.InternalError(c, "Failed to create business", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send business welcome email
	if err := h.service.EnqueueBusinessWelcomeEmail(req.BusinessEmail, businessData["business_name"], businessData["personal_name"]); err != nil {
		log.Printf("Failed to enqueue business welcome email: %v", err)
		if err := h.emailService.SendBusinessWelcome(email.BusinessWelcomeData{
			To:           req.BusinessEmail,
			BusinessName: businessData["business_name"],
			OwnerName:    businessData["personal_name"],
		}); err != nil {
			log.Printf("Failed to send business welcome email: %v", err)
		}
	}

	response.Success(c, "Business email verified successfully", gin.H{
		"business_name":  businessData["business_name"],
		"business_email": req.BusinessEmail,
		"status":         "verified",
		"message":        "Business email verified. Your business is now active.",
	})
}

// ================================================
// RESEND BUSINESS OTP HANDLER
// ================================================

// ResendBusinessOTP godoc
// @Summary Resend business OTP
// @Description Resend OTP to business email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ResendBusinessOTPRequest true "Resend business OTP details"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/auth/resend-business-otp [post]
func (h *AuthHandler) ResendBusinessOTP(c *gin.Context) {
	var req ResendBusinessOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get business data
	businessData, err := h.service.GetBusinessData(req.UserID)
	if err != nil {
		response.BadRequest(c, "Business data not found", gin.H{
			"user_id": req.UserID,
		})
		return
	}

	// Generate new OTP
	newOTP := h.service.GenerateOTP()

	// Store OTP
	if err := h.service.StoreBusinessOTP(req.BusinessEmail, newOTP); err != nil {
		log.Printf("Failed to store business OTP: %v", err)
		response.InternalError(c, "Failed to resend OTP", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send OTP email
	if err := h.service.EnqueueBusinessOTPEmail(req.BusinessEmail, businessData["business_name"], newOTP); err != nil {
		log.Printf("Failed to send business OTP: %v", err)
		if err := h.emailService.SendBusinessOTP(req.BusinessEmail, businessData["business_name"], newOTP, "10 minutes"); err != nil {
			log.Printf("Failed to send business OTP email: %v", err)
		}
	}

	response.Success(c, "OTP resent to business email", gin.H{
		"business_email": req.BusinessEmail,
		"expires_in":     600,
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
// @Router /api/v1/auth/resend-otp [post]
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

	if err := h.emailService.SendSignupOTP(req.Email, "User", newOTP, "5 minutes"); err != nil {
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
// @Description Authenticate user by email. Sends 2FA OTP to email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials (email)"
// @Success 200 {object} response.BaseResponse{data=TwoFactorAuthResponse}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	var user *models.User
	var err error

	if req.Email != "" {
		user, err = h.service.AuthenticateUserByEmail(req.Email, req.Password)
	} else {
		response.BadRequest(c, "Email required", nil)
		return
	}

	if err != nil {
		response.Unauthorized(c, "Invalid credentials", gin.H{
			"error": err.Error(),
		})
		return
	}

	// 2FA: Generate and send OTP
	otp := h.service.GenerateOTP()

	if err := h.service.StoreTwoFactorOTP(user.Email, otp); err != nil {
		response.InternalError(c, "Failed to process request", gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("2FA OTP for %s: %s", user.Email, otp)

	if err := h.service.EnqueueTwoFactorOTPEmail(user.Email, user.Name, otp); err != nil {
		log.Printf("Failed to enqueue 2FA OTP email: %v", err)
		if err := h.emailService.SendTwoFactorOTP(user.Email, user.Name, otp, "5 minutes"); err != nil {
			log.Printf("Failed to send 2FA OTP email: %v", err)
		}
	}

	response.Success(c, "2FA verification required", TwoFactorAuthResponse{
		Requires2FA: true,
		Email:       user.Email,
		ExpiresIn:   300,
	})
}

// ================================================
// VERIFY TWO-FACTOR OTP HANDLER
// ================================================

// VerifyTwoFactorOTP verifies OTP for two-factor authentication
// @Summary Verify 2FA OTP
// @Description Verify the OTP sent to user's email for two-factor authentication and issue access tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyTwoFactorOTPRequest true "2FA OTP verification details"
// @Success 200 {object} response.BaseResponse{data=LoginResponse}
// @Failure 400 {object} response.BaseResponse
// @Failure 401 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/auth/verify-2fa [post]
func (h *AuthHandler) VerifyTwoFactorOTP(c *gin.Context) {
	var req VerifyTwoFactorOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("Verifying 2FA OTP for %s", req.Email)

	storedOTP, err := h.service.GetTwoFactorOTP(req.Email)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired OTP", gin.H{
			"error": "OTP not found or expired. Please login again.",
		})
		return
	}

	if req.OTP != storedOTP {
		response.Unauthorized(c, "Invalid OTP", gin.H{
			"error": "The OTP you entered is incorrect. Please try again.",
		})
		return
	}

	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		response.NotFound(c, "User not found", nil)
		return
	}

	if err := h.service.DeleteTwoFactorOTP(req.Email); err != nil {
		log.Printf("Failed to delete 2FA OTP: %v", err)
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

	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, refreshToken)

	userInfo := h.service.BuildUserInfo(user)

	if err := h.service.EnqueueLoginNotificationEmail(
		user.Email,
		user.Name,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	); err != nil {
		log.Printf("Failed to enqueue login notification: %v", err)
	}

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
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var refreshToken string

	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		refreshToken = req.RefreshToken
	}

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

	h.setAccessTokenCookie(c, newAccessToken)
	h.setRefreshTokenCookie(c, newRefreshToken)

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
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var refreshToken string

	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		refreshToken = req.RefreshToken
	}

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
// @Router /api/v1/auth/forgot-password [post]
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
			"field":      "new_password",
			"error":      msg,
			"strength":   validator.Password.StrengthLabel(validator.Password.Score(req.NewPassword)),
		})
		return
	}

	// Check if user exists (without revealing existence)
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

	// Store OTP and new password using service method
	if err := h.service.StoreResetData(req.Email, otp, req.NewPassword); err != nil {
		log.Printf("Failed to store reset data: %v", err)
		response.InternalError(c, "Failed to process request", gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("Password reset OTP for %s: %s", req.Email, otp)

	// Enqueue password reset OTP email
	if err := h.service.EnqueuePasswordResetOTPEmail(req.Email, user.Name, otp); err != nil {
		log.Printf("Failed to enqueue password reset OTP email: %v", err)
		// Fallback: send synchronously
		if err := h.emailService.SendPasswordResetOTP(req.Email, user.Name, otp, "5 minutes"); err != nil {
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
// @Router /api/v1/auth/verify-reset-otp [post]
func (h *AuthHandler) VerifyResetOTP(c *gin.Context) {
	var req VerifyResetOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get stored data from Redis using service method
	data, err := h.service.GetResetData(req.Email)
	if err != nil {
		if err.Error() == "reset data not found" {
			response.Unauthorized(c, "Invalid or expired OTP", gin.H{
				"error": "OTP not found or expired",
			})
			return
		}
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

	// Get user
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

	// Update user in database
	if err := h.repo.UpdateUser(user); err != nil {
		response.InternalError(c, "Failed to update password", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete reset data from Redis using service method
	if err := h.service.DeleteResetData(req.Email); err != nil {
		log.Printf("Failed to delete reset data: %v", err)
	}

	// Enqueue password reset confirmation email
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