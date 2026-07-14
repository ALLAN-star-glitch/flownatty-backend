package main

import (
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/database/seeders"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := seeders.SeedAdmin(db); err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	log.Println("Seeding completed successfully")
}