package runal

import (
	"os"

	ansitoimage "github.com/pavelpatrin/go-ansi-to-image"
)

// SaveCanvas exports the current canvas to a file (png).
func (c *Canvas) SaveCanvas(filename string) {
	c.save = true
	c.saveFilename = filename
	c.Redraw()
}

// SavedCanvasFont sets a custom font (tff) file used for rendering text characters
// in exported images generated via SaveCanvas().
func (c *Canvas) SavedCanvasFont(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	config := newCaptureConfig(c.Width, c.Height)
	config.MonoRegularFontBytes = file
	c.capture, _ = ansitoimage.NewConverter(config)
}

func newCapture(width, height int) *ansitoimage.Converter {
	imageCapture, _ := ansitoimage.NewConverter(newCaptureConfig(width, height))
	return imageCapture
}

func newCaptureConfig(width, height int) ansitoimage.Config {
	captureConfig := ansitoimage.DefaultConfig
	captureConfig.Padding = 0
	captureConfig.PageCols = width - 2
	captureConfig.PageRows = height
	return captureConfig
}

func (c *Canvas) exportCanvasToPNG(frame string) {
	err := c.capture.Parse(frame)
	if err != nil {
		return
	}
	img, _ := c.capture.ToPNG()
	_ = os.WriteFile(c.saveFilename, img, 0o644)
}
