package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service *AuthService
	repo    *AuthRepository
	config  *config.Config
}

func NewAuthHandler(service *AuthService, repo *AuthRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		repo:    repo,
		config:  cfg,
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Phone    string `json:"phone_number" binding:"required" example:"+254712345678"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Password string `json:"password" binding:"required,min=6" example:"SecurePass123"`
	Role     string `json:"role" binding:"omitempty,oneof=consumer business" example:"consumer"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with phone, email, and password. Sends OTP via email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 200 {object} map[string]interface{} "OTP sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	if req.Role == "" {
		req.Role = "consumer"
	}

	existingUser, _ := h.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already registered",
		})
		return
	}

	existingPhone, _ := h.repo.GetUserByPhone(req.Phone)
	if existingPhone != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone number already registered",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process password",
		})
		return
	}

	otp := h.service.GenerateOTP()

	if h.config.Environment == "development" {
		log.Printf("DEV MODE: OTP for %s is %s", req.Email, otp)
	}

	otpRecord := &models.OTP{
		PhoneNumber:  req.Phone,
		Email:        req.Email,
		OTPCode:      otp,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		Purpose:      "signup",
		ExpiresAt:    time.Now().Add(5 * time.Minute),
		IsUsed:       false,
	}

	if err := h.repo.SaveOTP(otpRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save OTP",
		})
		return
	}

	response := gin.H{
		"message":    "OTP sent successfully",
		"email":      req.Email,
		"phone":      req.Phone,
		"expires_at": time.Now().Add(5 * time.Minute),
	}

	if h.config.Environment == "development" {
		response["otp"] = otp
	}

	c.JSON(http.StatusOK, response)
}

// VerifyOTPRequest represents the OTP verification request body
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// VerifyOTP godoc
// @Summary Verify OTP and create account
// @Description Verify the OTP sent to the user's email and create the account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "OTP verification details"
// @Success 201 {object} map[string]interface{} "Account created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid OTP"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	otpRecord, err := h.repo.GetOTPByEmail(req.Email, req.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	otpRecord.IsUsed = true
	if err := h.repo.UpdateOTP(otpRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update OTP",
		})
		return
	}

	now := time.Now()
	user := &models.User{
		PhoneNumber:     otpRecord.PhoneNumber,
		Email:           otpRecord.Email,
		Password:        otpRecord.PasswordHash,
		Name:            otpRecord.Name,
		Role:            otpRecord.Role,
		IsVerified:      true,
		IsEmailVerified: false,
		VerifiedAt:      &now,
	}

	if err := h.repo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	token, err := h.service.GenerateToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Account created successfully",
		"access_token": token,
		"user": gin.H{
			"id":    user.ID,
			"phone": user.PhoneNumber,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

// ResendOTPRequest represents the resend OTP request body
type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// ResendOTP godoc
// @Summary Resend OTP
// @Description Resend OTP to the user's email for registration
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ResendOTPRequest true "Email to resend OTP"
// @Success 200 {object} map[string]interface{} "OTP resent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req ResendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	existingUser, _ := h.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already registered",
		})
		return
	}

	otpRecord, err := h.repo.GetLatestOTP(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No pending registration found",
		})
		return
	}

	newOTP := h.service.GenerateOTP()

	otpRecord.OTPCode = newOTP
	otpRecord.ExpiresAt = time.Now().Add(5 * time.Minute)
	otpRecord.IsUsed = false

	if err := h.repo.UpdateOTP(otpRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to resend OTP",
		})
		return
	}

	if h.config.Environment == "development" {
		log.Printf("DEV MODE: New OTP for %s is %s", req.Email, newOTP)
	}

	response := gin.H{
		"message":    "OTP resent successfully",
		"email":      req.Email,
		"expires_at": time.Now().Add(5 * time.Minute),
	}

	if h.config.Environment == "development" {
		response["otp"] = newOTP
	}

	c.JSON(http.StatusOK, response)
}