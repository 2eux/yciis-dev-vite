package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	DatabaseURL       string
	RedisHost         string
	RedisPassword     string
	RedisDB           int
	JWTSecret         string
	JWTExpiry         time.Duration
	RefreshExpiry     time.Duration
	Environment       string
	Debug             bool
	AllowedOrigins    []string
	RateLimitRequests int
	RateLimitWindowMs int
	SMTPHost          string
	SMTPPort          int
	SMTPUser          string
	SMTPPass          string
	SMTPFrom          string
	WhatsAppAPI       string
	MidtransURL       string
	MidtransKey       string
	MidtransServerKey string
	XenditURL         string
	XenditKey         string
	S3Endpoint        string
	S3Bucket          string
	S3AccessKey       string
	S3SecretKey       string
	S3UseSSL          bool
}

// Load reads configuration from environment variables with validation.
// Critical secrets MUST be set; the server will refuse to start with defaults.
func Load() *Config {
	godotenv.Load()

	env := getEnv("ENVIRONMENT", "development")

	// In production, enforce mandatory secrets
	jwtSecret := getEnv("JWT_SECRET", "")
	databaseURL := getEnv("DATABASE_URL", "")

	if env == "production" {
		if jwtSecret == "" {
			log.Fatal("FATAL: JWT_SECRET environment variable is required in production")
		}
		if len(jwtSecret) < 32 {
			log.Fatal("FATAL: JWT_SECRET must be at least 32 characters long")
		}
		if databaseURL == "" {
			log.Fatal("FATAL: DATABASE_URL environment variable is required in production")
		}
	} else {
		// Development defaults with warnings
		if jwtSecret == "" {
			jwtSecret = "DEVELOPMENT-ONLY-INSECURE-SECRET-CHANGE-ME"
			log.Println("WARNING: Using insecure default JWT_SECRET. Set JWT_SECRET env var for production.")
		}
		if databaseURL == "" {
			databaseURL = "postgres://postgres:password@localhost:5432/edusys?sslmode=disable"
			log.Println("WARNING: Using default DATABASE_URL. Set DATABASE_URL env var for production.")
		}
	}

	// Parse allowed origins
	originsRaw := getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	origins := strings.Split(originsRaw, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	// Parse Redis URL into host:port format
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")
	redisHost := redisURL
	redisPassword := ""
	redisDB := 0

	// Strip redis:// prefix and parse
	redisHost = strings.TrimPrefix(redisHost, "redis://")
	if parts := strings.SplitN(redisHost, "@", 2); len(parts) == 2 {
		redisPassword = parts[0]
		redisHost = parts[1]
	}
	if parts := strings.SplitN(redisHost, "/", 2); len(parts) == 2 {
		redisHost = parts[0]
		if db, err := strconv.Atoi(parts[1]); err == nil {
			redisDB = db
		}
	}

	jwtExpiryMinutes := getEnvAsInt("JWT_EXPIRY_MINUTES", 15)
	refreshExpiryDays := getEnvAsInt("REFRESH_TOKEN_EXPIRY_DAYS", 30)

	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DatabaseURL:       databaseURL,
		RedisHost:         redisHost,
		RedisPassword:     redisPassword,
		RedisDB:           redisDB,
		JWTSecret:         jwtSecret,
		JWTExpiry:         time.Duration(jwtExpiryMinutes) * time.Minute,
		RefreshExpiry:     time.Duration(refreshExpiryDays) * 24 * time.Hour,
		Environment:       env,
		Debug:             getEnv("DEBUG", "false") == "true",
		AllowedOrigins:    origins,
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindowMs: getEnvAsInt("RATE_LIMIT_WINDOW_MS", 60000),
		SMTPHost:          getEnv("SMTP_HOST", ""),
		SMTPPort:          getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:          getEnv("SMTP_USER", ""),
		SMTPPass:          getEnv("SMTP_PASS", ""),
		SMTPFrom:          getEnv("SMTP_FROM", ""),
		WhatsAppAPI:       getEnv("WHATSAPP_API", ""),
		MidtransURL:       getEnv("MIDTRANS_URL", "https://api.midtrans.com"),
		MidtransKey:       getEnv("MIDTRANS_KEY", ""),
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		XenditURL:         getEnv("XENDIT_URL", "https://api.xendit.co"),
		XenditKey:         getEnv("XENDIT_KEY", ""),
		S3Endpoint:        getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3Bucket:          getEnv("S3_BUCKET", "edusys"),
		S3AccessKey:       getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:       getEnv("S3_SECRET_KEY", ""),
		S3UseSSL:          getEnv("S3_USE_SSL", "false") == "true",
	}
}

// IsProduction returns true if the environment is set to production.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Validate performs runtime checks on configuration integrity.
func (c *Config) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET cannot be empty")
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL cannot be empty")
	}
	return nil
}