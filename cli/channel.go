package cli

import (
	"os"
	"os/signal"
	"syscall"
)

func newChannel() (chan os.Signal, func()) {
	done := make(chan os.Signal, 1)
	signal.Notify(
		done,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	return done, func() {
		close(done)
	}
}
