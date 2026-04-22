package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort     string
	DatabaseURL   string
	RedisURL      string
	JWTSecret     string
	JWTExpiry     time.Duration
	RedisExpiry  time.Duration
	Environment  string
	Debug        bool
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPass     string
	SMTPFrom     string
	WhatsAppAPI  string
	MidtransURL  string
	MidtransKey  string
	MidtransServerKey string
	XenditURL    string
	XenditKey    string
	S3Endpoint  string
	S3Bucket    string
	S3AccessKey string
	S3SecretKey string
	S3UseSSL    bool
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/edusys?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:      15 * time.Minute,
		RedisExpiry:   24 * time.Hour,
		Environment:   getEnv("ENVIRONMENT", "development"),
		Debug:         getEnv("DEBUG", "false") == "true",
		SMTPHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPass:      getEnv("SMTP_PASS", ""),
		SMTPFrom:      getEnv("SMTP_FROM", "noreply@school.edu"),
		WhatsAppAPI:   getEnv("WHATSAPP_API", ""),
		MidtransURL:   getEnv("MIDTRANS_URL", "https://api.midtrans.com"),
		MidtransKey:   getEnv("MIDTRANS_KEY", ""),
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		XenditURL:     getEnv("XENDIT_URL", "https://api.xendit.co"),
		XenditKey:    getEnv("XENDIT_KEY", ""),
		S3Endpoint:   getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3Bucket:     getEnv("S3_BUCKET", "edusys"),
		S3AccessKey:  getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:  getEnv("S3_SECRET_KEY", ""),
		S3UseSSL:     getEnv("S3_SSL", "false") == "true",
	}
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