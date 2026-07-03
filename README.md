# Flownatty Backend API

Flownatty MVP1 Backend API - E-commerce marketplace built with Go and Gin.

## Tech Stack

- **Language:** Go 1.26.1
- **Framework:** Gin
- **Database:** PostgreSQL 15+
- **ORM:** GORM
- **Auth:** JWT
- **Payments:** M-Pesa Daraja API

## Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Git

## Installation

```bash
# Clone repository
git clone https://github.com/ALLAN-star-glitch/flownatty-backend.git
cd flownatty-backend

# Install dependencies
go mod download

# Create .env file
cp .env.example .env
# Edit .env with your credentials

# Create database
createdb flownatty

# Run migrations
goose -dir internal/database/migrations postgres "postgresql://localhost:5432/flownatty?sslmode=disable" up

# Run application
go run cmd/api/main.go