package runal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run(ctx context.Context, setup, draw func(c *Canvas), opts ...option) *sync.WaitGroup {
	config := newOptions()
	for _, opt := range opts {
		opt(config)
	}
	w, h := termSize()
	c := newCanvas(w, h)

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	tick := time.Tick(config.frameDuration)
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
			case <-tick:
				resetCursorPosition()
				update <- struct{}{}
				c.render()
			}
		}
	}()

	return &wg
}
