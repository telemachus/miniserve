package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Thanks to https://blog.questionable.services/article/guide-logging-middleware-go
// for ideas and code.

// responseWriter wraps http.ResponseWriter so that we can log key information.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func Logging(logger *slog.Logger, appName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			elapsed := time.Since(start)
			status := fmt.Sprintf(
				"%d %s",
				wrapped.status,
				http.StatusText(wrapped.status),
			)
			ctx := context.TODO()
			logger.LogAttrs(
				ctx,
				slog.LevelInfo,
				fmt.Sprintf("%s: request received and response sent", appName),
				slog.Group("request",
					slog.String("host", r.Host),
					slog.String("path", r.URL.EscapedPath()),
					slog.String("proto", r.Proto),
					slog.String("method", r.Method),
				),
				slog.Group("response",
					slog.String("status", status),
					slog.Any("duration", elapsed),
				),
			)
		})
	}
}
