package cli

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"golang.org/x/exp/slog"
)

const logMsg = "miniserve"

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

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if err, ok := err.(error); ok {
						w.WriteHeader(http.StatusInternalServerError)
						msg := fmt.Sprintf("%s: %v", appName, err)
						logger.Error(msg,
							slog.Any("err", err),
							"trace", debug.Stack(),
						)
					}
				}
			}()
			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Info(logMsg,
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}
		return http.HandlerFunc(fn)
	}
}

// NewLogger returns a configured slog logger or nil.
func (app *App) NewLogger(w io.Writer) *slog.Logger {
	if app.NoOp() {
		return nil
	}
	return slog.New(slog.NewTextHandler(w))
}
