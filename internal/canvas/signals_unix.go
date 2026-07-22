//go:build !windows

package canvas

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func listenForResize(ctx context.Context, wg *sync.WaitGroup) <-chan struct{} {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGWINCH)

	resize := make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer signal.Stop(signals)
		for {
			select {
			case <-ctx.Done():
				return
			case <-signals:
				notifyResize(resize)
			}
		}
	}()
	return resize
}
