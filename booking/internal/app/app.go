package app

import (
	"os"
)

type Config struct {
	DBUrl       string
	RabbitMQUrl string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		DBUrl:       os.Getenv("DB_URL"),
		RabbitMQUrl: os.Getenv("RABBITMQ_URL"),
		Port:        os.Getenv("PORT"),
	}
}
