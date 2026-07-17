// main.go

package main

import (
	"context"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/app"
)

func main() {
	// Create application - handles everything: config, Redis, DB, modules
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}
	defer application.Close()

	// Initialize modules (seed data, etc.)
	ctx := context.Background()
	if err := application.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Setup routes
	application.SetupRoutes()

	// Start server
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}