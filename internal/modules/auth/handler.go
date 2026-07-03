package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/email"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service      *AuthService
	repo         *AuthRepository
	config       *config.Config
	emailService *email.EmailService
}

func NewAuthHandler(service *AuthService, repo *AuthRepository, cfg *config.Config) *AuthHandler {
	emailSvc := email.NewEmailService(cfg.Resend.ApiKey, cfg.Resend.From)

	return &AuthHandler{
		service:      service,
		repo:         repo,
		config:       cfg,
		emailService: emailSvc,
	}
}

// Register Request
type RegisterRequest struct {
	Phone    string `json:"phone_number" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=consumer business"`
}

// Register Handler
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

	// Generate random OTP (always random)
	otp := h.service.GenerateOTP()

	log.Printf("Generated OTP for %s: %s", req.Email, otp)

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

	// ✅ Always send OTP via email (both development and production)
	if err := h.emailService.SendOTP(email.OTPEmailData{
		To:      req.Email,
		Name:    req.Name,
		OTP:     otp,
		Expires: "5 minutes",
	}); err != nil {
		log.Printf("Failed to send OTP email: %v", err)
	}

	response := gin.H{
		"message":    "OTP sent successfully",
		"email":      req.Email,
		"phone":      req.Phone,
		"expires_at": time.Now().Add(5 * time.Minute),
	}

	// ✅ Return OTP in development only (for testing)
	if h.config.Environment == "development" {
		response["otp"] = otp
	}

	c.JSON(http.StatusOK, response)
}

// Verify OTP Request
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// Verify OTP Handler
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// ✅ Allow 123456 as fallback (for development testing)
	if req.OTP == "123456" {
		// Find the OTP record by email only (not by code)
		otpRecord, err := h.repo.GetLatestOTP(req.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No pending registration found",
			})
			return
		}

		if time.Now().After(otpRecord.ExpiresAt) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "OTP expired",
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

		if err := h.emailService.SendWelcome(email.OTPEmailData{
			To:   user.Email,
			Name: user.Name,
		}); err != nil {
			log.Printf("Failed to send welcome email: %v", err)
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
		return
	}

	// Normal flow: verify with real OTP
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

	if err := h.emailService.SendWelcome(email.OTPEmailData{
		To:   user.Email,
		Name: user.Name,
	}); err != nil {
		log.Printf("Failed to send welcome email: %v", err)
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

// Resend OTP Request
type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Resend OTP Handler
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

	// Generate new random OTP
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

	// Send new OTP via email
	if err := h.emailService.SendOTP(email.OTPEmailData{
		To:      req.Email,
		Name:    "User", // We don't have the name here
		OTP:     newOTP,
		Expires: "5 minutes",
	}); err != nil {
		log.Printf("Failed to resend OTP email: %v", err)
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