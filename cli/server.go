package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	kitlog "github.com/go-kit/log"
)

const (
	clientHeaderTimeout = 60
	shutdownTimeout     = 5
)

func (app *App) NewServer(logger kitlog.Logger) *http.Server {
	if app.NoOp() {
		return nil
	}

	addr := fmt.Sprintf(":%s", app.Port)
	fs := http.FileServer(WrappedDir{http.Dir(app.Dir)})
	middleware := loggingMiddleware(logger)
	fsWithLogging := middleware(fs)

	return &http.Server{
		ReadHeaderTimeout: clientHeaderTimeout * time.Second,
		Addr:              addr,
		Handler:           fsWithLogging,
	}
}

func (app *App) StartAndShutdown(server *http.Server, logger kitlog.Logger) {
	if app.NoOp() {
		return
	}

	app.TrapSignals(os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer app.CloseTrap()

	go app.start(server, logger)

	logger.Log(
		"level", "info",
		"msg", <-app.Trap,
	)

	app.shutdown(context.Background(), server, logger)
}

func (app *App) start(server *http.Server, logger kitlog.Logger) {
	if app.NoOp() {
		return
	}

	logger.Log(
		"level", "info",
		"msg", fmt.Sprintf("starting %s", appName),
		"addr", server.Addr,
	)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log(
			"level", "error",
			"msg", err,
		)

		app.ExitValue = exitFailure

		// Don't hang the terminal if the server never starts.
		// Alternatively, I could simply call panic(err) here.
		app.Trap <- syscall.SIGINT

		return
	}

	logger.Log(
		"level", "info",
		"msg", fmt.Sprintf("attempting graceful shutdown for %s", appName),
	)
}

func (app *App) shutdown(ctx context.Context, server *http.Server, logger kitlog.Logger) {
	if app.NoOp() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		msg := fmt.Sprintf("%s failed to shut down cleanly: %v", appName, err)
		logger.Log(
			"level", "error",
			"msg", msg,
		)

		app.ExitValue = exitFailure

		// Don't hang the terminal if the server fails to shutdown cleanly.
		// Alternatively, I could simply call panic(err) here.
		app.Trap <- syscall.SIGINT

		return
	}

	logger.Log(
		"level", "info",
		"msg", fmt.Sprintf("%s successfully shut down", appName),
	)
}
