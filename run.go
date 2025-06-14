package runal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"
)

const (
	defaultFPS = 30
)

func Run(ctx context.Context, setup, draw func(c *Canvas)) {
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	Start(ctx, setup, draw, nil).Wait()
}

func Start(ctx context.Context, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) *sync.WaitGroup {
	ctx, cancel := context.WithCancel(ctx)
	w, h := termSize()
	c := newCanvas(w, h)

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)

	enterAltScreen()

	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	go func() {
		buf := make([]byte, 1)
		for {
			os.Stdin.Read(buf)
			if buf[0] == 3 { // ctrl+c
				cancel()
				return
			}
			if onKey != nil {
				onKey(c, string(buf))
			}
		}
	}()

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
