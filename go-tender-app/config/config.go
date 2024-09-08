package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	PostgresURL   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return &Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		PostgresURL:   os.Getenv("POSTGRES_CONN"),
	}
}
