//go:build js

package canvas

// On the js/wasm build, image/video export is not supported: wiring it up would
// pull the ansitoimage + freetype + fogleman/gg + golang.org/x/image stack (and,
// transitively, lipgloss/termenv) into the binary for no runtime benefit, since
// there is no filesystem to write to. The capture backend is therefore a no-op
// and the SaveCanvasTo*/SavedCanvasFont* methods are inert stubs so sketches
// that call them still run.

type nopCapturer struct{}

type capturer = *nopCapturer

func newCapture(width, height int) capturer { return nil }

func (c *Canvas) captureResize(width, height int) {}

// SaveCanvasToPNG is a no-op on the web backend.
func (c *Canvas) SaveCanvasToPNG(filename string) {}

// SaveCanvasToGIF is a no-op on the web backend.
func (c *Canvas) SaveCanvasToGIF(filename string, duration int) {}

// SaveCanvasToMP4 is a no-op on the web backend.
func (c *Canvas) SaveCanvasToMP4(filename string, duration int) {}

// SavedCanvasFont is a no-op on the web backend.
func (c *Canvas) SavedCanvasFont(filename string) {}

// SavedCanvasFontSize is a no-op on the web backend.
func (c *Canvas) SavedCanvasFontSize(size int) {}
