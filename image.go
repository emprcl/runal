package runal

import (
	"fmt"
	"os"
	"time"

	ansitoimage "github.com/pavelpatrin/go-ansi-to-image"
)

const screenshotFilename = "screenshot_%s.png"

func (c *Canvas) SaveCanvas() {
	c.save = true
	c.Redraw()
}

func newCapture(width, height int) *ansitoimage.Converter {
	captureConfig := ansitoimage.DefaultConfig
	captureConfig.Padding = 0
	captureConfig.PageCols = width - 1
	captureConfig.PageRows = height
	imageCapture, _ := ansitoimage.NewConverter(captureConfig)
	return imageCapture
}

func (c *Canvas) exportCanvasToPNG(frame string) {
	c.capture.Parse(frame)
	img, _ := c.capture.ToPNG()
	os.WriteFile(fmt.Sprintf(screenshotFilename, time.Now().Local().Format(time.RFC3339Nano)), img, 0644)
}
