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
    Supabase    SupabaseConfig
    JWT         JWTConfig
    Mpesa       MpesaConfig
    Twilio      TwilioConfig
    Firebase    FirebaseConfig
    R2          R2Config
    Redis       RedisConfig
    Environment string
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

type SupabaseConfig struct {
    URL        string
    AnonKey    string
    ServiceKey string
}

type JWTConfig struct {
    Secret          string
    Expiration      time.Duration
    RefreshExpiration time.Duration
}

type MpesaConfig struct {
    ConsumerKey    string
    ConsumerSecret string
    Shortcode      string
    Passkey        string
    Environment    string
}

type TwilioConfig struct {
    AccountSID  string
    AuthToken   string
    PhoneNumber string
}

type FirebaseConfig struct {
    ServerKey string
}

type R2Config struct {
    AccessKey string
    SecretKey string
    Bucket    string
    Endpoint  string
}

type RedisConfig struct {
    URL string
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
        Supabase: SupabaseConfig{
            URL:        getEnv("SUPABASE_URL", ""),
            AnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
            ServiceKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
        },
        JWT: JWTConfig{
            Secret:            getEnv("JWT_SECRET", "your-super-secret-key-change-this"),
            Expiration:        24 * time.Hour,
            RefreshExpiration: 168 * time.Hour,
        },
        Mpesa: MpesaConfig{
            ConsumerKey:    getEnv("MPESA_CONSUMER_KEY", ""),
            ConsumerSecret: getEnv("MPESA_CONSUMER_SECRET", ""),
            Shortcode:      getEnv("MPESA_SHORTCODE", ""),
            Passkey:        getEnv("MPESA_PASSKEY", ""),
            Environment:    getEnv("MPESA_ENVIRONMENT", "sandbox"),
        },
        Twilio: TwilioConfig{
            AccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
            AuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
            PhoneNumber: getEnv("TWILIO_PHONE_NUMBER", ""),
        },
        Firebase: FirebaseConfig{
            ServerKey: getEnv("FCM_SERVER_KEY", ""),
        },
        R2: R2Config{
            AccessKey: getEnv("R2_ACCESS_KEY", ""),
            SecretKey: getEnv("R2_SECRET_KEY", ""),
            Bucket:    getEnv("R2_BUCKET", "flownatty"),
            Endpoint:  getEnv("R2_ENDPOINT", ""),
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