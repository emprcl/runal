package runal

import (
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
