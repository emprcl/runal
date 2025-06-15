package runal

import (
	"fmt"
	"os"
	"time"

	ansitoimage "github.com/pavelpatrin/go-ansi-to-image"
)

const canvasFilename = "canvas_%s.png"

func (c *Canvas) SaveCanvas() {
	c.save = true
	c.Redraw()
}

func (c *Canvas) CanvasFont(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	config := newCaptureConfig(c.termWidth, c.termHeight)
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
	captureConfig.PageCols = width - 1
	captureConfig.PageRows = height
	return captureConfig
}

func (c *Canvas) exportCanvasToPNG(frame string) {
	c.capture.Parse(frame)
	img, _ := c.capture.ToPNG()
	os.WriteFile(fmt.Sprintf(canvasFilename, time.Now().Local().Format(time.RFC3339Nano)), img, 0644)
}
