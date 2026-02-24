package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"
)

// ctxKey untuk menyimpan request id di context
type ctxKey string

const RequestIDKey ctxKey = "request_id"

type Config struct {
	AppName   string // contoh: "absensi-king-royal-api"
	Env       string // "dev" | "prod"
	Level     string // "debug"|"info"|"warn"|"error"
	AddSource bool   // tampilkan file:line
}

// New membuat logger slog dengan handler sesuai env
func New(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Biar field time jadi rapi (RFC3339)
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.UTC().Format(time.RFC3339Nano))
				}
			}
			return a
		},
	}

	var handler slog.Handler
	if strings.ToLower(cfg.Env) == "dev" {
		// Text handler lebih enak dibaca di local dev
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		// JSON handler cocok untuk production (ELK/Loki/Cloud Logging)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	base := slog.New(handler).With(
		"app", cfg.AppName,
		"env", cfg.Env,
	)

	return base
}

// WithRequestID mengembalikan logger yang sudah punya request_id dari context
func WithRequestID(ctx context.Context, log *slog.Logger) *slog.Logger {
	if ctx == nil || log == nil {
		return log
	}
	if rid, ok := ctx.Value(RequestIDKey).(string); ok && rid != "" {
		return log.With("request_id", rid)
	}
	return log
}

// SetRequestID simpan request ID ke context
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func parseLevel(lvl string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(lvl)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
