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

    // Set default role
    if req.Role == "" {
        req.Role = "consumer"
    }

    // Check if email already exists
    existingUser, _ := h.repo.GetUserByEmail(req.Email)
    if existingUser != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Email already registered",
        })
        return
    }

    // Check if phone already exists
    existingPhone, _ := h.repo.GetUserByPhone(req.Phone)
    if existingPhone != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Phone number already registered",
        })
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to process password",
        })
        return
    }

    // Generate OTP
    otp := h.service.GenerateOTP()

    // Log OTP in development
    if h.config.Environment == "development" {
        log.Printf("DEV MODE: OTP for %s is %s", req.Email, otp)
    }

    // Store OTP in database
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

    // TODO: Send OTP via email (will implement later)
    // For now, we just log it

    response := gin.H{
        "message":    "OTP sent successfully",
        "email":      req.Email,
        "phone":      req.Phone,
        "expires_at": time.Now().Add(5 * time.Minute),
    }

    // Return OTP in development only
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

    // Get OTP record
    otpRecord, err := h.repo.GetOTPByEmail(req.Email, req.OTP)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // Mark OTP as used
    otpRecord.IsUsed = true
    if err := h.repo.UpdateOTP(otpRecord); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to update OTP",
        })
        return
    }

    // Create user
    user := &models.User{
        PhoneNumber:     otpRecord.PhoneNumber,
        Email:           otpRecord.Email,
        Password:        otpRecord.PasswordHash,
        Name:            otpRecord.Name,
        Role:            otpRecord.Role,
        IsVerified:      true,
        IsEmailVerified: false,
        VerifiedAt:      &[]time.Time{time.Now()}[0],
    }

    if err := h.repo.CreateUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create user",
        })
        return
    }

    // Generate JWT token
    token, err := h.service.GenerateToken(user.ID.String())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to generate token",
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Account created successfully",
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

    // Check if user already exists
    existingUser, _ := h.repo.GetUserByEmail(req.Email)
    if existingUser != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Email already registered",
        })
        return
    }

    // Get existing OTP
    otpRecord, err := h.repo.GetLatestOTP(req.Email)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "No pending registration found",
        })
        return
    }

    // Generate new OTP
    newOTP := h.service.GenerateOTP()

    // Update OTP record
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