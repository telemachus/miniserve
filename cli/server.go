package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	kitlog "github.com/go-kit/log"
)

const (
	clientHeaderTimeout = 60
	shutdownTimeout     = 5
)

func (app *App) NewServer(l kitlog.Logger) *http.Server {
	if app.NoOp() {
		return nil
	}

	addr := fmt.Sprintf(":%s", app.Port)
	fs := http.FileServer(wrappedDir{http.Dir(app.Dir)})
	middleware := loggingMiddleware(l)
	fsWithLogging := middleware(fs)

	return &http.Server{
		ReadHeaderTimeout: clientHeaderTimeout * time.Second,
		Addr:              addr,
		Handler:           fsWithLogging,
	}
}

func (app *App) StartAndShutdown(s *http.Server, l kitlog.Logger) {
	if app.NoOp() {
		return
	}

	go app.start(s, l)

	stopCh, closeCh := newChannel()
	defer closeCh()
	l.Log("level", "info", "msg", <-stopCh)

	app.shutdown(context.Background(), s, l)
}

func (app *App) start(s *http.Server, l kitlog.Logger) {
	if app.NoOp() {
		return
	}

	l.Log(
		"level", "info",
		"msg", "starting miniserve",
		"addr", s.Addr,
	)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Log(
			"level", "error",
			"msg", err,
		)

		app.ExitValue = exitFailure

		return
	}

	l.Log(
		"level", "info",
		"msg", "attempting graceful shutdown for miniserve",
	)
}

func (app *App) shutdown(ctx context.Context, s *http.Server, l kitlog.Logger) {
	if app.NoOp() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		msg := fmt.Sprintf("miniserve failed to shut down cleanly: %v", err)
		l.Log(
			"level", "error",
			"msg", msg,
		)

		app.ExitValue = exitFailure

		return
	}

	l.Log(
		"level", "info",
		"msg", "miniserve successfully shut down",
	)
}
