package shutdown

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

var ErrOSSignal = errors.New("operating system signal")

func ListenSignal(ctx context.Context) error {
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-ctx.Done():
		return nil
	case <-sigquit:
		return ErrOSSignal
	}
}
