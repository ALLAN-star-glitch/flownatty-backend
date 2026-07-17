package permissions

// GetRoleDisplayName returns a human-readable role name for UI
func GetRoleDisplayName(role string) string {
	displayNames := map[string]string{
		RoleSuperAdmin.String():        "Super Admin",
		RoleAdmin.String():             "Admin",
		RoleBusinessAdmin.String():     "Business Admin",
		RoleProductManager.String():    "Product Manager",
		RoleOrderManager.String():      "Order Manager",
		RoleContentManager.String():    "Content Manager",
		RoleServiceManager.String():    "Service Manager",
		RoleCustomerSupport.String():   "Customer Support",
		RolePremiumConsumer.String():   "Premium Consumer",
		RoleConsumer.String():          "Consumer",
		RoleGuest.String():             "Guest",
	}
	if name, ok := displayNames[role]; ok {
		return name
	}
	return role
}

// GetRoleLevel returns hierarchy level for UI ordering
func GetRoleLevel(role string) int {
	levels := map[string]int{
		RoleSuperAdmin.String():        100,
		RoleAdmin.String():             90,
		RoleBusinessAdmin.String():     80,
		RoleProductManager.String():    70,
		RoleOrderManager.String():      70,
		RoleContentManager.String():    70,
		RoleServiceManager.String():    70,
		RoleCustomerSupport.String():   60,
		RolePremiumConsumer.String():   20,
		RoleConsumer.String():          10,
		RoleGuest.String():             0,
	}
	if level, ok := levels[role]; ok {
		return level
	}
	return 0
}

// GetRoleColor returns a color for UI role badges
func GetRoleColor(role string) string {
	colors := map[string]string{
		RoleSuperAdmin.String():        "#EF4444", // Red
		RoleAdmin.String():             "#F59E0B", // Yellow
		RoleBusinessAdmin.String():     "#6366F1", // Indigo
		RoleProductManager.String():    "#3B82F6", // Blue
		RoleOrderManager.String():      "#10B981", // Green
		RoleContentManager.String():    "#8B5CF6", // Purple
		RoleServiceManager.String():    "#EC4899", // Pink
		RoleCustomerSupport.String():   "#F97316", // Orange
		RolePremiumConsumer.String():   "#2DD4BF", // Teal
		RoleConsumer.String():          "#94A3B8", // Light Slate
		RoleGuest.String():             "#CBD5E1", // Light Gray
	}
	if color, ok := colors[role]; ok {
		return color
	}
	return "#64748B"
}

// GetRoleIcon returns an icon name for UI role badges
func GetRoleIcon(role string) string {
	icons := map[string]string{
		RoleSuperAdmin.String():        "shield-check",
		RoleAdmin.String():             "shield",
		RoleBusinessAdmin.String():     "building-office",
		RoleProductManager.String():    "package",
		RoleOrderManager.String():      "truck",
		RoleContentManager.String():    "pencil",
		RoleServiceManager.String():    "wrench",
		RoleCustomerSupport.String():   "chat-bubble-left-right",
		RolePremiumConsumer.String():   "star",
		RoleConsumer.String():          "user",
		RoleGuest.String():             "user",
	}
	if icon, ok := icons[role]; ok {
		return icon
	}
	return "user"
}

// IsManagementRole checks if a role is a management-level role
func IsManagementRole(role string) bool {
	managementRoles := []string{
		RoleBusinessAdmin.String(),
		RoleProductManager.String(),
		RoleOrderManager.String(),
		RoleContentManager.String(),
		RoleServiceManager.String(),
	}
	for _, r := range managementRoles {
		if r == role {
			return true
		}
	}
	return false
}

// IsAdminRole checks if a role is an admin-level role
func IsAdminRole(role string) bool {
	adminRoles := []string{
		RoleSuperAdmin.String(),
		RoleAdmin.String(),
		RoleBusinessAdmin.String(),
	}
	for _, r := range adminRoles {
		if r == role {
			return true
		}
	}
	return false
}