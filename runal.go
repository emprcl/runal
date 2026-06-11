package runal

import (
	"context"
	"sync"

	"github.com/emprcl/runal/internal/canvas"
)

// Canvas is the main drawing surface.
type Canvas = canvas.Canvas

// Image represents a drawable image.
type Image = canvas.Image

// KeyEvent represents a keyboard event.
type KeyEvent = canvas.KeyEvent

// MouseEvent represents a mouse event.
type MouseEvent = canvas.MouseEvent

// CallbackOption is a functional option for Run and Start.
type CallbackOption = canvas.CallbackOption

// WithOnKey registers a callback for key events.
var WithOnKey = canvas.WithOnKey

// WithOnMouseMove registers a callback for mouse move events.
var WithOnMouseMove = canvas.WithOnMouseMove

// WithOnMouseClick registers a callback for mouse click events.
var WithOnMouseClick = canvas.WithOnMouseClick

// WithOnMouseRelease registers a callback for mouse release events.
var WithOnMouseRelease = canvas.WithOnMouseRelease

// WithOnMouseWheel registers a callback for mouse wheel events.
var WithOnMouseWheel = canvas.WithOnMouseWheel

// Run starts the runal event loop and blocks until the sketch exits.
func Run(ctx context.Context, setup, draw func(c *Canvas), opts ...CallbackOption) {
	canvas.Run(ctx, setup, draw, opts...)
}

// Start starts the runal event loop in a goroutine and returns a WaitGroup.
func Start(ctx context.Context, done chan struct{}, setup, draw func(c *Canvas), opts ...CallbackOption) *sync.WaitGroup {
	return canvas.Start(ctx, done, setup, draw, opts...)
}
