//go:build js

package canvas

// On the js/wasm build the ansitoimage + freetype + fogleman/gg +
// golang.org/x/image export stack is not linked (it would bloat the binary and
// there is no filesystem). Image/video export is handled entirely in the
// browser by the JS proxy (web/runal.js, web/gif.js), which owns the <canvas>.
// The Go-side SaveCanvasTo*/SavedCanvasFont* methods are therefore inert; on
// the web a sketch's calls are handled by the JS proxy's savePNG/recordGIF/
// recordVideo directly, not these.

type nopCapturer struct{}

type capturer = *nopCapturer

func newCapture(width, height int) capturer { return nil }

func (c *Canvas) captureResize(width, height int) {}

func (c *Canvas) SaveCanvasToPNG(filename string)               {}
func (c *Canvas) SaveCanvasToGIF(filename string, duration int) {}
func (c *Canvas) SaveCanvasToMP4(filename string, duration int) {}
func (c *Canvas) SavedCanvasFont(filename string)               {}
func (c *Canvas) SavedCanvasFontSize(size int)                  {}
