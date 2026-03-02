package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func StructuredLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		env := getEnvironment()

		if shouldSkipLog(c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()
		requestID := strings.TrimSpace(c.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = "req_" + uuid.NewString()
		}
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		reqCtx := logger.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(reqCtx)

		includePayload := env != "local"
		var requestPayload any
		responseBuffer := bytes.NewBuffer(nil)
		if includePayload {
			requestPayload = captureRequestSummary(c)
			blw := &bodyLogWriter{
				ResponseWriter: c.Writer,
				body:           responseBuffer,
			}
			c.Writer = blw
		}

		c.Next()

		statusCode := c.Writer.Status()
		latencyMs := float64(time.Since(start).Microseconds()) / 1000.0
		logLevel := resolveLogLevel(statusCode)

		var userID any
		if v, ok := c.Get("uid"); ok {
			userID = v
			if uidStr, ok := v.(string); ok && strings.TrimSpace(uidStr) != "" {
				reqCtx = logger.WithUserID(c.Request.Context(), uidStr)
				c.Request = c.Request.WithContext(reqCtx)
			}
		}

		fields := map[string]any{
			"user_id": userID,
			"http": map[string]any{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"status":     statusCode,
				"latency_ms": latencyMs,
				"ip":         c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			},
		}

		if includePayload {
			fields["payload"] = map[string]any{
				"request":  requestPayload,
				"response": redactPayload(responseBuffer.Bytes()),
			}
		}

		switch logLevel {
		case "ERROR":
			logger.Error(c.Request.Context(), "http.middleware", "HTTP Request Processed", fields)
		case "WARNING":
			logger.Warn(c.Request.Context(), "http.middleware", "HTTP Request Processed", fields)
		default:
			logger.Info(c.Request.Context(), "http.middleware", "HTTP Request Processed", fields)
		}
	}
}

func getEnvironment() string {
	env := config.GetEnv()
	if env == nil || strings.TrimSpace(env.Environment) == "" {
		return "local"
	}
	return strings.ToLower(strings.TrimSpace(env.Environment))
}

func shouldSkipLog(path string) bool {
	switch path {
	case "/metrics", "/docs", "/redoc", "/openapi.json", "/api/v1/health", "/api/v1":
		return true
	default:
		return false
	}
}

func resolveLogLevel(statusCode int) string {
	if statusCode >= 500 {
		return "ERROR"
	}
	if statusCode >= 400 {
		return "WARNING"
	}
	return "INFO"
}

func captureRequestSummary(c *gin.Context) any {
	contentType := strings.ToLower(c.GetHeader("Content-Type"))
	switch {
	case strings.Contains(contentType, "application/json"):
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return "[ERROR_READING_JSON_BODY]"
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		return redactPayload(body)
	case strings.Contains(contentType, "multipart/form-data"):
		return "[multipart/form-data]"
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		return "[form-urlencoded]"
	case contentType == "":
		return "[no-content-type]"
	default:
		return "[" + contentType + "]"
	}
}

func redactPayload(payload any) any {
	sensitiveKeys := map[string]struct{}{
		"password":         {},
		"current_password": {},
		"new_password":     {},
		"token":            {},
		"secret":           {},
		"authorization":    {},
		"access_token":     {},
		"refresh_token":    {},
	}

	var process func(any) any
	process = func(v any) any {
		switch value := v.(type) {
		case map[string]any:
			out := make(map[string]any, len(value))
			for k, val := range value {
				if _, ok := sensitiveKeys[strings.ToLower(k)]; ok {
					out[k] = "[REDACTED]"
					continue
				}
				out[k] = process(val)
			}
			return out
		case []any:
			if len(value) >= 10 {
				return "[Array(" + strconv.Itoa(len(value)) + ")]"
			}
			out := make([]any, 0, len(value))
			for _, item := range value {
				out = append(out, process(item))
			}
			return out
		default:
			return value
		}
	}

	switch v := payload.(type) {
	case []byte:
		raw := strings.TrimSpace(string(v))
		if raw == "" {
			return nil
		}
		var decoded any
		if err := json.Unmarshal(v, &decoded); err == nil {
			return process(decoded)
		}
		if len(raw) > 120 {
			return "[TEXT_CONTENT: " + raw[:120] + "...]"
		}
		return raw
	case string:
		raw := strings.TrimSpace(v)
		if raw == "" {
			return nil
		}
		var decoded any
		if err := json.Unmarshal([]byte(raw), &decoded); err == nil {
			return process(decoded)
		}
		if len(raw) > 120 {
			return "[TEXT_CONTENT: " + raw[:120] + "...]"
		}
		return raw
	default:
		return process(v)
	}
}
