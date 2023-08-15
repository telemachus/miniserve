package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/justinas/alice"
	"github.com/telemachus/miniserve/internal/middleware"
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
	handlerChain := alice.New(
		middleware.PanicRecovery(logger, appName),
		middleware.Logging(logger, appName),
	).Then(fs)
	return &http.Server{
		ReadHeaderTimeout: clientHeaderTimeout * time.Second,
		Addr:              addr,
		Handler:           handlerChain,
	}
}

func (app *App) StartAndShutdown(server *http.Server, logger *slog.Logger) {
	if app.NoOp() {
		return
	}
	app.TrapSignals(os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer app.CloseTrap()
	ctx := context.TODO()
	go app.start(server, logger)
	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		fmt.Sprintf("%s: caught signal", appName),
		slog.Any("signal", <-app.Trap),
	)
	app.shutdown(context.Background(), server, logger)
}

func (app *App) start(server *http.Server, logger *slog.Logger) {
	if app.NoOp() {
		return
	}
	ctx := context.TODO()
	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		fmt.Sprintf("%s: server started", appName),
		slog.String("addr", server.Addr),
	)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			fmt.Sprintf("%s: server failed to start", appName),
			slog.Any("err", err),
		)
		app.ExitValue = exitFailure
		// Don't hang the terminal if the server never starts.
		// Alternatively, I could simply call panic(err) here.
		app.Trap <- syscall.SIGINT
		return
	}
	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		fmt.Sprintf("%s: initiated graceful shutdown", appName),
	)
}

func (app *App) shutdown(ctx context.Context, server *http.Server, logger *slog.Logger) {
	if app.NoOp() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			fmt.Sprintf("%s: failed graceful shutdown", appName),
			slog.Any("err", err),
		)
		app.ExitValue = exitFailure
		// Don't hang the terminal if the server fails to shutdown cleanly.
		// Alternatively, I could simply call panic(err) here.
		app.Trap <- syscall.SIGINT
		return
	}
	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		fmt.Sprintf("%s: finished graceful shutdown", appName),
	)
}
