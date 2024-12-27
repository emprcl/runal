package runal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context, setup, draw func(c *Canvas), opts ...option) {
	config := newOptions()
	for _, opt := range opts {
		opt(config)
	}
	w, h := termSize()
	c := newCanvas(w, h)
	setup(c)

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	tick := time.Tick(config.frameDuration)

	enterAltScreen()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-resize:
				w, h := termSize()
				c.resize(w, h)
				clearScreen()
			case <-tick:
				resetCursorPosition()
				draw(c)
				c.render()
			}
		}
	}()
}
