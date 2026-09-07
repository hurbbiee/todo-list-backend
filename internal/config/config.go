package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const webhookEncryptionKeySize = 32

type RabbitMQConfig struct {
	URL string
}

type DatabaseConfig struct {
	URL            string
	ConnectTimeout time.Duration
	MaxConnections int32
	MinConnections int32
}

type SecurityConfig struct {
	WebhookEncryptionKey []byte
}

type Config struct {
	Port         string
	FrontendURL  string
	JWTSecret    []byte
	CrossOrigins string
	Database     DatabaseConfig
	RabbitMQ     RabbitMQConfig
	Security     SecurityConfig
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

	webhookEncryptionKey, err := loadWebhookEncryptionKey()
	if err != nil {
		return Config{}, err
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
		RabbitMQ: RabbitMQConfig{
			URL: getEnv(
				"RABBITMQ_URL",
				"amqp://todo:todo-secret@localhost:5672/",
			),
		},
		Security: SecurityConfig{
			WebhookEncryptionKey: webhookEncryptionKey,
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

	if cfg.RabbitMQ.URL == "" {
		return Config{}, errors.New("RABBITMQ_URL is required")
	}

	return *cfg, nil
}

func loadWebhookEncryptionKey() ([]byte, error) {
	encodedKey := os.Getenv("WEBHOOK_ENCRYPTION_KEY")
	if encodedKey == "" {
		return nil, errors.New("WEBHOOK_ENCRYPTION_KEY is required")
	}

	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf(
			"WEBHOOK_ENCRYPTION_KEY must be valid base64: %w",
			err,
		)
	}
	if len(key) != webhookEncryptionKeySize {
		return nil, fmt.Errorf(
			"WEBHOOK_ENCRYPTION_KEY must decode to %d bytes",
			webhookEncryptionKeySize,
		)
	}

	return key, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
