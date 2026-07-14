package permissions

// Purpose: Module initialization and lifecycle management.
// What it does:

// Creates and connects all permission components

// Runs seeder on initialization

// Provides access to enforcer and service

// Handles cleanup on shutdown

import (
	"context"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the permissions module
type Module struct {
	enforcer *Enforcer
	service  *Service
	seeder   *Seeder
}

// NewModule creates a new permissions module
func NewModule(db *gorm.DB, cfg *config.Config) (*Module, error) {
	// Initialize enforcer
	enforcer, err := InitEnforcer(db, cfg)
	if err != nil {
		return nil, err
	}

	// Initialize service
	service := NewService(enforcer)

	// Initialize seeder
	seeder := NewSeeder(enforcer, service)

	return &Module{
		enforcer: enforcer,
		service:  service,
		seeder:   seeder,
	}, nil
}

// Init initializes the module (seeds database if needed)
func (m *Module) Init(ctx context.Context) error {
	// Check if already seeded
	seeded, err := m.seeder.IsSeeded(ctx)
	if err != nil {
		log.Printf("Warning: failed to check if seeded: %v", err)
	}

	if !seeded {
		log.Println("Seeding permissions database...")
		if err := m.seeder.SeedDatabase(ctx); err != nil {
			return err
		}
	}

	return nil
}

// GetEnforcer returns the enforcer instance
func (m *Module) GetEnforcer() *Enforcer {
	return m.enforcer
}

// GetService returns the permission service
func (m *Module) GetService() *Service {
	return m.service
}

// RegisterRoutes registers permission-related routes (if any)
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	// No public routes for permissions module
	// All permission checks are done via middleware
}

// Close closes the module and cleans up resources
func (m *Module) Close() {
	if m.enforcer != nil {
		m.enforcer.Close()
	}
}