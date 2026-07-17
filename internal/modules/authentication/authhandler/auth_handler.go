// internal/modules/authentication/handler/auth_handler.go

package authhandler

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authrepo"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/authservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizservice"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/validators/bizvalidator"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/email"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var validator = validation.New()

// ================================================
// RESPONSE MODELS (For Swagger)
// ================================================

// OTPResponse represents OTP response
type OTPResponse struct {
	Email     string    `json:"email" example:"john@example.com"`
	Phone     string    `json:"phone" example:"+254712345678"`
	Role      string    `json:"role" example:"consumer"`
	ExpiresAt time.Time `json:"expires_at" example:"2026-07-06T17:40:00Z"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string               `json:"access_token,omitempty"`
	RefreshToken string               `json:"refresh_token,omitempty"`
	TokenType    string               `json:"token_type"`
	ExpiresIn    int64                `json:"expires_in"`
	User         authservice.UserInfo `json:"user"`
}

// RefreshResponse represents refresh response
type RefreshResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string `json:"refresh_token" example:"abc123xyz789..."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int64  `json:"expires_in" example:"86400"`
}

// TwoFactorAuthResponse represents 2FA response
type TwoFactorAuthResponse struct {
	Requires2FA bool   `json:"requires_2fa" example:"true"`
	Email       string `json:"email" example:"john@example.com"`
	ExpiresIn   int    `json:"expires_in" example:"300"`
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
// REQUEST MODELS
// ================================================

type RegisterRequest struct {
	Phone    string `json:"phone_number" binding:"required" example:"+254712345678"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Password string `json:"password" binding:"required,min=8" example:"SecurePass123!"`
	Role     string `json:"role" binding:"omitempty,oneof=consumer business_admin" example:"consumer"`

	BusinessType     string `json:"business_type"`
	BusinessName     string `json:"business_name"`
	BusinessCategory string `json:"business_category"`
	BusinessPhone    string `json:"business_phone"`
	BusinessEmail    string `json:"business_email"`
	BusinessAddress  string `json:"business_address"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

type VerifyBusinessEmailRequest struct {
	BusinessEmail string `json:"business_email" binding:"required,email" example:"info@business.com"`
	OTP           string `json:"otp" binding:"required,len=6" example:"123456"`
	UserID        string `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type ResendBusinessOTPRequest struct {
	BusinessEmail string `json:"business_email" binding:"required,email" example:"info@business.com"`
	UserID        string `json:"user_id" binding:"required"`
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePass123!"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"abc123xyz789..."`
}

type VerifyTwoFactorOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"abc123xyz789..."`
}

type ForgotPasswordRequest struct {
	Email       string `json:"email" binding:"required,email" example:"john@example.com"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"SecurePass123!"`
}

type VerifyResetOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// ================================================
// HANDLER STRUCT
// ================================================

type AuthHandler struct {
	service         *authservice.AuthService
	repo            *authrepo.AuthRepository
	config          *config.Config
	emailService    *email.EmailService
	permService     *permissions.Service
	tokenService    *authservice.TokenService
	businessService *bizservice.BusinessService
}

func NewAuthHandler(
	service *authservice.AuthService,
	repo *authrepo.AuthRepository,
	cfg *config.Config,
	permService *permissions.Service,
	tokenService *authservice.TokenService,
	businessService *bizservice.BusinessService,
) *AuthHandler {
	emailSvc := email.NewEmailService(cfg.Resend.ApiKey, cfg.Resend.From)
	return &AuthHandler{
		service:         service,
		repo:            repo,
		config:          cfg,
		emailService:    emailSvc,
		permService:     permService,
		tokenService:    tokenService,
		businessService: businessService,
	}
}

// ================================================
// COOKIE HELPERS
// ================================================

func (h *AuthHandler) setAccessTokenCookie(c *gin.Context, token string) {
	c.SetCookie("access_token", token, int(h.config.JWT.Expiration.Seconds()), "/", "", h.config.Environment == "production", true)
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetCookie("refresh_token", token, int(h.config.JWT.RefreshExpiration.Seconds()), "/auth/refresh", "", h.config.Environment == "production", true)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", h.config.Environment == "production", true)
	c.SetCookie("refresh_token", "", -1, "/auth/refresh", "", h.config.Environment == "production", true)
}

func (h *AuthHandler) getRefreshTokenFromCookie(c *gin.Context) (string, error) {
	token, err := c.Cookie("refresh_token")
	if err != nil || token == "" {
		return "", errors.New("refresh token not found")
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
		response.BadRequest(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	// Validate phone
	if !validator.Phone.Validate(req.Phone) {
		response.BadRequest(c, "Invalid Kenyan phone number", gin.H{
			"field":  "phone_number",
			"format": "Use format: 254XXXXXXXXX, +254XXXXXXXXX, 07XXXXXXXX, 01XXXXXXXX, or 02XXXXXXXX",
		})
		return
	}
	req.Phone = validator.Phone.Normalize(req.Phone)

	// Validate password
	if valid, msg := validator.Password.Validate(req.Password); !valid {
		response.BadRequest(c, "Weak password", gin.H{
			"field":    "password",
			"error":    msg,
			"strength": validator.Password.StrengthLabel(validator.Password.Score(req.Password)),
		})
		return
	}

	// Set default role
	if req.Role == "" {
		req.Role = permissions.RoleConsumer.String()
	}

	// Validate role
	if req.Role != permissions.RoleConsumer.String() && req.Role != permissions.RoleBusinessAdmin.String() {
		response.BadRequest(c, "Invalid role", gin.H{
			"allowed_roles": []string{permissions.RoleConsumer.String(), permissions.RoleBusinessAdmin.String()},
			"provided":      req.Role,
		})
		return
	}

	// Validate business fields if role is business_admin
	if req.Role == permissions.RoleBusinessAdmin.String() {
		bizValidator := bizvalidator.NewBusinessValidator()

		if req.BusinessType == "" || bizValidator.ValidateBusinessType(req.BusinessType) != nil {
			response.BadRequest(c, "Valid business type is required", gin.H{"field": "business_type"})
			return
		}
		if req.BusinessName == "" || bizValidator.ValidateBusinessName(req.BusinessName) != nil {
			response.BadRequest(c, "Valid business name is required", gin.H{"field": "business_name"})
			return
		}
		if req.BusinessCategory == "" || bizValidator.ValidateBusinessCategory(req.BusinessCategory) != nil {
			response.BadRequest(c, "Valid business category is required", gin.H{"field": "business_category"})
			return
		}
		if req.BusinessPhone == "" || bizValidator.ValidateBusinessPhone(req.BusinessPhone) != nil {
			response.BadRequest(c, "Valid business phone is required", gin.H{"field": "business_phone"})
			return
		}
		req.BusinessPhone = bizValidator.NormalizeBusinessPhone(req.BusinessPhone)
		if req.BusinessEmail == "" || bizValidator.ValidateBusinessEmailRequired(req.BusinessEmail) != nil {
			response.BadRequest(c, "Valid business email is required", gin.H{"field": "business_email"})
			return
		}
		if req.BusinessAddress == "" || bizValidator.ValidateBusinessAddress(req.BusinessAddress) != nil {
			response.BadRequest(c, "Valid business address is required", gin.H{"field": "business_address"})
			return
		}
	}

	// Prepare user data
	userData := map[string]string{
		"phone":    req.Phone,
		"email":    req.Email,
		"name":     req.Name,
		"password": req.Password,
		"role":     req.Role,
	}

	if req.Role == permissions.RoleBusinessAdmin.String() {
		userData["business_type"] = req.BusinessType
		userData["business_name"] = req.BusinessName
		userData["business_category"] = req.BusinessCategory
		userData["business_phone"] = req.BusinessPhone
		userData["business_email"] = req.BusinessEmail
		userData["business_address"] = req.BusinessAddress
	}

	// Delegate to service
	ctx := context.Background()
	if err := h.service.RegisterUser(ctx, userData); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	respData := gin.H{
		"email":      req.Email,
		"phone":      req.Phone,
		"role":       req.Role,
		"expires_at": time.Now().Add(5 * time.Minute),
		"message":    "OTP sent successfully. Verify to complete registration.",
	}

	if req.Role == permissions.RoleBusinessAdmin.String() {
		respData["business_type"] = req.BusinessType
		respData["business_name"] = req.BusinessName
		respData["message"] = "OTP sent to your email. After verification, complete your business setup."
		respData["next_step"] = "verify_personal_email"
		if req.BusinessEmail != "" {
			respData["business_email"] = req.BusinessEmail
			respData["message"] = "OTP sent to your email. After personal verification, you'll need to verify your business email."
		}
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
// @Success 201 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 400 {object} response.BaseResponse
// @Failure 409 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /api/v1/auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	user, business, result, err := h.service.VerifyOTPAndCreateUser(ctx, req.Email, req.OTP)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Set cookies
	if accessToken, ok := result["access_token"].(string); ok && accessToken != "" {
		h.setAccessTokenCookie(c, accessToken)
	}
	if refreshToken, ok := result["refresh_token"].(string); ok && refreshToken != "" {
		h.setRefreshTokenCookie(c, refreshToken)
	}

	// Build response
	responseData := gin.H{
		"token_type": "Bearer",
		"expires_in": int64(h.config.JWT.Expiration.Seconds()),
		"user": gin.H{
			"id":    user.ID,
			"phone": user.PhoneNumber,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	}

	// Merge result data
	for key, value := range result {
		if key != "user" && key != "business" && key != "business_data" {
			responseData[key] = value
		}
	}

	if business != nil {
		responseData["business_id"] = business.ID
		responseData["business_name"] = business.Name
	}

	response.Created(c, "Account created successfully", responseData)
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
		response.BadRequest(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	ctx := context.Background()
	business, result, err := h.service.CompleteBusinessVerification(ctx, userID, req.BusinessEmail, req.OTP)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Business email verified successfully", gin.H{
		"business_name":  result["business_name"],
		"business_email": result["business_email"],
		"business_id":    business.ID,
		"status":         result["status"],
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
		response.BadRequest(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	// Get business data
	businessData, err := h.service.GetBusinessData(req.UserID)
	if err != nil {
		response.BadRequest(c, "Business data not found", nil)
		return
	}

	// Generate new OTP
	newOTP := h.service.GenerateOTP()
	if err := h.service.StoreBusinessOTP(req.BusinessEmail, newOTP); err != nil {
		response.InternalError(c, "Failed to resend OTP", gin.H{"error": err.Error()})
		return
	}

	// Send OTP email
	if err := h.service.EnqueueBusinessOTPEmail(req.BusinessEmail, businessData["business_name"], newOTP); err != nil {
		log.Printf("Failed to send business OTP: %v", err)
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
		response.BadRequest(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	existingUser, _ := h.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		response.Conflict(c, "Email already registered", gin.H{"field": "email", "value": req.Email})
		return
	}

	newOTP := h.service.GenerateOTP()
	if err := h.service.StoreOTP(req.Email, newOTP); err != nil {
		response.InternalError(c, "Failed to resend OTP", gin.H{"error": err.Error()})
		return
	}

	if err := h.service.StoreUserData(req.Email, map[string]interface{}{"refresh": time.Now().Unix()}); err != nil {
		log.Printf("Failed to refresh user data TTL: %v", err)
	}

	if err := h.emailService.SendSignupOTP(req.Email, "User", newOTP, "5 minutes"); err != nil {
		log.Printf("Failed to resend OTP email: %v", err)
	}

	response.Success(c, "OTP resent successfully", gin.H{
		"email":      req.Email,
		"expires_at": time.Now().Add(5 * time.Minute),
	})
}

// ================================================
// LOGIN HANDLER
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
		response.BadRequest(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	user, _, err := h.service.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, "Invalid credentials", gin.H{"error": err.Error()})
		return
	}

	response.Success(c, "2FA verification required", gin.H{
		"requires_2fa": true,
		"email":        user.Email,
		"expires_in":   300,
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
		response.BadRequest(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	user, accessToken, refreshToken, err := h.service.VerifyTwoFactorAndLogin(ctx, req.Email, req.OTP)
	if err != nil {
		response.Unauthorized(c, err.Error(), nil)
		return
	}

	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, refreshToken)

	userInfo := h.service.BuildUserInfo(user)

	response.Success(c, "Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    int64(h.config.JWT.Expiration.Seconds()),
		"user":          userInfo,
	})
}

// ================================================
// REFRESH TOKEN HANDLER
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
		response.Unauthorized(c, "Refresh token required", nil)
		return
	}

	ctx := context.Background()
	newAccessToken, newRefreshToken, err := h.service.RefreshTokens(ctx, refreshToken, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		response.Unauthorized(c, "Invalid refresh token", gin.H{"error": err.Error()})
		return
	}

	h.setAccessTokenCookie(c, newAccessToken)
	h.setRefreshTokenCookie(c, newRefreshToken)

	response.Success(c, "Token refreshed successfully", gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
		"token_type":    "Bearer",
		"expires_in":    int64(h.config.JWT.Expiration.Seconds()),
	})
}

// ================================================
// LOGOUT HANDLER
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
		ctx := context.Background()
		if err := h.service.RevokeToken(ctx, refreshToken); err != nil {
			log.Printf("Failed to revoke refresh token: %v", err)
		}
	}

	h.clearAuthCookies(c)
	response.Success(c, "Logged out successfully", nil)
}

// ================================================
// FORGOT PASSWORD HANDLER
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
		response.BadRequest(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	if valid, msg := validator.Password.Validate(req.NewPassword); !valid {
		response.BadRequest(c, "Weak password", gin.H{
			"field":    "new_password",
			"error":    msg,
			"strength": validator.Password.StrengthLabel(validator.Password.Score(req.NewPassword)),
		})
		return
	}

	ctx := context.Background()
	if err := h.service.InitiatePasswordReset(ctx, req.Email, req.NewPassword); err != nil {
		response.InternalError(c, "Failed to process request", gin.H{"error": err.Error()})
		return
	}

	response.Success(c, "OTP sent to your email", gin.H{
		"message":    "Check your email for the OTP",
		"expires_in": 300,
	})
}

// ================================================
// VERIFY RESET OTP HANDLER
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
		response.BadRequest(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	if err := h.service.VerifyResetOTPAndResetPassword(ctx, req.Email, req.OTP); err != nil {
		response.Unauthorized(c, err.Error(), nil)
		return
	}

	response.Success(c, "Password reset successfully", gin.H{
		"message": "Password reset successfully",
	})
}