package config

import (
    "log"
    "os"
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
    Environment string
}

type OTPConfig struct {
    TTL time.Duration  // How long OTP is valid
}

type RedisConfig struct {
    URL string
}

type ResendConfig struct {
	ApiKey string
	From   string
}

type ServerConfig struct {
    Port string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}

type JWTConfig struct {
    Secret          string
    Expiration      time.Duration
    RefreshExpiration time.Duration
}

type EmailConfig struct {
    APIKey string
    From   string
}

func Load() *Config {
    // Load .env file
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
            TTL: 5 * time.Minute,  // OTP expires in 5 minutes
        },
        Redis: RedisConfig{
            URL: getEnv("REDIS_URL", "localhost:6379"),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}