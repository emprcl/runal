package runal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/charmbracelet/x/input"
)

const (
	defaultFPS = 30
)

func Run(ctx context.Context, setup, draw func(c *Canvas), onKey func(c *Canvas, e KeyEvent), onMouse func(c *Canvas, e MouseEvent)) {
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	Start(ctx, nil, setup, draw, onKey, onMouse).Wait()
}

func Start(ctx context.Context, done chan os.Signal, setup, draw func(c *Canvas), onKey func(c *Canvas, e KeyEvent), onMouse func(c *Canvas, e MouseEvent)) *sync.WaitGroup {
	w, h := termSize()
	c := newCanvas(w, h)

	resize := listenForResize()
	inputEvents := listenForInputEvents(ctx)

	enterAltScreen()
	enableMouse()

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
		disableMouse()
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
				case input.MouseClickEvent:
					if onMouse != nil {
						mouse := e.Mouse()
						onMouse(c, MouseEvent{X: mouse.X, Y: mouse.Y})
					}
				case input.KeyEvent:
					switch e.String() {
					case "ctrl+c":
						exit()
						if done != nil {
							done <- os.Interrupt
						}
						return
					default:
						if onKey != nil {
							onKey(c, KeyEvent{
								Key:  e.Key().String(),
								Code: int(e.Key().Code),
							})
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
