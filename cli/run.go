// Package cli organizes and implements a command line program.
package cli

import (
	"os"
)

const (
	exitSuccess = 0
	exitFailure = 1
	appName     = "miniserve"
	appVersion  = "v0.0.2"
	appUsage    = `usage: miniserve [-port PORT] [-d DIR]

options:
    -port PORT    Specify the port to run on (default is 8080)
    -dir DIR      Specify the directory to serve (default is ".")
    -help, -h     Show this message
    -version      Show version`
)

// Run creates an App, performs the App's tasks, and returns an exit value.
func Run(args []string) int {
	app := &App{ExitValue: exitSuccess}

	app.ParseFlags(args)
	logger := app.NewLogger(os.Stderr)
	server := app.NewServer(logger)
	app.StartAndShutdown(server, logger)

	return app.ExitValue
}
