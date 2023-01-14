package cli

import (
	"os"
	"os/signal"
)

// I don't need to check app.NoOp for these functions because they appear
// only in a function that has already made such a check.

// TrapSignals sets app.Trap to watch for a variable number of signals.
func (app *App) TrapSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)

	app.Trap = ch
}

// CloseTrap closes the os.Signal channel.
func (app *App) CloseTrap() {
	// This should never happen, but let's be extra careful.
	if app.Trap == nil {
		return
	}

	close(app.Trap)
}
