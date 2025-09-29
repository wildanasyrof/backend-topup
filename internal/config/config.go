package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Db     DatabaseConfig
	JWT    JwtConfig
	Oauth  GoogleOauthConfig
}

type ServerConfig struct {
	Port           string
	RequestTimeOut int
	Env            string
	UploadDir      string
}

type GoogleOauthConfig struct {
	ClientID     string
	ClientSecret string
	CallbackUrl  string
}

type DatabaseConfig struct {
	DbUrl string
}

type JwtConfig struct {
	JWTSecret          string
	AccessTokenMinutes int
	RefreshTokenDays   int
}

// LoadConfig reads configuration from environment variables and returns a Config struct.
func LoadConfig() (*Config, error) {
	// Attempt to load .env file. It's OK if it doesn't exist, as it might be
	// running in a production environment with env variables already set.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	// Helper function to get a string environment variable with a default value.
	getEnv := func(key, defaultValue string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return defaultValue
	}

	// Helper function to get an int environment variable.
	getEnvAsInt := func(key string) (int, error) {
		valueStr := os.Getenv(key)
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value, nil
		}
		// Return 0 and an error if the key doesn't exist or is not an integer.
		return 0, nil
	}

	requestTimeOut, err := getEnvAsInt("HTTP_REQUEST_TIME_OUT")
	if err != nil {
		log.Fatal("HTTP_REQUEST_TIME_OUT must be a valid integer")
	}

	accessTokenMinutes, err := getEnvAsInt("ACCESS_TOKEN_MINUTES")
	if err != nil {
		log.Fatal("ACCESS_TOKEN_MINUTES must be a valid integer")
	}

	refreshTokenDays, err := getEnvAsInt("REFRESH_TOKEN_DAYS")
	if err != nil {
		log.Fatal("REFRESH_TOKEN_DAYS must be a valid integer")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:           getEnv("PORT", "8080"),
			RequestTimeOut: requestTimeOut,
			Env:            getEnv("ENV", "development"),
			UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		},
		Db: DatabaseConfig{
			DbUrl: getEnv("DATABASE_URL", ""),
		},
		JWT: JwtConfig{
			JWTSecret:          getEnv("JWT_SECRET", ""),
			AccessTokenMinutes: accessTokenMinutes,
			RefreshTokenDays:   refreshTokenDays,
		},
		Oauth: GoogleOauthConfig{
			ClientID:     getEnv("GOOGLE_OAUTH_CLIENT_ID", "719094701967-kfms9f9b27pahl58a64lm1uj5bq9b1sk.apps.googleusercontent.com"),
			ClientSecret: getEnv("GOOGLE_OAUTH_SECRET_ID", "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"),
			CallbackUrl:  getEnv("GOOGLE_OAUTH_CALLBACK_URL", ""),
		},
	}

	return cfg, nil
}
