package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration structs
type Config struct {
	Server ServerConfig
	Db     DatabaseConfig
	JWT    JwtConfig
	Oauth  GoogleOauthConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port           string
	RequestTimeOut int
	Env            string
	UploadDir      string
}

// GoogleOauthConfig holds Google OAuth specific configuration
type GoogleOauthConfig struct {
	ClientID     string
	ClientSecret string
	CallbackUrl  string
}

// DatabaseConfig holds database connection details
type DatabaseConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	// DbUrl can be constructed if needed, but keeping components is often better
	// DbUrl string
}

// JwtConfig holds JWT secret and lifetime configuration
type JwtConfig struct {
	AccessSecret       string // Renamed from JWTSecret
	RefreshSecret      string // Added for refresh token
	AccessTokenMinutes int
	RefreshTokenDays   int
}

// LoadConfig reads configuration from environment variables and returns a Config struct.
func LoadConfig() (*Config, error) {
	// Attempt to load .env file. It's OK if it doesn't exist.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, reading from environment variables")
	}

	// Helper function to get a string environment variable with a default value.
	getEnv := func(key, defaultValue string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		log.Printf("Environment variable %s not found, using default value: %s", key, defaultValue)
		return defaultValue
	}

	// Helper function to get an int environment variable with a default value and error handling.
	getEnvAsInt := func(key string, defaultValue int) (int, error) {
		valueStr := os.Getenv(key)
		if valueStr == "" {
			log.Printf("Environment variable %s not found, using default value: %d", key, defaultValue)
			return defaultValue, nil
		}
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			// Return an error instead of fatal logging
			return 0, fmt.Errorf("invalid integer value for %s: %w", key, err)
		}
		return value, nil
	}

	// --- Server Config ---
	requestTimeOut, err := getEnvAsInt("HTTP_REQUEST_TIME_OUT", 15) // Default to 15
	if err != nil {
		return nil, err // Return error to caller
	}

	// --- JWT Config ---
	accessTokenMinutes, err := getEnvAsInt("ACCESS_TOKEN_MINUTES", 15) // Default to 15
	if err != nil {
		return nil, err
	}

	refreshTokenDays, err := getEnvAsInt("REFRESH_TOKEN_DAYS", 30) // Default to 30
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:           getEnv("PORT", "8080"), // Default to 8080 if not set
			RequestTimeOut: requestTimeOut,
			Env:            getEnv("ENV", "development"),
			UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		},
		Db: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "5432"),
			Database: getEnv("DB_DATABASE", ""), // No sensible default for database name
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", ""), // No default for password
			// DbUrl: dbURL, // You can assign the constructed URL here if needed elsewhere
		},
		JWT: JwtConfig{
			AccessSecret:       getEnv("ACCESS_SECRET", ""),  // No default for secrets
			RefreshSecret:      getEnv("REFRESH_SECRET", ""), // No default for secrets
			AccessTokenMinutes: accessTokenMinutes,
			RefreshTokenDays:   refreshTokenDays,
		},
		Oauth: GoogleOauthConfig{
			ClientID:     getEnv("GOOGLE_OAUTH_CLIENT_ID", ""),     // No default for client ID
			ClientSecret: getEnv("GOOGLE_OAUTH_CLIENT_SECRET", ""), // No default for client secret
			CallbackUrl:  getEnv("GOOGLE_OAUTH_CALLBACK_URL", ""),  // No default for callback URL
		},
	}

	// Basic validation for essential empty values
	if cfg.JWT.AccessSecret == "" {
		log.Println("Warning: ACCESS_SECRET is not set.")
	}
	if cfg.JWT.RefreshSecret == "" {
		log.Println("Warning: REFRESH_SECRET is not set.")
	}
	if cfg.Db.Database == "" {
		log.Println("Warning: DB_DATABASE is not set.")
	}
	// Add more checks as needed, potentially returning errors for critical missing values

	return cfg, nil
}
