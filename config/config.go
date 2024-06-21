package config

import (
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	Port string
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &AppConfig{
		Port: os.Getenv("PORT"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	return config, nil
}
