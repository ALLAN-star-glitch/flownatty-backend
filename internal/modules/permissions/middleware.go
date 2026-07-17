package permissions

import (
	"net/http"
	"strings"

	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// Context keys for storing user information
const (
	ContextKeyUserID     = "user_id"
	ContextKeyUserRole   = "user_role"
	ContextKeyDomain     = "domain"
	ContextKeyBusinessID = "business_id"
	ContextKeyUserRoles  = "user_roles"
)

// AuthorizationMiddleware creates a Gin middleware for authorization
func AuthorizationMiddleware(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			response.Unauthorized(c, "User not authenticated", gin.H{
				"reason": "user_id not found in context",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			response.Unauthorized(c, "Invalid user ID", gin.H{
				"reason": "user_id is not a string",
			})
			c.Abort()
			return
		}

		// Determine domain from request
		domain := getDomainFromRequest(c)

		// Determine resource from request path
		resource := getResourceFromRequest(c)

		// Determine action from HTTP method
		action := getActionFromRequest(c)

		// Get user's roles for this domain for context
		roles := enforcer.GetRolesForUserInDomain(userIDStr, domain)
		c.Set(ContextKeyUserRoles, roles)

		// If domain is a business domain, store business ID
		if IsBusinessDomain(domain) {
			c.Set(ContextKeyBusinessID, ExtractBusinessID(domain))
		}

		// Enforce permission
		allowed, err := enforcer.EnforceWithContext(userIDStr, domain, resource, action)
		if err != nil {
			response.InternalError(c, "Authorization error", gin.H{
				"error":    err.Error(),
				"user":     userIDStr,
				"domain":   domain,
				"resource": resource,
				"action":   action,
			})
			c.Abort()
			return
		}

		if !allowed {
			response.Forbidden(c, "Insufficient permissions", gin.H{
				"user":     userIDStr,
				"domain":   domain,
				"resource": resource,
				"action":   action,
				"roles":    roles,
			})
			c.Abort()
			return
		}

		// Store domain in context for downstream handlers
		c.Set(ContextKeyDomain, domain)
		c.Next()
	}
}

// RequireRoles creates middleware that requires specific roles
func RequireRoles(enforcer *Enforcer, roles ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			response.Unauthorized(c, "User not authenticated", gin.H{
				"reason": "user_id not found in context",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			response.Unauthorized(c, "Invalid user ID", gin.H{
				"reason": "user_id is not a string",
			})
			c.Abort()
			return
		}

		// Get domain from context
		domain, exists := c.Get(ContextKeyDomain)
		if !exists {
			domain = DomainPlatform
		}
		domainStr, ok := domain.(string)
		if !ok {
			response.InternalError(c, "Invalid domain type", gin.H{
				"reason": "domain is not a string",
			})
			c.Abort()
			return
		}

		// Get user's roles in this domain
		userRoles := enforcer.GetRolesForUserInDomain(userIDStr, domainStr)

		// Check if user has any of the required roles
		roleMap := make(map[string]bool)
		for _, r := range userRoles {
			roleMap[r] = true
		}

		// Check for each required role
		for _, required := range roles {
			if roleMap[required.String()] {
				c.Next()
				return
			}
		}

		// Build list of required role names for error message
		requiredRoleNames := make([]string, len(roles))
		for i, r := range roles {
			requiredRoleNames[i] = r.String()
		}

		response.Forbidden(c, "Insufficient roles", gin.H{
			"user":           userIDStr,
			"domain":         domainStr,
			"required_roles": requiredRoleNames,
			"user_roles":     userRoles,
		})
		c.Abort()
	}
}

// RequireBusinessAdmin creates middleware to check if user is a business admin
// This replaces the old RequireBusinessOwner middleware
func RequireBusinessAdmin(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			response.Unauthorized(c, "User not authenticated", gin.H{
				"reason": "user_id not found in context",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			response.Unauthorized(c, "Invalid user ID", gin.H{
				"reason": "user_id is not a string",
			})
			c.Abort()
			return
		}

		// Get business ID from path
		businessID := c.Param("businessId")
		if businessID == "" {
			businessID = c.Param("id") // Fallback
		}

		if businessID == "" {
			// Try to get from query param
			businessID = c.Query("businessId")
		}

		if businessID == "" {
			response.BadRequest(c, "Business ID required", gin.H{
				"reason": "businessId not found in path or query",
			})
			c.Abort()
			return
		}

		// Check if user has business_admin role
		if !enforcer.IsBusinessAdmin(userIDStr, businessID) {
			response.Forbidden(c, "Not a business admin", gin.H{
				"user":          userIDStr,
				"business_id":   businessID,
				"required_role": RoleBusinessAdmin.String(),
			})
			c.Abort()
			return
		}

		c.Set(ContextKeyBusinessID, businessID)
		c.Next()
	}
}

// RequireBusinessRole creates middleware to check if user has ANY business role
func RequireBusinessRole(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			response.Unauthorized(c, "User not authenticated", gin.H{
				"reason": "user_id not found in context",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			response.Unauthorized(c, "Invalid user ID", gin.H{
				"reason": "user_id is not a string",
			})
			c.Abort()
			return
		}

		// Get business ID from path
		businessID := c.Param("businessId")
		if businessID == "" {
			businessID = c.Param("id")
		}

		if businessID == "" {
			businessID = c.Query("businessId")
		}

		if businessID == "" {
			response.BadRequest(c, "Business ID required", gin.H{
				"reason": "businessId not found in path or query",
			})
			c.Abort()
			return
		}

		domain := BusinessDomain(businessID)
		roles := enforcer.GetRolesForUserInDomain(userIDStr, domain)

		// Check if user has any business role
		hasAccess := false
		for _, role := range roles {
			switch role {
			case RoleBusinessAdmin.String(),
				RoleProductManager.String(),
				RoleOrderManager.String(),
				RoleContentManager.String(),
				RoleServiceManager.String(),
				RoleCustomerSupport.String():
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			response.Forbidden(c, "No business access", gin.H{
				"user":        userIDStr,
				"business_id": businessID,
				"user_roles":  roles,
			})
			c.Abort()
			return
		}

		c.Set(ContextKeyBusinessID, businessID)
		c.Set(ContextKeyDomain, domain)
		c.Next()
	}
}

// RequirePlatformRole creates middleware that requires a platform-level role
func RequirePlatformRole(enforcer *Enforcer, roles ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			response.Unauthorized(c, "User not authenticated", gin.H{
				"reason": "user_id not found in context",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			response.Unauthorized(c, "Invalid user ID", gin.H{
				"reason": "user_id is not a string",
			})
			c.Abort()
			return
		}

		// Get user's platform roles
		userRoles := enforcer.GetUserPlatformRoles(userIDStr)

		// Check if user has any of the required roles
		roleMap := make(map[string]bool)
		for _, r := range userRoles {
			roleMap[r] = true
		}

		for _, required := range roles {
			if roleMap[required.String()] {
				c.Next()
				return
			}
		}

		requiredRoleNames := make([]string, len(roles))
		for i, r := range roles {
			requiredRoleNames[i] = r.String()
		}

		response.Forbidden(c, "Insufficient platform roles", gin.H{
			"user":           userIDStr,
			"required_roles": requiredRoleNames,
			"user_roles":     userRoles,
		})
		c.Abort()
	}
}

// RequireBusinessAccess creates middleware that checks if user has ANY business role
// Use this for /businesses/me and /businesses/my endpoints
func RequireBusinessAccess(enforcer *Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			response.Unauthorized(c, "User not authenticated", gin.H{
				"reason": "user_id not found in context",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			response.Unauthorized(c, "Invalid user ID", gin.H{
				"reason": "user_id is not a string",
			})
			c.Abort()
			return
		}

		// Check if user has ANY business role
		hasBusinessAccess := enforcer.HasAnyBusinessRole(userIDStr)

		if !hasBusinessAccess {
			// Get user's businesses for debugging
			businesses := enforcer.GetUserBusinesses(userIDStr)
			platformRoles := enforcer.GetUserPlatformRoles(userIDStr)

			response.Forbidden(c, "User does not belong to any business", gin.H{
				"user":           userIDStr,
				"businesses":     businesses,
				"platform_roles": platformRoles,
				"message":        "User must be a member of at least one business to access this endpoint",
			})
			c.Abort()
			return
		}

		// Store user ID in context for downstream handlers
		c.Set(ContextKeyUserID, userIDStr)

		// Get user's businesses and store the first one as default if needed
		businesses := enforcer.GetUserBusinesses(userIDStr)
		if len(businesses) > 0 {
			c.Set(ContextKeyBusinessID, businesses[0])
		}

		c.Next()
	}
}

// getDomainFromRequest extracts the domain from the request


func getDomainFromRequest(c *gin.Context) string {
	// Check if domain is in context (set by previous middleware)
	if domain, exists := c.Get(ContextKeyDomain); exists {
		if domainStr, ok := domain.(string); ok {
			return domainStr
		}
	}

	path := c.Request.URL.Path

	// ✅ Check business routes FIRST (before user routes)
	if strings.Contains(path, "/businesses/me") ||
		strings.Contains(path, "/business/me") ||
		strings.Contains(path, "/businesses/my") ||
		strings.Contains(path, "/business/my") {
		return DomainPlatform
	}

	// ✅ Check for business ID in path - check BOTH "id" and "businessId"
	businessID := c.Param("businessId")
	if businessID == "" {
		businessID = c.Param("id") // ✅ Also check "id" parameter
	}
	if businessID != "" {
		return BusinessDomain(businessID)
	}

	// Check if business ID is in query
	businessID = c.Query("businessId")
	if businessID != "" {
		return BusinessDomain(businessID)
	}

	// Check if it's a user-specific request
	userID := c.Param("userId")
	if userID != "" {
		return UserDomain(userID)
	}

	// Check if it's the current user profile
	if strings.Contains(path, "/profile") && !strings.Contains(path, "/businesses/") {
		if userID, exists := c.Get(ContextKeyUserID); exists {
			if userIDStr, ok := userID.(string); ok {
				return UserDomain(userIDStr)
			}
		}
	}

	// Check if it's a me endpoint (but NOT business me - already handled above)
	if strings.Contains(path, "/me") && !strings.Contains(path, "/businesses/") {
		if userID, exists := c.Get(ContextKeyUserID); exists {
			if userIDStr, ok := userID.(string); ok {
				return UserDomain(userIDStr)
			}
		}
	}

	// Default to platform domain
	return DomainPlatform
}

// getResourceFromRequest extracts the resource from the request path


func getResourceFromRequest(c *gin.Context) string {
	path := c.Request.URL.Path

	// Remove API prefix
	path = strings.TrimPrefix(path, "/api/v1/")
	path = strings.TrimPrefix(path, "/api/v1/business/")
	path = strings.TrimPrefix(path, "/api/v1/consumer/")

	// Get the first segment
	segments := strings.Split(path, "/")
	if len(segments) > 0 && segments[0] != "" {
		resource := segments[0]

		// Handle business "me" endpoints
		if resource == "businesses" || resource == "business" {
			if len(segments) > 1 && (segments[1] == "me" || segments[1] == "my") {
				return ResourceBusiness.String()
			}
			// ✅ Handle member endpoints - check for "members" in any position
			for _, segment := range segments {
				if segment == "members" {
					return ResourceMember.String()
				}
			}
		}

		// Map common URL patterns to resources
		switch resource {
		case "profile", "me":
			return ResourceUser.String()
		case "businesses", "business":
			return ResourceBusiness.String()
		case "products":
			return ResourceProduct.String()
		case "orders":
			return ResourceOrder.String()
		case "bookings":
			return ResourceBooking.String()
		case "posts":
			return ResourcePost.String()
		case "chat", "chats":
			return ResourceChat.String()
		case "invoices":
			return ResourceInvoice.String()
		case "leads":
			return ResourceLead.String()
		case "customers":
			return ResourceCustomer.String()
		case "payments":
			return ResourcePayment.String()
		case "cart":
			return ResourceCart.String()
		case "wishlist":
			return ResourceWishlist.String()
		case "follow", "following":
			return ResourceFollow.String()
		case "notifications":
			return ResourceNotification.String()
		case "dashboard":
			return ResourceDashboard.String()
		case "analytics":
			return ResourceAnalytics.String()
		case "members":
			return ResourceMember.String()
		default:
			return resource
		}
	}

	return ""
}

// getActionFromRequest maps HTTP method to action
func getActionFromRequest(c *gin.Context) string {
	method := c.Request.Method

	switch method {
	case http.MethodGet:
		return ActionRead.String()
	case http.MethodPost:
		return ActionCreate.String()
	case http.MethodPut, http.MethodPatch:
		return ActionUpdate.String()
	case http.MethodDelete:
		return ActionDelete.String()
	default:
		return ActionRead.String()
	}
}