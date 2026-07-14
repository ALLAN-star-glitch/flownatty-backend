package authentication

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/handler"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.RouterGroup,
	authHandler *handler.AuthHandler,
) {
	auth := router.Group("/auth")
	{
		// ================================================
		// REGISTRATION FLOW (Public)
		// ================================================
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/resend-otp", authHandler.ResendOTP)

		// ================================================
		// BUSINESS EMAIL VERIFICATION (Public)
		// ================================================
		auth.POST("/verify-business-email", authHandler.VerifyBusinessEmail)
		auth.POST("/resend-business-otp", authHandler.ResendBusinessOTP)

		// ================================================
		// AUTHENTICATION FLOW (Public)
		// ================================================
		auth.POST("/login", authHandler.Login)
		auth.POST("/verify-2fa", authHandler.VerifyTwoFactorOTP)

		// ================================================
		// TOKEN MANAGEMENT (Public - Cookie Based)
		// ================================================
		auth.POST("/refresh", authHandler.RefreshToken)

		// ================================================
		// SESSION MANAGEMENT (Cookie Based)
		// ================================================
		auth.POST("/logout", authHandler.Logout)

		// ================================================
		// PASSWORD RESET FLOW (Public - OTP Based)
		// ================================================
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/verify-reset-otp", authHandler.VerifyResetOTP)
	}
}