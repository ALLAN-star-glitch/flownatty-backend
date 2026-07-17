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

	// Business roles (platform permission level)
	RoleBusinessAdmin    Role = "business_admin"     // Full business management
	RoleProductManager   Role = "product_manager"    // Manage products, inventory, categories
	RoleOrderManager     Role = "order_manager"      // Manage orders, shipping, refunds
	RoleContentManager   Role = "content_manager"    // Manage posts, promotions, announcements
	RoleServiceManager   Role = "service_manager"    // Manage services, bookings, appointments
	RoleCustomerSupport  Role = "customer_support"   // Customer service, chat, support tickets

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
	ResourcePromotion Resource = "promotion"
	ResourceCategory  Resource = "category"
	ResourceMember    Resource = "member" 
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
func BusinessDomain(businessID string) string {
	if businessID == "" {
		return ""
	}
	return "business:" + businessID
}

func UserDomain(userID string) string {
	if userID == "" {
		return ""
	}
	return "user:" + userID
}

func IsBusinessDomain(domain string) bool {
	return len(domain) > 9 && domain[:9] == "business:"
}

func IsUserDomain(domain string) bool {
	return len(domain) > 5 && domain[:5] == "user:"
}

func IsPlatformDomain(domain string) bool {
	return domain == DomainPlatform
}

func ExtractBusinessID(domain string) string {
	if IsBusinessDomain(domain) {
		return domain[9:]
	}
	return ""
}

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
		RoleBusinessAdmin,
		RoleProductManager,
		RoleOrderManager,
		RoleContentManager,
		RoleServiceManager,
		RoleCustomerSupport,
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
		ResourcePromotion,
		ResourceCategory,
		ResourceMember,
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