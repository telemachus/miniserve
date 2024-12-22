package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"
)

// PanicRecovery wraps an http.Handler with a middleware to handle panics more gracefully.
func PanicRecovery(logger *slog.Logger, appName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					http.Error(
						w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
					errorLog(logger, appName, fmt.Errorf("%s", err))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func errorLog(l *slog.Logger, appName string, err error) {
	ctx := context.TODO()
	if !l.Enabled(ctx, slog.LevelError) {
		return
	}
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	var pcs [1]uintptr
	// Skip runtime.Callers, errorLog, and PanicRecovery to reach the true
	// caller.
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(
		time.Now(),
		slog.LevelError,
		fmt.Sprintf("%s: internal server error", appName),
		pcs[0],
	)
	r.AddAttrs(slog.String("err", trace))
	_ = l.Handler().Handle(ctx, r)
}
