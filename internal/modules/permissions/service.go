package permissions


// Purpose: Business logic layer for permission management.

// What it does:

// Wraps enforcer methods with business logic

// Handles role assignments for different user types

// Adds business-specific policies when a business is created

// Provides clean API for other modules to use

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

// AssignBusinessOwnerRole assigns the business owner role to a user
func (s *Service) AssignBusinessOwnerRole(ctx context.Context, userID string, businessID string) error {
	log.Printf("Assigning business owner role to user: %s for business: %s", userID, businessID)
	_, err := s.enforcer.AddBusinessRole(userID, businessID, RoleBusinessOwner)
	if err != nil {
		return fmt.Errorf("failed to assign business owner role: %w", err)
	}

	// Add business policies for this business
	err = s.AddBusinessPolicies(ctx, businessID)
	if err != nil {
		return fmt.Errorf("failed to add business policies: %w", err)
	}

	return nil
}

// AssignBusinessStaffRole assigns the business staff role to a user
func (s *Service) AssignBusinessStaffRole(ctx context.Context, userID string, businessID string) error {
	log.Printf("Assigning business staff role to user: %s for business: %s", userID, businessID)
	_, err := s.enforcer.AddBusinessRole(userID, businessID, RoleBusinessStaff)
	if err != nil {
		return fmt.Errorf("failed to assign business staff role: %w", err)
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

// AddBusinessPolicies adds default policies for a new business
func (s *Service) AddBusinessPolicies(ctx context.Context, businessID string) error {
	domain := BusinessDomain(businessID)
	log.Printf("Adding business policies for business: %s", businessID)

	// Define policies for business owner
	policies := [][]string{
		// Product management
		{RoleBusinessOwner.String(), domain, ResourceProduct.String(), ActionCreate.String()},
		{RoleBusinessOwner.String(), domain, ResourceProduct.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourceProduct.String(), ActionUpdate.String()},
		{RoleBusinessOwner.String(), domain, ResourceProduct.String(), ActionDelete.String()},

		// Order management
		{RoleBusinessOwner.String(), domain, ResourceOrder.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourceOrder.String(), ActionUpdate.String()},
		{RoleBusinessOwner.String(), domain, ResourceOrder.String(), ActionDelete.String()},

		// Booking management
		{RoleBusinessOwner.String(), domain, ResourceBooking.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourceBooking.String(), ActionUpdate.String()},
		{RoleBusinessOwner.String(), domain, ResourceBooking.String(), ActionDelete.String()},
		{RoleBusinessOwner.String(), domain, ResourceBooking.String(), "confirm"},
		{RoleBusinessOwner.String(), domain, ResourceBooking.String(), "complete"},
		{RoleBusinessOwner.String(), domain, ResourceBooking.String(), "cancel"},

		// Post management
		{RoleBusinessOwner.String(), domain, ResourcePost.String(), ActionCreate.String()},
		{RoleBusinessOwner.String(), domain, ResourcePost.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourcePost.String(), ActionUpdate.String()},
		{RoleBusinessOwner.String(), domain, ResourcePost.String(), ActionDelete.String()},

		// Chat management
		{RoleBusinessOwner.String(), domain, ResourceChat.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourceChat.String(), ActionCreate.String()},
		{RoleBusinessOwner.String(), domain, ResourceChat.String(), ActionUpdate.String()},

		// Business profile
		{RoleBusinessOwner.String(), domain, ResourceBusiness.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourceBusiness.String(), ActionUpdate.String()},
		{RoleBusinessOwner.String(), domain, ResourceBusiness.String(), ActionDelete.String()},

		// Invoice management
		{RoleBusinessOwner.String(), domain, ResourceInvoice.String(), ActionCreate.String()},
		{RoleBusinessOwner.String(), domain, ResourceInvoice.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourceInvoice.String(), ActionUpdate.String()},
		{RoleBusinessOwner.String(), domain, ResourceInvoice.String(), ActionDelete.String()},
		{RoleBusinessOwner.String(), domain, ResourceInvoice.String(), "send"},
		{RoleBusinessOwner.String(), domain, ResourceInvoice.String(), "mark_paid"},

		// Lead management
		{RoleBusinessOwner.String(), domain, ResourceLead.String(), ActionCreate.String()},
		{RoleBusinessOwner.String(), domain, ResourceLead.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourceLead.String(), ActionUpdate.String()},
		{RoleBusinessOwner.String(), domain, ResourceLead.String(), ActionDelete.String()},
		{RoleBusinessOwner.String(), domain, ResourceLead.String(), "convert"},

		// Customer management
		{RoleBusinessOwner.String(), domain, ResourceCustomer.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourceCustomer.String(), ActionUpdate.String()},
		{RoleBusinessOwner.String(), domain, ResourceCustomer.String(), ActionDelete.String()},

		// Payment
		{RoleBusinessOwner.String(), domain, ResourcePayment.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourcePayment.String(), "refund"},

		// Dashboard and analytics
		{RoleBusinessOwner.String(), domain, ResourceDashboard.String(), ActionRead.String()},
		{RoleBusinessOwner.String(), domain, ResourceAnalytics.String(), ActionRead.String()},

		// Staff policies
		{RoleBusinessStaff.String(), domain, ResourceOrder.String(), ActionRead.String()},
		{RoleBusinessStaff.String(), domain, ResourceOrder.String(), ActionUpdate.String()},
		{RoleBusinessStaff.String(), domain, ResourceBooking.String(), ActionRead.String()},
		{RoleBusinessStaff.String(), domain, ResourceBooking.String(), ActionUpdate.String()},
		{RoleBusinessStaff.String(), domain, ResourceChat.String(), ActionRead.String()},
		{RoleBusinessStaff.String(), domain, ResourceChat.String(), ActionCreate.String()},
		{RoleBusinessStaff.String(), domain, ResourceCustomer.String(), ActionRead.String()},
		{RoleBusinessStaff.String(), domain, ResourceProduct.String(), ActionRead.String()},
		{RoleBusinessStaff.String(), domain, ResourcePayment.String(), ActionRead.String()},
		{RoleBusinessStaff.String(), domain, ResourceDashboard.String(), ActionRead.String()},
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

	// Get all policies for this domain
	policies, err := s.enforcer.GetFilteredPolicy(1, domain)
	if err != nil {
		return fmt.Errorf("failed to get filtered policies: %w", err)
	}

	// Remove all policies
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