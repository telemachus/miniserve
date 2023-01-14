package cli

import (
	"io"
	stdlog "log"
	"net/http"
	"runtime/debug"
	"time"

	kitlog "github.com/go-kit/log"
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

func loggingMiddleware(logger kitlog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Log(
						"level", "err",
						"msg", err,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Log(
				"level", "info",
				"msg", "http logline",
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}

// NewLogger returns a configured go-kit logger.
func (app *App) NewLogger(w io.Writer) kitlog.Logger {
	if app.NoOp() {
		return nil
	}

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(w))
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger))
	logger = kitlog.With(logger,
		"ts", kitlog.DefaultTimestamp,
	)

	return logger
}
