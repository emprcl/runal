//go:build !js

package canvas

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/input"
	charmterm "github.com/charmbracelet/x/term"
	"golang.org/x/term"
)

func Run(ctx context.Context, setup, draw func(c *Canvas), opts ...CallbackOption) {
	Start(ctx, nil, setup, draw, opts...).Wait()
}

func Start(ctx context.Context, done chan struct{}, setup, draw func(c *Canvas), opts ...CallbackOption) *sync.WaitGroup {
	if setup == nil {
		log.Fatal("setup method is required")
	}
	if draw == nil {
		log.Fatal("draw method is required")
	}
	ctx, cancel := context.WithCancel(ctx)
	w, h := termSize()
	c := newCanvas(w, h)
	wg := sync.WaitGroup{}

	eventCallbacks := callbacks{}
	for _, opt := range opts {
		opt(&eventCallbacks)
	}

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
				case "exit":
					if done != nil {
						done <- struct{}{}
					}
					cancel()
					return
				}
			case event := <-inputEvents:
				switch e := event.(type) {
				case input.MouseMotionEvent:
					c.setMousePostion(e.X, e.Y)
					if eventCallbacks.onMouseMove != nil {
						eventCallbacks.onMouseMove(c, MouseEvent{
							X: c.MouseX,
							Y: c.MouseY,
						})
					}
				case input.MouseClickEvent:
					if eventCallbacks.onMouseClick != nil {
						c.setMousePostion(e.X, e.Y)
						eventCallbacks.onMouseClick(c, MouseEvent{
							X:      c.MouseX,
							Y:      c.MouseY,
							Button: e.Button.String(),
						})
					}
				case input.MouseReleaseEvent:
					if eventCallbacks.onMouseRelease != nil {
						c.setMousePostion(e.X, e.Y)
						eventCallbacks.onMouseRelease(c, MouseEvent{
							X:      c.MouseX,
							Y:      c.MouseY,
							Button: e.Button.String(),
						})
					}
				case input.MouseWheelEvent:
					if eventCallbacks.onMouseWheel != nil {
						c.setMousePostion(e.X, e.Y)
						eventCallbacks.onMouseWheel(c, MouseEvent{
							X:      c.MouseX,
							Y:      c.MouseY,
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
						if eventCallbacks.onKey != nil {
							eventCallbacks.onKey(c, KeyEvent{
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

func termSize() (int, int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal("can't read terminal size")
	}
	return w, h
}

func listenForInputEvents(ctx context.Context, wg *sync.WaitGroup) chan input.Event {
	inputEvents := make(chan input.Event, 2048)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(inputEvents)

		state, err := charmterm.MakeRaw(os.Stdin.Fd())
		if err != nil {
			log.Fatal(err)
		}
		defer charmterm.Restore(os.Stdin.Fd(), state) // nolint: errcheck

		reader, err := input.NewReader(os.Stdin, os.Getenv("TERM"), input.FlagMouseMode)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		readerEvents := make(chan []input.Event, 8)
		readerErrors := make(chan error, 8)
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				events, err := reader.ReadEvents()
				if err != nil {
					readerErrors <- err
					return
				}
				readerEvents <- events
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case events := <-readerEvents:
				for _, ev := range events {
					inputEvents <- ev
				}
			case err := <-readerErrors:
				log.Fatal(err)
			}
		}
	}()
	return inputEvents
}
