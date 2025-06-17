package runal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	defaultFPS = 30
)

func Run(ctx context.Context, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) {
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	Start(ctx, nil, setup, draw, onKey).Wait()
}

func Start(ctx context.Context, done chan os.Signal, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) *sync.WaitGroup {
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

	exit := func() {
		ticker.Stop()
		clearScreen()
		resetCursorPosition()
		showCursor()
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			_ = keyboard.Close()
		}()
		keyEvent, _ := keyboard.GetKeys(1)
		for {
			select {
			case <-ctx.Done():
				exit()
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
				case "start":
					ticker.Reset(newFramerate(defaultFPS))
				case "render":
					render()
				}
			case event := <-keyEvent:
				// ctrl+c
				if event.Key == keyboard.KeyCtrlC {
					exit()
					if done != nil {
						done <- os.Interrupt
					}
					return
				}
				// NOTE: keyboard package has a small bug on
				// space key not filling the Rune attribute.
				if event.Key == keyboard.KeySpace {
					event.Rune = ' '
				}
				if onKey != nil {
					onKey(c, string(event.Rune))
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
