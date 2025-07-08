package runal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/input"
	"github.com/charmbracelet/x/term"
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

	resize := listenForResize()

	oldState, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(os.Stdin.Fd(), oldState)

	reader, err := input.NewReader(os.Stdin, os.Getenv("TERM"), 0)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	inputEvents := make(chan input.Event, 2048)

	go func() {
		for {
			events, _ := reader.ReadEvents()
			for _, ev := range events {
				inputEvents <- ev
			}
		}
	}()

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
		resetCursorPosition()
		clearScreen()
		showCursor()
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
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
			case event := <-inputEvents:
				switch e := event.(type) {
				case input.KeyEvent:
					switch e.String() {
					case "ctrl+c":
						exit()
						if done != nil {
							done <- os.Interrupt
						}
						fmt.Println("COUCOU")
						return
					default:
						if onKey != nil {
							onKey(c, e.String())
						}
					}
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
