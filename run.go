package runal

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/x/input"
)

const (
	defaultFPS = 30
)

func Run(ctx context.Context, setup, draw func(c *Canvas), onKey func(c *Canvas, e KeyEvent), onMouse func(c *Canvas, e MouseEvent)) {
	Start(ctx, nil, setup, draw, onKey, onMouse).Wait()
}

func Start(ctx context.Context, done chan struct{}, setup, draw func(c *Canvas), onKey func(c *Canvas, e KeyEvent), onMouse func(c *Canvas, e MouseEvent)) *sync.WaitGroup {
	ctx, cancel := context.WithCancel(ctx)
	w, h := termSize()
	c := newCanvas(w, h)
	wg := sync.WaitGroup{}

	ticker := time.NewTicker(newFramerate(defaultFPS))

	exit := func() {
		ticker.Stop()
		resetCursorPosition()
		clearScreen()
		showCursor()
		disableMouse()
	}

	defer func() {
		if r := recover(); r != nil {
			exit()
			panic(r)
		}
	}()

	resize := listenForResize()
	inputEvents := listenForInputEvents(ctx, &wg)

	enterAltScreen()
	enableMouse()

	setup(c)
	render := func() {
		resetCursorPosition()
		draw(c)
		c.render()
	}
	render()

	wg.Add(1)
	go func() {
		defer func() {
			exit()
			wg.Done()
		}()
		for {
			select {
			case <-ctx.Done():
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
				case input.MouseMotionEvent:
					c.MouseX = e.X
					c.MouseY = e.Y
				case input.MouseClickEvent:
					if onMouse != nil {
						onMouse(c, MouseEvent{
							X:      e.X,
							Y:      e.Y,
							Button: e.Button.String(),
						})
					}
				case input.KeyEvent:
					switch e.String() {
					case "ctrl+c":
						if done != nil {
							done <- struct{}{}
						}
						cancel()
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
