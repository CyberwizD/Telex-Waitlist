package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	AppName        string
	Port           string
	DatabaseURL    string
	AdminAPIKey    string
	AllowedOrigins []string

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	EmailEnabled bool
	LogLevel     string
}

// Load parses environment variables into a Config struct.
// It panics if required values are missing.
func Load() (*Config, error) {
	appName := getEnv("APP_NAME", "Telex Waitlist")
	port := getEnv("PORT", "8080")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := getEnv("SMTP_PORT", "587")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")

	emailEnabled := strings.EqualFold(getEnv("EMAIL_ENABLED", "true"), "true")
	logLevel := getEnv("LOG_LEVEL", "info")

	var allowedOrigins []string
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		for _, o := range strings.Split(origins, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}

	return &Config{
		AppName:        appName,
		Port:           port,
		DatabaseURL:    dbURL,
		AdminAPIKey:    os.Getenv("ADMIN_API_KEY"),
		AllowedOrigins: allowedOrigins,
		SMTPHost:       smtpHost,
		SMTPPort:       smtpPort,
		SMTPUsername:   smtpUser,
		SMTPPassword:   smtpPass,
		SMTPFrom:       smtpFrom,
		EmailEnabled:   emailEnabled,
		LogLevel:       logLevel,
	}, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
