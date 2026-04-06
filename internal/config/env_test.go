package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv_PrefersDotEnvOverProcessEnv(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() err = %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("os.Chdir() err = %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	restoreEnv := snapshotEnv(
		"APP_NAME",
		"ENVIRONMENT",
		"DATABASE_URL",
		"ACCESS_KEY",
		"SERVER_BASE_URL",
		"PORT",
	)
	t.Cleanup(restoreEnv)

	// Process env has one set of values...
	mustSetenv(t, "ENVIRONMENT", "production")
	mustSetenv(t, "DATABASE_URL", "mysql://from-process-env")
	mustSetenv(t, "SERVER_BASE_URL", "https://from-process-env.example")

	// ...but .env should win.
	dotEnv := []byte(
		"ENVIRONMENT=staging\n" +
			"DATABASE_URL=mysql://from-dotenv\n" +
			"SERVER_BASE_URL=https://from-dotenv.example\n" +
			"PORT=9090\n",
	)
	if err := os.WriteFile(filepath.Join(tmpDir, ".env"), dotEnv, 0o600); err != nil {
		t.Fatalf("os.WriteFile(.env) err = %v", err)
	}

	cfg, err := LoadEnv()
	if err != nil {
		t.Fatalf("LoadEnv() err = %v", err)
	}
	if cfg.Environment != "staging" {
		t.Fatalf("Environment = %q, want %q", cfg.Environment, "staging")
	}
	if cfg.DatabaseURL != "mysql://from-dotenv" {
		t.Fatalf("DatabaseURL = %q, want %q", cfg.DatabaseURL, "mysql://from-dotenv")
	}
	if cfg.ServerBaseURL != "https://from-dotenv.example" {
		t.Fatalf("ServerBaseURL = %q, want %q", cfg.ServerBaseURL, "https://from-dotenv.example")
	}
	if cfg.Port != ":9090" {
		t.Fatalf("Port = %q, want %q", cfg.Port, ":9090")
	}
}

func mustSetenv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("os.Setenv(%q) err = %v", key, err)
	}
}

func snapshotEnv(keys ...string) func() {
	type envValue struct {
		value string
		ok    bool
	}
	original := make(map[string]envValue, len(keys))
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		original[key] = envValue{value: value, ok: ok}
	}

	return func() {
		for key, v := range original {
			if !v.ok {
				_ = os.Unsetenv(key)
				continue
			}
			_ = os.Setenv(key, v.value)
		}
	}
}

