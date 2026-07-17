package permissions

import (
	"context"
	"fmt"
	"log"
)

// Service provides high-level permission management
type Service struct {
	enforcer *Enforcer
}

// NewService creates a new permission service
func NewService(enforcer *Enforcer) *Service {
	return &Service{
		enforcer: enforcer,
	}
}

// AssignConsumerRole assigns the consumer role to a user
func (s *Service) AssignConsumerRole(ctx context.Context, userID string) error {
	log.Printf("Assigning consumer role to user: %s", userID)
	_, err := s.enforcer.AddPlatformRole(userID, RoleConsumer)
	if err != nil {
		return fmt.Errorf("failed to assign consumer role: %w", err)
	}
	return nil
}

// AssignBusinessAdminRole assigns the business admin role to a user
func (s *Service) AssignBusinessAdminRole(ctx context.Context, userID string, businessID string) error {
	log.Printf("Assigning business admin role to user: %s for business: %s", userID, businessID)
	_, err := s.enforcer.AddBusinessRole(userID, businessID, RoleBusinessAdmin)
	if err != nil {
		return fmt.Errorf("failed to assign business admin role: %w", err)
	}

	// Add business policies for this business
	err = s.AddBusinessPolicies(ctx, businessID)
	if err != nil {
		return fmt.Errorf("failed to add business policies: %w", err)
	}

	return nil
}

// AssignProductManagerRole assigns the product manager role to a user
func (s *Service) AssignProductManagerRole(ctx context.Context, userID string, businessID string) error {
	log.Printf("Assigning product manager role to user: %s for business: %s", userID, businessID)
	_, err := s.enforcer.AddBusinessRole(userID, businessID, RoleProductManager)
	if err != nil {
		return fmt.Errorf("failed to assign product manager role: %w", err)
	}
	return nil
}

// AssignOrderManagerRole assigns the order manager role to a user
func (s *Service) AssignOrderManagerRole(ctx context.Context, userID string, businessID string) error {
	log.Printf("Assigning order manager role to user: %s for business: %s", userID, businessID)
	_, err := s.enforcer.AddBusinessRole(userID, businessID, RoleOrderManager)
	if err != nil {
		return fmt.Errorf("failed to assign order manager role: %w", err)
	}
	return nil
}

// AssignContentManagerRole assigns the content manager role to a user
func (s *Service) AssignContentManagerRole(ctx context.Context, userID string, businessID string) error {
	log.Printf("Assigning content manager role to user: %s for business: %s", userID, businessID)
	_, err := s.enforcer.AddBusinessRole(userID, businessID, RoleContentManager)
	if err != nil {
		return fmt.Errorf("failed to assign content manager role: %w", err)
	}
	return nil
}

// AssignServiceManagerRole assigns the service manager role to a user
func (s *Service) AssignServiceManagerRole(ctx context.Context, userID string, businessID string) error {
	log.Printf("Assigning service manager role to user: %s for business: %s", userID, businessID)
	_, err := s.enforcer.AddBusinessRole(userID, businessID, RoleServiceManager)
	if err != nil {
		return fmt.Errorf("failed to assign service manager role: %w", err)
	}
	return nil
}

// AssignCustomerSupportRole assigns the customer support role to a user
func (s *Service) AssignCustomerSupportRole(ctx context.Context, userID string, businessID string) error {
	log.Printf("Assigning customer support role to user: %s for business: %s", userID, businessID)
	_, err := s.enforcer.AddBusinessRole(userID, businessID, RoleCustomerSupport)
	if err != nil {
		return fmt.Errorf("failed to assign customer support role: %w", err)
	}
	return nil
}

// RemoveBusinessRole removes a business role from a user
func (s *Service) RemoveBusinessRole(ctx context.Context, userID string, businessID string, role Role) error {
	log.Printf("Removing role %s from user: %s for business: %s", role, userID, businessID)
	_, err := s.enforcer.RemoveBusinessRole(userID, businessID, role)
	if err != nil {
		return fmt.Errorf("failed to remove business role: %w", err)
	}
	return nil
}

// RemoveAllBusinessRoles removes all business roles from a user
func (s *Service) RemoveAllBusinessRoles(ctx context.Context, userID string, businessID string) error {
	log.Printf("Removing all roles for user: %s from business: %s", userID, businessID)
	
	domain := BusinessDomain(businessID)
	roles := s.enforcer.GetRolesForUserInDomain(userID, domain)
	
	for _, role := range roles {
		_, err := s.enforcer.RemoveBusinessRole(userID, businessID, Role(role))
		if err != nil {
			return fmt.Errorf("failed to remove role %s: %w", role, err)
		}
	}
	
	return nil
}

// AddBusinessPolicies adds default policies for a new business
func (s *Service) AddBusinessPolicies(ctx context.Context, businessID string) error {
	domain := BusinessDomain(businessID)
	log.Printf("Adding business policies for business: %s", businessID)

	// Define policies for business admain
	policies := [][]string{ // this data structure is a slice of slices of strings, where each inner slice represents a policy rule
		// Product Management
		{RoleBusinessAdmin.String(), domain, ResourceProduct.String(), ActionCreate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceProduct.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourceProduct.String(), ActionUpdate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceProduct.String(), ActionDelete.String()},

		// Order Management
		{RoleBusinessAdmin.String(), domain, ResourceOrder.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourceOrder.String(), ActionUpdate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceOrder.String(), ActionDelete.String()},

		// Booking Management
		{RoleBusinessAdmin.String(), domain, ResourceBooking.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourceBooking.String(), ActionUpdate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceBooking.String(), ActionDelete.String()},
		{RoleBusinessAdmin.String(), domain, ResourceBooking.String(), "confirm"},
		{RoleBusinessAdmin.String(), domain, ResourceBooking.String(), "complete"},
		{RoleBusinessAdmin.String(), domain, ResourceBooking.String(), "cancel"},

		// Post Management
		{RoleBusinessAdmin.String(), domain, ResourcePost.String(), ActionCreate.String()},
		{RoleBusinessAdmin.String(), domain, ResourcePost.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourcePost.String(), ActionUpdate.String()},
		{RoleBusinessAdmin.String(), domain, ResourcePost.String(), ActionDelete.String()},

		// Chat Management
		{RoleBusinessAdmin.String(), domain, ResourceChat.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourceChat.String(), ActionCreate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceChat.String(), ActionUpdate.String()},

		// Business Profile Management
		{RoleBusinessAdmin.String(), domain, ResourceBusiness.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourceBusiness.String(), ActionUpdate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceBusiness.String(), ActionDelete.String()},

		// Invoice Management
		{RoleBusinessAdmin.String(), domain, ResourceInvoice.String(), ActionCreate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceInvoice.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourceInvoice.String(), ActionUpdate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceInvoice.String(), ActionDelete.String()},
		{RoleBusinessAdmin.String(), domain, ResourceInvoice.String(), "send"},
		{RoleBusinessAdmin.String(), domain, ResourceInvoice.String(), "mark_paid"},

		// Lead Management
		{RoleBusinessAdmin.String(), domain, ResourceLead.String(), ActionCreate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceLead.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourceLead.String(), ActionUpdate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceLead.String(), ActionDelete.String()},
		{RoleBusinessAdmin.String(), domain, ResourceLead.String(), "convert"},

		// Customer Management
		{RoleBusinessAdmin.String(), domain, ResourceCustomer.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourceCustomer.String(), ActionUpdate.String()},
		{RoleBusinessAdmin.String(), domain, ResourceCustomer.String(), ActionDelete.String()},

		// Payment Management
		{RoleBusinessAdmin.String(), domain, ResourcePayment.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourcePayment.String(), "refund"},

		// Dashboard and Analytics
		{RoleBusinessAdmin.String(), domain, ResourceDashboard.String(), ActionRead.String()},
		{RoleBusinessAdmin.String(), domain, ResourceAnalytics.String(), ActionRead.String()},
	}

	// Add all policies using AddPolicies for better performance
	_, err := s.enforcer.AddPolicies(policies)
	if err != nil {
		return fmt.Errorf("failed to add business policies: %w", err)
	}

	log.Printf("Added %d policies for business: %s", len(policies), businessID)
	return nil
}

// RemoveBusinessPolicies removes all policies for a business
func (s *Service) RemoveBusinessPolicies(ctx context.Context, businessID string) error {
	domain := BusinessDomain(businessID)
	log.Printf("Removing policies for business: %s", businessID)

	policies, err := s.enforcer.GetFilteredPolicy(1, domain)
	if err != nil {
		return fmt.Errorf("failed to get filtered policies: %w", err)
	}

	if len(policies) > 0 {
		_, err := s.enforcer.RemovePolicies(policies)
		if err != nil {
			return fmt.Errorf("failed to remove policies: %w", err)
		}
	}

	log.Printf("Removed %d policies for business: %s", len(policies), businessID)
	return nil
}

// CanAccess checks if a user can perform an action on a resource
func (s *Service) CanAccess(ctx context.Context, userID string, domain string, resource Resource, action Action) (bool, error) {
	return s.enforcer.Enforce(userID, domain, resource, action)
}

// CanAccessResource checks if a user can perform an action on a resource (string version)
func (s *Service) CanAccessResource(ctx context.Context, userID string, domain string, resource string, action string) (bool, error) {
	return s.enforcer.EnforceWithContext(userID, domain, resource, action)
}

// GetUserPermissions returns all permissions for a user in a domain
func (s *Service) GetUserPermissions(ctx context.Context, userID string, domain string) ([][]string, error) {
	return s.enforcer.GetImplicitPermissionsForUser(userID, domain)
}

// GetUserRoles returns all roles for a user in a domain
func (s *Service) GetUserRoles(ctx context.Context, userID string, domain string) []string {
	return s.enforcer.GetRolesForUserInDomain(userID, domain)
}

// GetUserImplicitRoles returns all roles including inherited for a user in a domain
func (s *Service) GetUserImplicitRoles(ctx context.Context, userID string, domain string) ([]string, error) {
	return s.enforcer.GetImplicitRolesForUser(userID, domain)
}