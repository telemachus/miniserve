package cli

import (
	"io"
	"log/slog"

	"github.com/telemachus/humane"
)

// NewLogger returns a configured slog logger or nil.
func (app *App) NewLogger(w io.Writer) *slog.Logger {
	if app.NoOp() {
		return nil
	}
	return slog.New(humane.NewHandler(w, nil))
}
