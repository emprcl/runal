package runal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	defaultFPS = 20
)

func Run(ctx context.Context, setup, draw func(c *Canvas)) *sync.WaitGroup {
	w, h := termSize()
	c := newCanvas(w, h)

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	ticker := time.NewTicker(newFramerate(defaultFPS))

	update := make(chan struct{})

	enterAltScreen()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		setup(c)
		for {
			select {
			case _, ok := <-update:
				if !ok {
					return
				}
				draw(c)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				close(update)
				clearScreen()
				resetCursorPosition()
				showCursor()
				return
			case <-resize:
				clearScreen()
				if !c.autoResize {
					continue
				}
				w, h := termSize()
				c.resize(w, h)
			case event := <-c.bus:
				switch event.name {
				case "fps":
					ticker.Reset(newFramerate(event.value))
				}
			case <-ticker.C:
				resetCursorPosition()
				update <- struct{}{}
				c.render()
			}
		}
	}()

	return &wg
}

func newFramerate(fps int) time.Duration {
	return time.Second / time.Duration(fps)
}
