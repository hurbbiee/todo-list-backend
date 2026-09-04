package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	URL            string
	ConnectTimeout time.Duration
	MaxConnections int32
	MinConnections int32
}

type Config struct {
	Port         string
	FrontendURL  string
	JWTSecret    []byte
	CrossOrigins string
	Database     DatabaseConfig
}

func Load() (Config, error) {
	// ถ้าไม่มี .env เช่นตอนรันใน Docker ก็ให้อ่านจาก environment ต่อ
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env file not found, using system env")
	}

	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	cfg := &Config{
		Port:         getEnv("PORT", "3000"),
		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:5173"),
		CrossOrigins: getEnv("CROSS_ORIGINS", "http://localhost:5173"),
		Database: DatabaseConfig{
			URL:            databaseURL,
			ConnectTimeout: 5 * time.Second,
			MaxConnections: 10,
			MinConnections: 1,
		},
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if os.Getenv("APP_ENV") == "production" {
			log.Fatal("JWT_SECRET is required in production")
		}
		log.Println("JWT_SECRET not set, using dev default")
		jwtSecret = "dev-super-secret-key"
	}

	cfg.JWTSecret = []byte(jwtSecret)

	return *cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
