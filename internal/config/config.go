package config

// Purpose: Loads configuration from environment variables.
// What it does:

// Reads environment variables (or .env file)

// Provides Casbin configuration to the permissions module

// Tells Casbin where to find the model and policy files

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	Email       EmailConfig
	Resend      ResendConfig
	Redis       RedisConfig
	OTP         OTPConfig
	Casbin      CasbinConfig
	Environment string
}

// OTPConfig defines OTP configuration
type OTPConfig struct {
	TTL time.Duration
}

// RedisConfig defines Redis configuration
type RedisConfig struct {
	URL string
}

// ResendConfig defines Resend email configuration
type ResendConfig struct {
	ApiKey string
	From   string
}

// ServerConfig defines server configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig defines database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig defines JWT configuration
type JWTConfig struct {
	Secret            string
	Expiration        time.Duration // Access token expiration time - 24 hours
	RefreshExpiration time.Duration // Refresh token expiration time - 7 days
}

// EmailConfig defines email configuration
type EmailConfig struct {
	APIKey string
	From   string
}

// CasbinConfig defines Casbin authorization configuration
type CasbinConfig struct {
	ModelPath        string
	PolicyPath       string
	AutoLoad         bool
	AutoLoadInterval time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file from project root
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "flownatty"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "your-super-secret-key-change-this"),
			Expiration:        24 * time.Hour,
			RefreshExpiration: 168 * time.Hour,
		},
		Email: EmailConfig{
			APIKey: getEnv("SENDGRID_API_KEY", ""),
			From:   getEnv("EMAIL_FROM", "noreply@flownatty.com"),
		},
		Resend: ResendConfig{
			ApiKey: getEnv("RESEND_API_KEY", ""),
			From:   getEnv("EMAIL_FROM", "noreply@flownatty.com"),
		},
		OTP: OTPConfig{
			TTL: 5 * time.Minute,
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "localhost:6379"),
		},
		Casbin: CasbinConfig{
			ModelPath:        getEnv("CASBIN_MODEL", "configs/casbin/model.conf"),
			PolicyPath:       getEnv("CASBIN_POLICY", "configs/casbin/policies/default_policies.csv"),
			AutoLoad:         getEnvBool("CASBIN_AUTO_LOAD", true),
			AutoLoadInterval: getEnvDuration("CASBIN_AUTO_LOAD_INTERVAL", 10*time.Second),
		},
	}
}

// getEnv returns environment variable or default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool returns boolean environment variable or default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return parsed
	}
	return defaultValue
}

// getEnvDuration returns duration environment variable or default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil {
			return defaultValue
		}
		return parsed
	}
	return defaultValue
}