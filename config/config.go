package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	Port  string
	DBUrl string
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &AppConfig{
		Port:  os.Getenv("PORT"),
		DBUrl: os.Getenv("DB_URL"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	if config.DBUrl == "" {
		return nil, errors.New("DB_URL is not set")
	}

	return config, nil
}
