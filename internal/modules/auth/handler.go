package auth

import (
	"log"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/email"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
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
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Role == "" {
		req.Role = "consumer"
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

	// Generate random OTP
	otp := h.service.GenerateOTP()

	log.Printf("Generated OTP for %s: %s", req.Email, otp)

	// Store OTP in Redis (primary)
	if err := h.service.StoreOTP(req.Email, otp); err != nil {
		log.Printf("Failed to store OTP in Redis: %v", err)
		response.InternalError(c, "Failed to process request", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Store user data in Redis (with same TTL as OTP)
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

	// Enqueue OTP email via Asynq (non-blocking)
	if err := h.service.EnqueueOTPEmail(req.Email, req.Name, otp); err != nil {
		log.Printf("Failed to enqueue OTP email: %v", err)
		// Fallback: send synchronously
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
		"expires_at": time.Now().Add(5 * time.Minute),
	}

	response.Success(c, "OTP sent successfully", respData)
}

// Verify OTP Request
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
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

	// Check if user already exists by phone
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

	token, err := h.service.GenerateToken(user.ID.String())
	if err != nil {
		response.InternalError(c, "Failed to generate token", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Enqueue Welcome Email
	if err := h.service.EnqueueWelcomeEmail(user.Email, user.Name); err != nil {
		log.Printf("Failed to enqueue welcome email: %v", err)
		// Fallback: send synchronously
		if err := h.emailService.SendWelcome(email.WelcomeEmailData{
			To:   user.Email,
			Name: user.Name,
		}); err != nil {
			log.Printf("Failed to send welcome email: %v", err)
		}
	}

	response.Created(c, "Account created successfully", gin.H{
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

// Verify OTP Handler
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("Verifying OTP for %s", req.Email)

	// Check Redis for OTP
	storedOTP, err := h.service.GetOTP(req.Email)

	if err != nil {
		// OTP not found in Redis - check if 123456 fallback is allowed
		// In production, 123456 is NOT a valid OTP
		response.BadRequest(c, "Invalid or expired OTP", gin.H{
			"email": req.Email,
		})
		return
	}

	// Verify OTP
	if req.OTP != storedOTP {
		response.BadRequest(c, "Invalid OTP", gin.H{
			"email":         req.Email,
			"provided_otp": req.OTP,
		})
		return
	}

	log.Printf("OTP verified for %s", req.Email)

	// Get user data from Redis
	userData, err := h.service.GetUserData(req.Email)
	if err != nil {
		response.BadRequest(c, "Registration data not found. Please register again.", gin.H{
			"email": req.Email,
		})
		return
	}

	// Delete OTP and user data from Redis
	h.service.DeleteOTP(req.Email)
	h.service.DeleteUserData(req.Email)

	// Create user
	h.createUserFromData(c, req.Email, userData)
}

// Resend OTP Request
type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Resend OTP Handler
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

	// Generate new random OTP
	newOTP := h.service.GenerateOTP()

	// Store new OTP in Redis
	if err := h.service.StoreOTP(req.Email, newOTP); err != nil {
		log.Printf("Failed to store OTP in Redis: %v", err)
		response.InternalError(c, "Failed to resend OTP", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Refresh user data TTL in Redis
	if err := h.service.StoreUserData(req.Email, map[string]interface{}{
		"refresh": time.Now().Unix(),
	}); err != nil {
		log.Printf("Failed to refresh user data TTL: %v", err)
	}

	log.Printf("New OTP generated for %s: %s", req.Email, newOTP)

	// Send new OTP via email
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