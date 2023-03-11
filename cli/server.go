package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"golang.org/x/exp/slog"
)

const (
	clientHeaderTimeout = 60
	shutdownTimeout     = 5
)

func (app *App) NewServer(logger *slog.Logger) *http.Server {
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

func (app *App) StartAndShutdown(server *http.Server, logger *slog.Logger) {
	if app.NoOp() {
		return
	}
	app.TrapSignals(os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer app.CloseTrap()
	go app.start(server, logger)
	logger.Info(logMsg,
		"msg", <-app.Trap,
	)
	app.shutdown(context.Background(), server, logger)
}

func (app *App) start(server *http.Server, logger *slog.Logger) {
	if app.NoOp() {
		return
	}
	logger.Info(fmt.Sprintf("starting %s", appName),
		"addr", server.Addr,
	)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(logMsg, slog.Any("err", err))
		app.ExitValue = exitFailure
		// Don't hang the terminal if the server never starts.
		// Alternatively, I could simply call panic(err) here.
		app.Trap <- syscall.SIGINT
		return
	}
	logger.Info(fmt.Sprintf("attempting graceful shutdown for %s", appName))
}

func (app *App) shutdown(ctx context.Context, server *http.Server, logger *slog.Logger) {
	if app.NoOp() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		errMsg := fmt.Errorf("%s failed to shut down cleanly: %w", appName, err)
		logger.Error(logMsg, slog.Any("err", errMsg))
		app.ExitValue = exitFailure
		// Don't hang the terminal if the server fails to shutdown cleanly.
		// Alternatively, I could simply call panic(err) here.
		app.Trap <- syscall.SIGINT
		return
	}
	logger.Info(fmt.Sprintf("%s successfully shut down", appName))
}
