package runal

import (
	"context"
	"log"
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

func Run(ctx context.Context, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) {
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	Start(ctx, setup, draw, onKey).Wait()
}

func Start(ctx context.Context, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) *sync.WaitGroup {
	w, h := termSize()
	c := newCanvas(w, h)

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)

	input := inputChannel(ctx)

	enterAltScreen()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	setup(c)
	render := func() {
		resetCursorPosition()
		draw(c)
		c.render()
	}
	render()

	exit := func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		clearScreen()
		resetCursorPosition()
		showCursor()
	}

	ticker := time.NewTicker(newFramerate(defaultFPS))
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
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
				case "render":
					render()
				}
			case char := <-input:
				// ctrl+c
				if char == 3 {
					exit()
					return
				}
				if onKey != nil {
					onKey(c, string(char))
				}
			case <-ticker.C:
				render()
			}
		}
	}()

	return &wg
}

func inputChannel(ctx context.Context) <-chan byte {
	ch := make(chan byte)
	go func() {
		defer close(ch)
		buf := make([]byte, 1)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				continue
			}
			ch <- buf[0]
		}
	}()
	return ch
}

func newFramerate(fps int) time.Duration {
	return time.Second / time.Duration(fps)
}
