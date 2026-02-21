package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppName     string
	DatabaseURL string
	Port        string
}

func LoadEnv() (*EnvConfig, error) {
	_ = godotenv.Load()

	appName := strings.TrimSpace(os.Getenv("APP_NAME"))
	if appName == "" {
		appName = "absensi-king-royal-api"
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return &EnvConfig{
		AppName:     appName,
		DatabaseURL: databaseURL,
		Port:        port,
	}, nil
}
