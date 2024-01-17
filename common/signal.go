package common

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitSignal() <-chan os.Signal {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	return osSignals
}
