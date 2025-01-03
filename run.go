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

	enterAltScreen()

	setup(c)
	render := func() {
		resetCursorPosition()
		draw(c)
		c.render()
	}
	render()

	ticker := time.NewTicker(newFramerate(defaultFPS))
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				clearScreen()
				resetCursorPosition()
				showCursor()
				return
			case <-resize:
				clearScreen()
				w, h := termSize()
				c.termWidth = w
				c.termHeight = h
				if c.autoResize {
					c.resize(w, h)
				}
				render()
			case event := <-c.bus:
				switch event.name {
				case "fps":
					ticker.Reset(newFramerate(event.value))
				case "stop":
					ticker.Stop()
				}
			case <-ticker.C:
				render()
			}
		}
	}()

	return &wg
}

func newFramerate(fps int) time.Duration {
	return time.Second / time.Duration(fps)
}
