package cli

import (
	"io"

	"github.com/telemachus/humane"
	"golang.org/x/exp/slog"
)

// NewLogger returns a configured slog logger or nil.
func (app *App) NewLogger(w io.Writer) *slog.Logger {
	if app.NoOp() {
		return nil
	}
	return slog.New(humane.NewHandler(w))
}
