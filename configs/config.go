package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
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
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := getEnv("JWT_SECRET", "dev-secret-key")

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("incorrect  DB_PORT: %v", err)
	}

	dbPassword := getEnv("DB_PASSWORD", "")
	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: dbPassword,
		DBName:     getEnv("DB_NAME", "postgres"),
		JWTSecret:  jwtSecret,
	}, nil
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
