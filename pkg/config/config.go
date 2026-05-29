package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GiphyKey string
	Port     string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		GiphyKey: os.Getenv("GIPHY_API_KEY"),
		Port:     port,
	}, nil
}
