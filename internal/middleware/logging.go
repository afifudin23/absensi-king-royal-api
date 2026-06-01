package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
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

		includePayload := env != "local" || config.GetEnv().LogDetail
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
			respContentType := strings.ToLower(c.Writer.Header().Get("Content-Type"))
			var responsePayload any
			if isBinaryContentType(respContentType) {
				responsePayload = "[BINARY_CONTENT]"
			} else {
				responsePayload = redactPayload(responseBuffer.Bytes())
			}
			fields["payload"] = map[string]any{
				"request":  requestPayload,
				"response": responsePayload,
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

		if c.Request.Method != "GET" {
			go saveActivityLog(c.Request.Method, c.Request.URL.Path, statusCode, latencyMs, userID, c.ClientIP())
		}
	}
}

func saveActivityLog(method, path string, statusCode int, latencyMs float64, userID any, ip string) {
	db := config.GetDB()
	if db == nil {
		return
	}
	repo := repository.NewActivityLogRepository(db)

	entry := &model.ActivityLog{
		Method:      method,
		Path:        path,
		Description: describeActivity(method, path),
		StatusCode:  statusCode,
		IPAddress:   ip,
		LatencyMs:   latencyMs,
	}

	if uid, ok := userID.(string); ok && strings.TrimSpace(uid) != "" {
		entry.UserID = &uid
		var user model.User
		if err := db.Select("full_name").Where("id = ?", uid).First(&user).Error; err == nil {
			entry.UserName = &user.FullName
		}
	}

	_ = repo.Create(context.Background(), entry)
}

func describeActivity(method, path string) string {
	base := path
	if idx := strings.Index(path, "/api/v1"); idx >= 0 {
		base = path[idx+7:]
	}

	// Special cases — cek dulu sebelum rules umum
	if method == "POST" && strings.HasSuffix(base, "/send") {
		return "Kirim slip gaji via email"
	}
	if method == "POST" && strings.HasSuffix(base, "/reset-password") {
		return "Reset password karyawan"
	}
	if strings.Contains(base, "/profile-picture") {
		if method == "PUT" {
			return "Ganti foto profil"
		}
		if method == "DELETE" {
			return "Hapus foto profil"
		}
	}

	type rule struct {
		method  string
		pattern string
		desc    string
	}
	rules := []rule{
		// Auth
		{"POST", "/auth/login", "Login"},
		{"POST", "/auth/logout", "Logout"},
		{"POST", "/auth/forgot-password", "Minta reset password"},
		{"POST", "/auth/reset-password", "Reset password via OTP"},
		// Profil sendiri
		{"GET", "/users/me", "Lihat profil saya"},
		{"PUT", "/users/me/password", "Ganti password"},
		{"PUT", "/users/me", "Edit profil saya"},
		// Manajemen karyawan
		{"GET", "/users", "Lihat daftar karyawan"},
		{"POST", "/users", "Tambah karyawan"},
		{"GET", "/users/", "Lihat detail karyawan"},
		{"PUT", "/users/", "Edit data karyawan"},
		{"DELETE", "/users/", "Hapus karyawan"},
		// Absensi
		{"POST", "/attendances/check-in", "Absen masuk"},
		{"POST", "/attendances/check-out", "Absen pulang"},
		{"GET", "/attendances/recap", "Lihat rekap absensi"},
		{"GET", "/attendances/logs", "Lihat log absensi"},
		{"GET", "/attendances", "Lihat daftar absensi"},
		{"POST", "/attendances", "Tambah data absensi"},
		{"PUT", "/attendances/", "Edit data absensi harian"},
		{"DELETE", "/attendances/", "Hapus data absensi"},
		// Pengajuan izin
		{"GET", "/attendance-requests", "Lihat daftar pengajuan"},
		{"POST", "/attendance-requests", "Ajukan izin / cuti / lembur"},
		{"PUT", "/attendance-requests/", "Update status pengajuan"},
		{"DELETE", "/attendance-requests/", "Hapus pengajuan"},
		// Payroll
		{"POST", "/payrolls/generate-all", "Generate semua slip gaji"},
		{"POST", "/payrolls/generate/", "Generate slip gaji"},
		{"GET", "/payrolls/me", "Lihat slip gaji saya"},
		{"GET", "/payrolls", "Lihat daftar slip gaji"},
		{"PUT", "/payrolls/", "Edit komponen slip gaji"},
		// File
		{"POST", "/files", "Upload file"},
		{"DELETE", "/files/", "Hapus file"},
		{"GET", "/files/", "Lihat file"},
		// Log
		{"GET", "/activity-logs", "Lihat log aktivitas"},
	}

	for _, r := range rules {
		if r.method != method {
			continue
		}
		if base == r.pattern || strings.HasPrefix(base, r.pattern) {
			return r.desc
		}
	}

	return method + " " + base
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

func isBinaryContentType(contentType string) bool {
	binaryPrefixes := []string{"image/", "video/", "audio/", "application/octet-stream", "application/pdf"}
	for _, prefix := range binaryPrefixes {
		if strings.HasPrefix(contentType, prefix) {
			return true
		}
	}
	return false
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
