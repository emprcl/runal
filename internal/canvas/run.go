package canvas

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/input"
)

const (
	defaultFPS = 30
	minFPS     = 1
	maxFPS     = 240
)

type callbacks struct {
	onKey          func(c *Canvas, e KeyEvent)
	onMouseMove    func(c *Canvas, e MouseEvent)
	onMouseClick   func(c *Canvas, e MouseEvent)
	onMouseRelease func(c *Canvas, e MouseEvent)
	onMouseWheel   func(c *Canvas, e MouseEvent)
}

type CallbackOption func(*callbacks)

func WithOnKey(onKey func(c *Canvas, e KeyEvent)) CallbackOption {
	return func(c *callbacks) {
		c.onKey = onKey
	}
}

func WithOnMouseMove(onMouseMove func(c *Canvas, e MouseEvent)) CallbackOption {
	return func(c *callbacks) {
		c.onMouseMove = onMouseMove
	}
}

func WithOnMouseClick(onMouseClick func(c *Canvas, e MouseEvent)) CallbackOption {
	return func(c *callbacks) {
		c.onMouseClick = onMouseClick
	}
}

func WithOnMouseRelease(onMouseRelease func(c *Canvas, e MouseEvent)) CallbackOption {
	return func(c *callbacks) {
		c.onMouseRelease = onMouseRelease
	}
}

func WithOnMouseWheel(onMouseWheel func(c *Canvas, e MouseEvent)) CallbackOption {
	return func(c *callbacks) {
		c.onMouseWheel = onMouseWheel
	}
}

func Run(ctx context.Context, setup, draw func(c *Canvas), opts ...CallbackOption) {
	wg, err := Start(ctx, nil, setup, draw, opts...)
	if err != nil {
		log.Error(err)
		return
	}
	wg.Wait()
}

func Start(ctx context.Context, done chan struct{}, setup, draw func(c *Canvas), opts ...CallbackOption) (*sync.WaitGroup, error) {
	if setup == nil {
		return nil, errors.New("setup method is required")
	}
	if draw == nil {
		return nil, errors.New("draw method is required")
	}

	w, h, err := tryTermSize()
	if err != nil {
		return nil, fmt.Errorf("can't read terminal size: %w", err)
	}
	c, err := newCanvas(w, h)
	if err != nil {
		return nil, err
	}
	in, err := newInputReader()
	if err != nil {
		return nil, fmt.Errorf("can't initialize input: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	eventCallbacks := callbacks{}
	for _, opt := range opts {
		opt(&eventCallbacks)
	}

	ticker := time.NewTicker(newFramerate(defaultFPS))

	exit := func() {
		ticker.Stop()
		in.close()
		resetCursorPosition(c.output)
		exitAltScreen(c.output)
		showCursor(c.output)
		disableMouse(c.output)
	}

	defer func() {
		if r := recover(); r != nil {
			exit()
			panic(r)
		}
	}()

	resize := listenForResize(ctx, &wg)
	in.listen(ctx, &wg)

	enterAltScreen(c.output)
	enableMouse(c.output)

	setup(c)
	render := func() {
		resetCursorPosition(c.output)
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
				eraseScreen(c.output)
				if w, h, err := tryTermSize(); err == nil {
					c.termWidth = w
					c.termHeight = h
					if c.autoResize {
						c.resize(w, h)
					}
				}
				render()
			case event := <-c.bus:
				switch event.name {
				case "fps":
					ticker.Reset(newFramerate(event.value))
				case "stop":
					ticker.Stop()
				case "start":
					ticker.Reset(newFramerate(c.fps))
				case "render":
					render()
				case "exit":
					if done != nil {
						done <- struct{}{}
					}
					cancel()
					return
				}
			case event, ok := <-in.events:
				// stdin is gone; shut down the same way ctrl+c does.
				if !ok {
					if done != nil {
						done <- struct{}{}
					}
					cancel()
					return
				}
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

	return &wg, nil
}

// notifyResize signals a pending resize without blocking.
func notifyResize(resize chan struct{}) {
	select {
	case resize <- struct{}{}:
	default:
	}
}

func newFramerate(fps int) time.Duration {
	return time.Second / time.Duration(clamp(fps, minFPS, maxFPS))
}
