package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DBHost string `json:"DBHost"`

	DBPort     int    `json:"DBPort"`
	DBUser     string `json:"DBUser"`
	DBPassword string `json:"-"`
	DBName     string `json:"DBName"`
	JWTSecret  string `json:"-"`

	// Keep-alive configuration
	KeepAliveEnabled  bool   `json:"KeepAliveEnabled"`
	KeepAliveURL      string `json:"KeepAliveURL"`
	KeepAliveInterval int    `json:"KeepAliveInterval"` // in minutes
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	jwtSecret := getEnv("JWT_SECRET", "dev-secret-key")

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("incorrect  DB_PORT: %v", err)
	}

	dbPassword := getEnv("DB_PASSWORD", "")
	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	// Keep-alive configuration
	keepAliveEnabled := getEnv("KEEP_ALIVE_ENABLED", "true") == "true"
	keepAliveURL := getEnv("KEEP_ALIVE_URL", "http://localhost:8080")
	keepAliveInterval, err := strconv.Atoi(getEnv("KEEP_ALIVE_INTERVAL", "10"))
	if err != nil {
		keepAliveInterval = 10 // default to 10 minutes
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: dbPassword,
		DBName:     getEnv("DB_NAME", "microshopik"),
		JWTSecret:  jwtSecret,

		KeepAliveEnabled:  keepAliveEnabled,
		KeepAliveURL:      keepAliveURL,
		KeepAliveInterval: keepAliveInterval,
	}, nil
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
