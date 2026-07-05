package middleware

import (
	"strings"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/auth/permissions"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT from cookie or Authorization header
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Try to get token from cookie (preferred - more secure)
		tokenCookie, err := c.Cookie("access_token")
		if err == nil && tokenCookie != "" {
			tokenString = tokenCookie
		} else {
			// 2. Fallback: Get token from Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		// 3. No token found
		if tokenString == "" {
			response.Unauthorized(c, "Authentication required", gin.H{
				"reason": "no access token found in cookie or Authorization header",
			})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "Invalid or expired token", gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c, "Invalid token claims", nil)
			c.Abort()
			return
		}

		// Check token type
		tokenType, _ := claims["type"].(string)
		if tokenType != "access" {
			response.Unauthorized(c, "Invalid token type", gin.H{
				"expected": "access",
				"got":      tokenType,
			})
			c.Abort()
			return
		}

		// Get user ID from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			response.Unauthorized(c, "User ID not found in token", nil)
			c.Abort()
			return
		}

		// Get user role if present
		userRole, _ := claims["role"].(string)

		// Set user context using permissions package constants
		c.Set(permissions.ContextKeyUserID, userID)
		c.Set(permissions.ContextKeyUserRole, userRole)

		c.Next()
	}
}