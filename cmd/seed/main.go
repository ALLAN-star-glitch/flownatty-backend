// cmd/seed/main.go
package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database/seeders"
	"gorm.io/gorm"
)

type SeederLog struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SeederName  string    `gorm:"uniqueIndex"`
	Checksum    string    `gorm:"type:varchar(64)"`
	ExecutedAt  time.Time `gorm:"autoCreateTime"`
}

func (SeederLog) TableName() string {
	return "seeder_logs"
}

// SeederInfo holds information about a seeder
type SeederInfo struct {
	Name        string
	Fn          func(*gorm.DB) error
	Description string
	Deps        []string
}

func main() {
	// Parse command line flags
	var (
		only        = flag.String("only", "", "Run only specific seeder")
		env         = flag.String("env", "development", "Environment (development, test, production)")
		force       = flag.Bool("force", false, "Force re-seed even if already run")
		dryRun      = flag.Bool("dry-run", false, "Preview what would be seeded")
		skipConfirm = flag.Bool("skip-confirm", false, "Skip confirmation prompts (for CI/CD)")
		help        = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// Load config
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Auto-migrate seeder logs table (don't fatal on error)
	if err := db.AutoMigrate(&SeederLog{}); err != nil {
		log.Printf("⚠️  Warning: Failed to create seeder_logs table: %v", err)
	}

	// Define seeders with dependencies
	seedersList := []SeederInfo{
		{
			Name:        "business_types",
			Fn:          seeders.SeedBusinessTypes,
			Description: "Seed business legal structures",
			Deps:        []string{},
		},
		{
			Name:        "business_sectors",
			Fn:          seeders.SeedBusinessSectors,
			Description: "Seed business industry sectors",
			Deps:        []string{},
		},
		{
			Name:        "business_subcategories",
			Fn:          seeders.SeedBusinessSubcategories,
			Description: "Seed detailed business subcategories",
			Deps:        []string{"business_sectors"},
		},
		{
			Name:        "product_service_categories",
			Fn:          seeders.SeedProductServiceCategories,
			Description: "Seed product and service categories",
			Deps:        []string{},
		},
		{
			Name:        "product_service_subcategories",
			Fn:          seeders.SeedProductServiceSubcategories,
			Description: "Seed product and service subcategories",
			Deps:        []string{"product_service_categories"},
		},
		{
			Name:        "establishment_types",
			Fn:          seeders.SeedEstablishmentTypes,
			Description: "Seed establishment types",
			Deps:        []string{},
		},
		{
			Name:        "admin",
			Fn:          seeders.SeedAdmin,
			Description: "Seed admin user and permissions",
			Deps:        []string{},
		},
	}

	// Always validate dependencies
	if err := validateDependencies(seedersList); err != nil {
		log.Fatalf("❌ Dependency validation failed: %v", err)
	}

	// Dry run mode
	if *dryRun {
		runDryRun(seedersList, *only)
		return
	}

	// Production safety check
	if *env == "production" {
		log.Println("⚠️  Running in PRODUCTION mode")
		if !*skipConfirm {
			log.Print("Are you sure you want to seed in production? (yes/no): ")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				log.Println("❌ Seeding cancelled")
				os.Exit(1)
			}
		} else {
			log.Println("⚠️  Skipping confirmation (CI/CD mode)")
		}
	}

	// Run specific seeder if requested
	if *only != "" {
		runSpecificSeeder(db, seedersList, *only, *env, *force)
		return
	}

	// Run all seeders in order
	runAllSeeders(db, seedersList, *env, *force)
}

func runAllSeeders(db *gorm.DB, seedersList []SeederInfo, env string, force bool) {
	log.Printf("🚀 Starting all seeders (env: %s)", env)

	total := len(seedersList)
	success := 0

	for i, s := range seedersList {
		log.Printf("[%d/%d] Running seeder: %s", i+1, total, s.Name)

		// Check if already run in production
		if env == "production" && !force {
			if alreadyRun(db, s.Name) {
				log.Printf("⏭️  Seeder %s already run in production, skipping...", s.Name)
				success++
				continue
			}
		}

		// Run seeder with transaction
		var seededData interface{}
		err := db.Transaction(func(tx *gorm.DB) error {
			// Run the seeder
			if err := s.Fn(tx); err != nil {
				return err
			}
			seededData = s.Name
			return nil
		})

		if err != nil {
			log.Fatalf("❌ Failed to run seeder %s: %v", s.Name, err)
		}

		// Log successful seed with checksum
		if err := logSeeded(db, s.Name, seededData); err != nil {
			log.Printf("⚠️  Failed to log seeder: %v", err)
		}

		success++
	}

	log.Printf("✅ All seeders completed successfully (%d/%d)", success, total)
}

func runSpecificSeeder(db *gorm.DB, seedersList []SeederInfo, seederName, env string, force bool) {
	found := false

	for _, s := range seedersList {
		if s.Name == seederName {
			found = true
			log.Printf("🚀 Running seeder: %s (env: %s)", s.Name, env)

			// Check if already run in production
			if env == "production" && !force {
				if alreadyRun(db, s.Name) {
					log.Printf("⏭️  Seeder %s already run in production, skipping...", s.Name)
					return
				}
			}

			// Run with transaction
			var seededData interface{}
			err := db.Transaction(func(tx *gorm.DB) error {
				if err := s.Fn(tx); err != nil {
					return err
				}
				seededData = s.Name
				return nil
			})

			if err != nil {
				log.Fatalf("❌ Failed to run seeder %s: %v", s.Name, err)
			}

			if err := logSeeded(db, s.Name, seededData); err != nil {
				log.Printf("⚠️  Failed to log seeder: %v", err)
			}

			log.Printf("✅ Seeder %s completed successfully", s.Name)
			break
		}
	}

	if !found {
		log.Printf("❌ Unknown seeder: %s", seederName)
		printAvailableSeeders(seedersList)
		os.Exit(1)
	}
}

func runDryRun(seedersList []SeederInfo, only string) {
	log.Println("🔍 DRY RUN - Would seed:")
	if only != "" {
		for _, s := range seedersList {
			if s.Name == only {
				log.Printf("  - %s: %s", s.Name, s.Description)
				if len(s.Deps) > 0 {
					log.Printf("    Depends on: %v", s.Deps)
				}
			}
		}
	} else {
		for _, s := range seedersList {
			log.Printf("  - %s: %s", s.Name, s.Description)
			if len(s.Deps) > 0 {
				log.Printf("    Depends on: %v", s.Deps)
			}
		}
	}
	log.Println("✅ Dry run completed (no changes made)")
}

func alreadyRun(db *gorm.DB, seederName string) bool {
	var count int64
	db.Model(&SeederLog{}).Where("seeder_name = ?", seederName).Count(&count)
	return count > 0
}

// cmd/seed/main.go

func logSeeded(db *gorm.DB, seederName string, data interface{}) error {
	checksum := calculateChecksum(data)

	// ✅ Check if log already exists
	var existingLog SeederLog
	err := db.Where("seeder_name = ?", seederName).First(&existingLog).Error
	if err == nil {
		// ✅ Log exists, update it
		existingLog.Checksum = checksum
		existingLog.ExecutedAt = time.Now()
		return db.Save(&existingLog).Error
	}
	
	if err != gorm.ErrRecordNotFound {
		// Some other error
		return err
	}

	// ✅ Log doesn't exist, create it
	log := SeederLog{
		SeederName: seederName,
		Checksum:   checksum,
		ExecutedAt: time.Now(),
	}
	return db.Create(&log).Error
}

func calculateChecksum(data interface{}) string {
	if data == nil {
		return ""
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	hash := md5.Sum(jsonData)
	return hex.EncodeToString(hash[:])
}

func validateDependencies(seeders []SeederInfo) error {
	// Build a map of available seeders
	available := make(map[string]bool)
	for _, s := range seeders {
		available[s.Name] = true
	}

	// Check each seeder's dependencies
	for _, s := range seeders {
		for _, dep := range s.Deps {
			if !available[dep] {
				return fmt.Errorf("seeder '%s' depends on '%s' which doesn't exist", s.Name, dep)
			}
		}
	}
	return nil
}

func printHelp() {
	fmt.Println("Flownatty Database Seeder")
	fmt.Println("\nUsage:")
	fmt.Println("  go run cmd/seed/main.go [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -only <seeder>    Run only a specific seeder")
	fmt.Println("  -env <env>       Environment (development, test, production)")
	fmt.Println("  -force           Force re-seed even if already run")
	fmt.Println("  -dry-run        Preview what would be seeded")
	fmt.Println("  -skip-confirm   Skip confirmation prompts (for CI/CD)")
	fmt.Println("  -help            Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  # Run all seeders in development")
	fmt.Println("  go run cmd/seed/main.go")
	fmt.Println("")
	fmt.Println("  # Run only business types in production")
	fmt.Println("  go run cmd/seed/main.go -only=business_types -env=production -force")
	fmt.Println("")
	fmt.Println("  # Preview what would be seeded")
	fmt.Println("  go run cmd/seed/main.go -dry-run")
	fmt.Println("")
	fmt.Println("  # Run in CI/CD (no prompts)")
	fmt.Println("  go run cmd/seed/main.go -env=production -skip-confirm")
}

func printAvailableSeeders(seeders []SeederInfo) {
	fmt.Println("\nAvailable Seeders:")
	for _, s := range seeders {
		fmt.Printf("  - %-30s %s\n", s.Name, s.Description)
		if len(s.Deps) > 0 {
			fmt.Printf("    Depends on: %v\n", s.Deps)
		}
	}
}