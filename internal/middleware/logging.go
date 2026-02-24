package middleware

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/afifudin23/absensi-king-royal-api/pkg/logger"
	"github.com/google/uuid"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

// RequestLogger menambahkan request_id + log setiap request
func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Ambil request id dari header atau generate
			rid := r.Header.Get("X-Request-Id")
			if rid == "" {
				rid = uuid.NewString()
			}

			// Simpan ke context
			ctx := logger.SetRequestID(r.Context(), rid)
			r = r.WithContext(ctx)

			sw := &statusWriter{ResponseWriter: w}

			next.ServeHTTP(sw, r)

			l := logger.WithRequestID(r.Context(), log)
			l.Info("http_request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", sw.status,
				"bytes", sw.bytes,
				"duration_ms", time.Since(start).Milliseconds(),
				"ip", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}
