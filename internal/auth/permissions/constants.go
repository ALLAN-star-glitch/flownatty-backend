package permissions

// Role represents a user role in the system
type Role string

const (
	// Platform-level roles
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"

	// Consumer roles
	RoleConsumer      Role = "consumer"
	RolePremiumConsumer Role = "premium_consumer"

	// Business roles (domain-specific)
	RoleBusinessOwner Role = "business_owner"
	RoleBusinessStaff Role = "business_staff"

	// System
	RoleGuest Role = "guest"
)

// Resource represents a resource/object being accessed
type Resource string

const (
	ResourceProduct   Resource = "product"
	ResourceOrder     Resource = "order"
	ResourceBooking   Resource = "booking"
	ResourcePost      Resource = "post"
	ResourceBusiness  Resource = "business"
	ResourceUser      Resource = "user"
	ResourceChat      Resource = "chat"
	ResourceInvoice   Resource = "invoice"
	ResourceLead      Resource = "lead"
	ResourcePlatform  Resource = "platform"
	ResourceAnalytics Resource = "analytics"
	ResourceService   Resource = "service"
	ResourceCustomer  Resource = "customer"
	ResourcePayment   Resource = "payment"
	ResourceCart      Resource = "cart"
	ResourceWishlist  Resource = "wishlist"
	ResourceFollow    Resource = "follow"
	ResourceNotification Resource = "notification"
	ResourceDashboard Resource = "dashboard"
)

// Action represents an operation that can be performed
type Action string

const (
	ActionCreate   Action = "create"
	ActionRead     Action = "read"
	ActionUpdate   Action = "update"
	ActionDelete   Action = "delete"
	ActionManage   Action = "manage" // Full CRUD
	ActionExport   Action = "export"
	ActionConfirm  Action = "confirm"
	ActionComplete Action = "complete"
	ActionCancel   Action = "cancel"
	ActionSend     Action = "send"
	ActionMarkPaid Action = "mark_paid"
	ActionConvert  Action = "convert"
	ActionRefund   Action = "refund"
	ActionLike     Action = "like"
	ActionComment  Action = "comment"
)

// Domain constants
const (
	DomainPlatform = "platform"
)

// Domain helpers

// BusinessDomain returns the domain string for a business
func BusinessDomain(businessID string) string {
	if businessID == "" {
		return ""
	}
	return "business:" + businessID
}

// UserDomain returns the domain string for a user
func UserDomain(userID string) string {
	if userID == "" {
		return ""
	}
	return "user:" + userID
}

// IsBusinessDomain checks if a domain is a business domain
func IsBusinessDomain(domain string) bool {
	return len(domain) > 9 && domain[:9] == "business:"
}

// IsUserDomain checks if a domain is a user domain
func IsUserDomain(domain string) bool {
	return len(domain) > 5 && domain[:5] == "user:"
}

// IsPlatformDomain checks if a domain is the platform domain
func IsPlatformDomain(domain string) bool {
	return domain == DomainPlatform
}

// ExtractBusinessID extracts the business ID from a domain
func ExtractBusinessID(domain string) string {
	if IsBusinessDomain(domain) {
		return domain[9:]
	}
	return ""
}

// ExtractUserID extracts the user ID from a domain
func ExtractUserID(domain string) string {
	if IsUserDomain(domain) {
		return domain[5:]
	}
	return ""
}

// GetAllRoles returns all defined roles
func GetAllRoles() []Role {
	return []Role{
		RoleSuperAdmin,
		RoleAdmin,
		RoleConsumer,
		RolePremiumConsumer,
		RoleBusinessOwner,
		RoleBusinessStaff,
		RoleGuest,
	}
}

// GetAllResources returns all defined resources
func GetAllResources() []Resource {
	return []Resource{
		ResourceProduct,
		ResourceOrder,
		ResourceBooking,
		ResourcePost,
		ResourceBusiness,
		ResourceUser,
		ResourceChat,
		ResourceInvoice,
		ResourceLead,
		ResourcePlatform,
		ResourceAnalytics,
		ResourceService,
		ResourceCustomer,
		ResourcePayment,
		ResourceCart,
		ResourceWishlist,
		ResourceFollow,
		ResourceNotification,
		ResourceDashboard,
	}
}

// GetAllActions returns all defined actions
func GetAllActions() []Action {
	return []Action{
		ActionCreate,
		ActionRead,
		ActionUpdate,
		ActionDelete,
		ActionManage,
		ActionExport,
		ActionConfirm,
		ActionComplete,
		ActionCancel,
		ActionSend,
		ActionMarkPaid,
		ActionConvert,
		ActionRefund,
		ActionLike,
		ActionComment,
	}
}

// String returns the string representation of a Role
func (r Role) String() string {
	return string(r)
}

// String returns the string representation of a Resource
func (r Resource) String() string {
	return string(r)
}

// String returns the string representation of an Action
func (a Action) String() string {
	return string(a)
}