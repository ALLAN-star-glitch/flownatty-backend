package permissions

// Purpose: The core Casbin engine.

// The enforcer package provides a singleton Casbin enforcer with thread-safe methods for permission management, including auto-reloading policies from the database. It supports role-based access control (RBAC) with domain support, allowing for flexible permission checks across different business contexts.

// Initializes Casbin with model.conf and database adapter

// Loads policies from database (casbin_rule table)

// Provides methods to check permissions

// Provides methods to manage roles and policies

// Auto-reloads policies if enabled

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"gorm.io/gorm"
)

var (
	instance *Enforcer
	once     sync.Once
	initErr  error
)

// Enforcer wraps the Casbin enforcer with additional functionality - singleton pattern, auto-reload, and thread safety
type Enforcer struct {
	*casbin.Enforcer
	mu      sync.RWMutex
	cfg     *config.CasbinConfig
	db      *gorm.DB
	ctx     context.Context
	cancel  context.CancelFunc
	stopped bool
}

// InitEnforcer initializes the Casbin enforcer (singleton)
func InitEnforcer(db *gorm.DB, cfg *config.Config) (*Enforcer, error) {
	once.Do(func() {
		instance, initErr = newEnforcer(db, cfg)
	})

	if initErr != nil {
		return nil, initErr
	}
	return instance, nil
}

func newEnforcer(db *gorm.DB, cfg *config.Config) (*Enforcer, error) {
	// Create GORM adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin adapter: %w", err)
	}

	// Create enforcer with model and adapter
	e, err := casbin.NewEnforcer(cfg.Casbin.ModelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin enforcer: %w", err)
	}

	// Enable auto-save for policy changes
	e.EnableAutoSave(true)

	// Load policies from database
	err = e.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to load policies: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	enforcer := &Enforcer{
		Enforcer: e,
		cfg:      &cfg.Casbin,
		db:       db,
		ctx:      ctx,
		cancel:   cancel,
	}

	log.Println("Casbin enforcer initialized successfully")

	// Start auto-reload if enabled
	if cfg.Casbin.AutoLoad {
		go enforcer.autoLoadPolicies()
	}

	return enforcer, nil
}

// GetEnforcer returns the singleton enforcer instance
func GetEnforcer() *Enforcer {
	if instance == nil {
		panic("Casbin enforcer not initialized")
	}
	return instance
}

// Close stops the enforcer and cleans up resources
func (e *Enforcer) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.stopped {
		e.cancel()
		e.stopped = true
	}
}

// autoLoadPolicies periodically reloads policies from the database
func (e *Enforcer) autoLoadPolicies() {
	ticker := time.NewTicker(e.cfg.AutoLoadInterval)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			log.Println("Auto-load stopped")
			return
		case <-ticker.C:
			e.mu.Lock()
			err := e.LoadPolicy()
			if err != nil {
				log.Printf("Failed to auto-load policies: %v", err)
			} else {
				log.Println("Policies auto-loaded successfully")
			}
			e.mu.Unlock()
		}
	}
}

// Enforce checks if a user has permission in a domain
func (e *Enforcer) Enforce(userID string, domain string, resource Resource, action Action) (bool, error) {
	return e.EnforceWithContext(userID, domain, resource.String(), action.String())
}

// EnforceWithContext checks permission with explicit resource and action strings
func (e *Enforcer) EnforceWithContext(userID string, domain string, resource string, action string) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.Enforce(userID, domain, resource, action)
}

// BatchEnforce checks multiple permissions in one call
func (e *Enforcer) BatchEnforce(requests [][]interface{}) ([]bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.BatchEnforce(requests)
}

// AddPolicy adds a new policy rule
func (e *Enforcer) AddPolicy(sub, dom, obj, act string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.Enforcer.AddPolicy(sub, dom, obj, act)
}

// RemovePolicy removes a policy rule
func (e *Enforcer) RemovePolicy(sub, dom, obj, act string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.Enforcer.RemovePolicy(sub, dom, obj, act)
}

// AddPolicies adds multiple policy rules
func (e *Enforcer) AddPolicies(rules [][]string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.Enforcer.AddPolicies(rules)
}

// RemovePolicies removes multiple policy rules
func (e *Enforcer) RemovePolicies(rules [][]string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.Enforcer.RemovePolicies(rules)
}

// AddGroupingPolicies adds multiple grouping policies (role assignments)
func (e *Enforcer) AddGroupingPolicies(rules [][]string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.Enforcer.AddGroupingPolicies(rules)
}

// RemoveGroupingPolicies removes multiple grouping policies
func (e *Enforcer) RemoveGroupingPolicies(rules [][]string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.Enforcer.RemoveGroupingPolicies(rules)
}

// AddRoleForUserInDomain adds a role for a user in a specific domain
func (e *Enforcer) AddRoleForUserInDomain(userID, role, domain string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.Enforcer.AddGroupingPolicy(userID, role, domain)
}

// RemoveRoleForUserInDomain removes a role for a user in a specific domain
func (e *Enforcer) RemoveRoleForUserInDomain(userID, role, domain string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.Enforcer.RemoveGroupingPolicy(userID, role, domain)
}

// GetRolesForUserInDomain returns all roles for a user in a domain
func (e *Enforcer) GetRolesForUserInDomain(userID, domain string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.GetRolesForUserInDomain(userID, domain)
}

// GetImplicitRolesForUser returns all roles for a user including inherited ones
func (e *Enforcer) GetImplicitRolesForUser(userID string, domain string) ([]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.GetImplicitRolesForUser(userID, domain)
}

// GetImplicitPermissionsForUser returns all permissions for a user including inherited ones
func (e *Enforcer) GetImplicitPermissionsForUser(userID string, domain string) ([][]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.GetImplicitPermissionsForUser(userID, domain)
}

// HasRoleForUserInDomain checks if a user has a specific role in a domain
func (e *Enforcer) HasRoleForUserInDomain(userID, role, domain string) bool {
	roles := e.GetRolesForUserInDomain(userID, domain)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasImplicitRoleForUserInDomain checks if a user has a role (including inherited) in a domain
func (e *Enforcer) HasImplicitRoleForUserInDomain(userID, role, domain string) (bool, error) {
	roles, err := e.GetImplicitRolesForUser(userID, domain)
	if err != nil {
		return false, err
	}
	for _, r := range roles {
		if r == role {
			return true, nil
		}
	}
	return false, nil
}

// AddBusinessRole adds a business role for a user
func (e *Enforcer) AddBusinessRole(userID string, businessID string, role Role) (bool, error) {
	domain := BusinessDomain(businessID)
	return e.AddRoleForUserInDomain(userID, role.String(), domain)
}

// RemoveBusinessRole removes a business role from a user
func (e *Enforcer) RemoveBusinessRole(userID string, businessID string, role Role) (bool, error) {
	domain := BusinessDomain(businessID)
	return e.RemoveRoleForUserInDomain(userID, role.String(), domain)
}

// GetUserBusinessRoles returns all roles a user has in a business
func (e *Enforcer) GetUserBusinessRoles(userID string, businessID string) []string {
	domain := BusinessDomain(businessID)
	return e.GetRolesForUserInDomain(userID, domain)
}

// GetUserImplicitBusinessRoles returns all roles including inherited for a user in a business
func (e *Enforcer) GetUserImplicitBusinessRoles(userID string, businessID string) ([]string, error) {
	domain := BusinessDomain(businessID)
	return e.GetImplicitRolesForUser(userID, domain)
}

// AddPlatformRole adds a platform-level role for a user
func (e *Enforcer) AddPlatformRole(userID string, role Role) (bool, error) {
	return e.AddRoleForUserInDomain(userID, role.String(), DomainPlatform)
}

// RemovePlatformRole removes a platform-level role from a user
func (e *Enforcer) RemovePlatformRole(userID string, role Role) (bool, error) {
	return e.RemoveRoleForUserInDomain(userID, role.String(), DomainPlatform)
}

// GetUserPlatformRoles returns all platform-level roles for a user
func (e *Enforcer) GetUserPlatformRoles(userID string) []string {
	return e.GetRolesForUserInDomain(userID, DomainPlatform)
}

// IsBusinessOwner checks if a user is the owner of a business
func (e *Enforcer) IsBusinessOwner(userID string, businessID string) bool {
	return e.HasRoleForUserInDomain(userID, RoleBusinessOwner.String(), BusinessDomain(businessID))
}

// IsBusinessStaff checks if a user is staff of a business
func (e *Enforcer) IsBusinessStaff(userID string, businessID string) bool {
	return e.HasRoleForUserInDomain(userID, RoleBusinessStaff.String(), BusinessDomain(businessID))
}

// IsConsumer checks if a user is a consumer
func (e *Enforcer) IsConsumer(userID string) bool {
	return e.HasRoleForUserInDomain(userID, RoleConsumer.String(), DomainPlatform)
}

// IsAdmin checks if a user is a platform admin
func (e *Enforcer) IsAdmin(userID string) bool {
	return e.HasRoleForUserInDomain(userID, RoleAdmin.String(), DomainPlatform)
}

// IsSuperAdmin checks if a user is a super admin
func (e *Enforcer) IsSuperAdmin(userID string) bool {
	return e.HasRoleForUserInDomain(userID, RoleSuperAdmin.String(), DomainPlatform)
}

// GetFilteredPolicy gets policies filtered by field values
func (e *Enforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) ([][]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
}

// HasPolicy checks if a policy exists
func (e *Enforcer) HasPolicy(sub, dom, obj, act string) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.HasPolicy(sub, dom, obj, act)
}

// HasGroupingPolicy checks if a grouping policy exists (role assignment)
func (e *Enforcer) HasGroupingPolicy(user, role, domain string) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.HasGroupingPolicy(user, role, domain)
}

// GetPolicy returns all policies (with error handling)
func (e *Enforcer) GetPolicy() ([][]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.GetPolicy()
}

// GetGroupingPolicy returns all grouping policies (role assignments)
func (e *Enforcer) GetGroupingPolicy() ([][]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Enforcer.GetGroupingPolicy()
}