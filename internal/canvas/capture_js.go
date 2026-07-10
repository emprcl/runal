//go:build js

package canvas

// On the js/wasm build, export is done entirely in the browser (see web/runal.js
// and web/gif.js) rather than through the ansitoimage + freetype + fogleman/gg +
// golang.org/x/image stack, which would bloat the wasm binary. These methods
// just signal the JS renderer, which owns the actual <canvas> and encodes/
// downloads the result. There is no local filesystem, so `filename` is used as
// the download name.

type nopCapturer struct{}

type capturer = *nopCapturer

func newCapture(width, height int) capturer { return nil }

func (c *Canvas) captureResize(width, height int) {}

func (c *Canvas) displayReady() bool {
	return display != nil && !display.runal.IsUndefined()
}

// SaveCanvasToPNG downloads the current frame as a PNG.
func (c *Canvas) SaveCanvasToPNG(filename string) {
	if !c.displayReady() {
		return
	}
	display.runal.Call("savePNG", display.el, filename)
}

// SaveCanvasToGIF records `duration` seconds of frames and downloads an
// animated GIF.
func (c *Canvas) SaveCanvasToGIF(filename string, duration int) {
	if !c.displayReady() {
		return
	}
	display.runal.Call("recordGIF", display.el, filename, duration, c.fps)
}

// SaveCanvasToMP4 records `duration` seconds of the canvas and downloads a
// video (mp4 where the browser supports it, otherwise webm).
func (c *Canvas) SaveCanvasToMP4(filename string, duration int) {
	if !c.displayReady() {
		return
	}
	display.runal.Call("recordVideo", display.el, filename, duration)
}

// SavedCanvasFont is a no-op on the web backend: the exported image uses the
// canvas font, not a separate rasterization font.
func (c *Canvas) SavedCanvasFont(filename string) {}

// SavedCanvasFontSize is a no-op on the web backend.
func (c *Canvas) SavedCanvasFontSize(size int) {}
