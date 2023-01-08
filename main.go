// Package main is a minimal driver for miniserve, a small server for testing.
package main

import (
	"os"

	"github.com/telemachus/miniserve/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
