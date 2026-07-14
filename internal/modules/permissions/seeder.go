package permissions

// Purpose: Seeds initial policies into the database.
// What it does:

// Checks if policies already exist in database

// If not, loads policies from default_policies.csv

// Loads role hierarchy from role_hierarchy.csv

// Only runs once (idempotent)

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Seeder handles initial policy seeding
type Seeder struct {
	enforcer *Enforcer
	service  *Service
}

// NewSeeder creates a new seeder
func NewSeeder(enforcer *Enforcer, service *Service) *Seeder {
	return &Seeder{
		enforcer: enforcer,
		service:  service,
	}
}

// SeedDatabase seeds the database with initial policies and roles
func (s *Seeder) SeedDatabase(ctx context.Context) error {
	log.Println("Starting database seeding...")

	// Check if policies already exist
	policies, err := s.enforcer.GetPolicy()
	if err != nil {
		return fmt.Errorf("failed to get policies: %w", err)
	}

	if len(policies) > 0 {
		log.Printf("Policies already exist (%d rules), skipping seed", len(policies))
		return nil
	}

	// 1. Seed role hierarchy from CSV
	if err := s.seedRoleHierarchyFromCSV(ctx); err != nil {
		return err
	}

	// 2. Seed platform policies from CSV
	if err := s.seedPoliciesFromCSV(ctx); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// seedPoliciesFromCSV seeds policies from default_policies.csv
func (s *Seeder) seedPoliciesFromCSV(ctx context.Context) error {
	file, err := os.Open("configs/casbin/policies/default_policies.csv")
	if err != nil {
		return fmt.Errorf("failed to open policies CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1 // Allow variable fields

	var policies [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV: %w", err)
		}

		// Skip empty lines
		if len(record) == 0 || (len(record) == 1 && strings.TrimSpace(record[0]) == "") {
			continue
		}

		// Format: p, sub, dom, obj, act
		if len(record) >= 5 {
			policyType := strings.TrimSpace(record[0])
			if policyType == "p" {
				sub := strings.TrimSpace(record[1])
				dom := strings.TrimSpace(record[2])
				obj := strings.TrimSpace(record[3])
				act := strings.TrimSpace(record[4])

				// Skip template policies with {{.BusinessID}}
				if strings.Contains(dom, "{{.BusinessID}}") {
					continue
				}

				policies = append(policies, []string{sub, dom, obj, act})
			}
		}
	}

	if len(policies) > 0 {
		_, err := s.enforcer.AddPolicies(policies)
		if err != nil {
			return fmt.Errorf("failed to add policies: %w", err)
		}
		log.Printf("Seeded %d policies from CSV", len(policies))
	}

	return nil
}

// seedRoleHierarchyFromCSV seeds role hierarchy from role_hierarchy.csv
func (s *Seeder) seedRoleHierarchyFromCSV(ctx context.Context) error {
	file, err := os.Open("configs/casbin/policies/role_hierarchy.csv")
	if err != nil {
		return fmt.Errorf("failed to open role hierarchy CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	var rules [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV: %w", err)
		}

		// Skip empty lines
		if len(record) == 0 || (len(record) == 1 && strings.TrimSpace(record[0]) == "") {
			continue
		}

		// Format: g, user, role, domain
		if len(record) >= 4 {
			policyType := strings.TrimSpace(record[0])
			if policyType == "g" {
				user := strings.TrimSpace(record[1])
				role := strings.TrimSpace(record[2])
				domain := strings.TrimSpace(record[3])

				rules = append(rules, []string{user, role, domain})
			}
		}
	}

	if len(rules) > 0 {
		_, err := s.enforcer.AddGroupingPolicies(rules)
		if err != nil {
			return fmt.Errorf("failed to add role hierarchy: %w", err)
		}
		log.Printf("Seeded %d role hierarchy entries from CSV", len(rules))
	}

	return nil
}

// IsSeeded checks if the database has been seeded
func (s *Seeder) IsSeeded(ctx context.Context) (bool, error) {
	policies, err := s.enforcer.GetPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to get policies: %w", err)
	}
	return len(policies) > 0, nil
}