//go:build windows

package canvas

import (
	"context"
	"sync"
	"time"
)

const resizePollInterval = 250 * time.Millisecond

// listenForResize polls the terminal size, since Windows has no SIGWINCH.
func listenForResize(ctx context.Context, wg *sync.WaitGroup) <-chan struct{} {
	resize := make(chan struct{}, 1)

	lastWidth, lastHeight, err := tryTermSize()
	if err != nil {
		return resize
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(resizePollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w, h, err := tryTermSize()
				if err != nil || (w == lastWidth && h == lastHeight) {
					continue
				}
				lastWidth, lastHeight = w, h
				notifyResize(resize)
			}
		}
	}()
	return resize
}
