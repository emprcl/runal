package canvas

import (
	"time"
)

const (
	defaultFPS = 30
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

func newFramerate(fps int) time.Duration {
	return time.Second / time.Duration(fps)
}
