package cli

import (
	"os"
)

// App stores information about the application's state.
type App struct {
	ExitValue     int
	Trap          chan os.Signal
	HelpWanted    bool
	VersionWanted bool
	Port          string
	Dir           string
}

// NoOp determines whether an App should bail out.
func (app *App) NoOp() bool {
	return app.ExitValue != exitSuccess || app.HelpWanted || app.VersionWanted
}
