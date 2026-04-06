package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppName       string
	Environment   string
	DatabaseURL   string
	AccessKey     string
	Port          string
	ServerBaseURL string
}

func LoadEnv() (*EnvConfig, error) {
	// Prioritize values from .env over existing process environment variables.
	// If .env doesn't exist, fall back to the current process environment.
	_ = godotenv.Overload()

	appName := strings.TrimSpace(os.Getenv("APP_NAME"))
	if appName == "" {
		appName = "absensi-king-royal-api"
	}

	environment := strings.ToLower(strings.TrimSpace(os.Getenv("ENVIRONMENT")))
	if environment == "" {
		return nil, fmt.Errorf("ENVIRONMENT is required")
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	accessKey := strings.TrimSpace(os.Getenv("ACCESS_KEY"))
	if accessKey == "" {
		accessKey = "dev-secret-change-me"
	}

	serverBaseURL := strings.TrimSpace(os.Getenv("SERVER_BASE_URL"))
	if serverBaseURL == "" {
		return nil, fmt.Errorf("SERVER_BASE_URL is required")
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return &EnvConfig{
		AppName:       appName,
		Environment:   environment,
		DatabaseURL:   databaseURL,
		AccessKey:     accessKey,
		Port:          port,
		ServerBaseURL: serverBaseURL,
	}, nil
}
