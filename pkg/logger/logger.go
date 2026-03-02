package logger

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

type ctxKey string

const (
	requestIDKey ctxKey = "request_id"
	userIDKey    ctxKey = "user_id"
)

type outputLog struct {
	Timestamp   string         `json:"timestamp"`
	Level       string         `json:"level"`
	Service     string         `json:"service"`
	Environment string         `json:"environment"`
	RequestID   string         `json:"request_id,omitempty"`
	UserID      any            `json:"user_id,omitempty"`
	Message     string         `json:"message"`
	LoggerName  string         `json:"logger_name"`
	HTTP        any            `json:"http,omitempty"`
	Payload     any            `json:"payload,omitempty"`
	Error       any            `json:"error,omitempty"`
	Event       any            `json:"event,omitempty"`
	Context     map[string]any `json:"context,omitempty"`
}

var (
	mu          sync.RWMutex
	serviceName = "absensi-king-royal-api"
	environment = "local"
)

func Configure(service string, env string) {
	mu.Lock()
	defer mu.Unlock()

	if strings.TrimSpace(service) != "" {
		serviceName = strings.TrimSpace(service)
	}

	normalizedEnv := strings.ToLower(strings.TrimSpace(env))
	if normalizedEnv == "" {
		normalizedEnv = "local"
	}
	environment = normalizedEnv
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, requestIDKey, strings.TrimSpace(requestID))
}

func WithUserID(ctx context.Context, userID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, userIDKey, strings.TrimSpace(userID))
}

func Info(ctx context.Context, loggerName, message string, fields map[string]any) {
	log(ctx, "INFO", loggerName, message, fields)
}

func Warn(ctx context.Context, loggerName, message string, fields map[string]any) {
	log(ctx, "WARNING", loggerName, message, fields)
}

func Error(ctx context.Context, loggerName, message string, fields map[string]any) {
	log(ctx, "ERROR", loggerName, message, fields)
}

func log(ctx context.Context, level, loggerName, message string, fields map[string]any) {
	svc, env := currentConfig()
	entry := outputLog{
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
		Level:       level,
		Service:     svc,
		Environment: env,
		Message:     message,
		LoggerName:  loggerName,
	}

	if ctx != nil {
		if v, ok := ctx.Value(requestIDKey).(string); ok && strings.TrimSpace(v) != "" {
			entry.RequestID = v
		}
		if v, ok := ctx.Value(userIDKey).(string); ok && strings.TrimSpace(v) != "" {
			entry.UserID = v
		}
	}

	if fields != nil {
		// Reserved structured fields
		if v, ok := fields["http"]; ok {
			entry.HTTP = v
			delete(fields, "http")
		}
		if v, ok := fields["payload"]; ok && env != "local" {
			entry.Payload = v
			delete(fields, "payload")
		}
		if v, ok := fields["error"]; ok {
			entry.Error = v
			delete(fields, "error")
		}
		if v, ok := fields["event"]; ok {
			entry.Event = v
			delete(fields, "event")
		}
		if v, ok := fields["request_id"]; ok {
			entry.RequestID = anyToString(v)
			delete(fields, "request_id")
		}
		if v, ok := fields["user_id"]; ok {
			entry.UserID = v
			delete(fields, "user_id")
		}
		if len(fields) > 0 {
			entry.Context = fields
		}
	}

	write(entry, env)
}

func currentConfig() (string, string) {
	mu.RLock()
	defer mu.RUnlock()
	return serviceName, environment
}

func write(entry outputLog, env string) {
	var out []byte
	var err error

	if env == "local" {
		out, err = json.MarshalIndent(entry, "", "  ")
	} else {
		out, err = json.Marshal(entry)
	}
	if err != nil {
		return
	}
	_, _ = os.Stdout.Write(append(out, '\n'))
}

func anyToString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return strings.Trim(string(b), "\"")
	}
}
