package cli

import (
	"os"
)

// App stores information about the application's state.
type App struct {
	Trap          chan os.Signal
	Port          string
	Dir           string
	ExitValue     int
	HelpWanted    bool
	VersionWanted bool
}

// NoOp determines whether an App should bail out.
func (app *App) NoOp() bool {
	return app.ExitValue != exitSuccess || app.HelpWanted || app.VersionWanted
}
