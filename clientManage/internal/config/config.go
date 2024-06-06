package app

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	Env         string
	DBUrl       string
	RabbitMQUrl string
}

func LoadConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "4000"),
		Env:  getEnv("ENV", "development"),
		// $env:DSN="postgres://postgres:12345@localhost/clientdb?sslmode=disable"
		DBUrl: getEnv("DB_URL", "postgres://postgres:12345@localhost/clientdb?sslmode=disable"),
		// $env:RABBITMQ_URL="amqp://guest:guest@localhost:5672/"
		RabbitMQUrl: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Warning: %s environment variable not set. Using default value: %s", key, defaultValue)
		return defaultValue
	}
	return value
}
