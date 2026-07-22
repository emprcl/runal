package canvas

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"os"

	"github.com/charmbracelet/log"

	ansitoimage "github.com/xaviergodart/go-ansi-to-image"
)

// SaveCanvasToPNG exports the canvas to a png image file.
func (c *Canvas) SaveCanvasToPNG(filename string) {
	c.save = true
	c.saveFilename = filename
	c.Redraw()
}

// SaveCanvasToGIF exports the canvas to an animated gif
// for a given duration (in seconds).
func (c *Canvas) SaveCanvasToGIF(filename string, duration int) {
	if c.frames != nil {
		return
	}
	c.videoFormat = videoFormatGif
	c.frames = make([]image.Image, 0, captureFrames(duration, c.fps))
	c.saveFilename = filename
}

// SaveCanvasToMP4 exports the canvas to a mp4 (h264) video
// for a given duration (in seconds).
// Depends on ffmpeg.
func (c *Canvas) SaveCanvasToMP4(filename string, duration int) {
	if !checkFFMPEG() {
		log.Error("Can't use SaveCanvasToMP4(). ffmpeg is not installed.")
		c.DisableRendering()
		return
	}
	if c.frames != nil {
		return
	}
	c.videoFormat = videoFormatMp4
	c.frames = make([]image.Image, 0, captureFrames(duration, c.fps))
	c.saveFilename = filename
}

func captureFrames(duration, fps int) int {
	return max(duration, 1) * clamp(fps, minFPS, maxFPS)
}

// SavedCanvasFont sets a custom font (tff) file used for rendering text characters
// in exported images generated via SaveCanvasTo...().
func (c *Canvas) SavedCanvasFont(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Errorf("can't read font file: %v", err)
		c.DisableRendering()
		return
	}
	config := c.capture.Config()
	config.MonoRegularFontBytes = file
	c.setCapture(config)
}

// SavedCanvasFontSize sets the font size used for rendering text characters
// in exported images generated via SaveCanvas().
func (c *Canvas) SavedCanvasFontSize(size int) {
	config := c.capture.Config()
	config.MonoRegularFontPoints = float64(size)
	c.setCapture(config)
}

func (c *Canvas) setCapture(config ansitoimage.Config) {
	capture, err := ansitoimage.NewConverter(config)
	if err != nil {
		log.Errorf("can't configure canvas export: %v", err)
		c.DisableRendering()
		return
	}
	c.capture = capture
}

func newCapture(width, height int) (*ansitoimage.Converter, error) {
	return ansitoimage.NewConverter(newCaptureConfig(width, height))
}

func newCaptureConfig(width, height int) ansitoimage.Config {
	captureConfig := ansitoimage.DefaultConfig
	captureConfig.Padding = 0
	captureConfig.PageCols = max(width-2, 1)
	captureConfig.PageRows = max(height, 1)
	return captureConfig
}

func (c *Canvas) captureResize(width, height int) {
	config := c.capture.Config()
	config.PageCols = max(width-2, 1)
	config.PageRows = max(height, 1)
	c.setCapture(config)
}

func (c *Canvas) generateFrame(frame string) error {
	if err := c.capture.Parse(frame); err != nil {
		return fmt.Errorf("can't generate frame: %w", err)
	}
	return nil
}

func (c *Canvas) captureFailed(err error) {
	log.Error(err)
	c.DisableRendering()
}

func (c *Canvas) exportCanvasToPNG(frame string) {
	fmt.Fprintln(c.output, "Saving png...")
	if err := c.generateFrame(frame); err != nil {
		c.captureFailed(err)
		return
	}
	img, err := c.capture.ToPNG()
	if err != nil {
		c.captureFailed(fmt.Errorf("can't encode png: %w", err))
		return
	}
	if err := os.WriteFile(c.saveFilename, img, 0o644); err != nil {
		c.captureFailed(fmt.Errorf("can't create png file: %w", err))
	}
}

func gifDelay(fps int) int {
	return max(int(math.Round(100/float64(clamp(fps, minFPS, maxFPS)))), 1)
}

func (c *Canvas) exportFramesToGIF() error {
	outGif := &gif.GIF{}

	delay := gifDelay(c.fps)

	for _, img := range c.frames {
		bounds := img.Bounds()
		paletted := image.NewPaletted(bounds, palette.Plan9)
		draw.FloydSteinberg.Draw(paletted, bounds, img, image.Point{})
		outGif.Image = append(outGif.Image, paletted)
		outGif.Delay = append(outGif.Delay, delay)
	}

	file, err := os.Create(c.saveFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	return gif.EncodeAll(file, outGif)
}

func (c *Canvas) exportFramesToMP4() error {
	dir := randomDir()
	defer os.RemoveAll(dir)

	for i, img := range c.frames {
		f, err := os.Create(fmt.Sprintf("%s/frame_%d.png", dir, i))
		if err != nil {
			return fmt.Errorf("can't create frame file: %w", err)
		}

		if err := png.Encode(f, img); err != nil {
			f.Close()
			return fmt.Errorf("can't encode frame: %w", err)
		}
		f.Close()
	}

	return framesToMP4Videos(c.fps, fmt.Sprintf("%s/frame_%%d.png", dir), c.saveFilename)
}

func (c *Canvas) recordFrame(output string) {
	if len(c.frames) >= cap(c.frames) {
		var err error
		switch c.videoFormat {
		case videoFormatGif:
			fmt.Fprintln(c.output, "Saving GIF...")
			err = c.exportFramesToGIF()
		case videoFormatMp4:
			fmt.Fprintln(c.output, "Saving MP4...")
			err = c.exportFramesToMP4()
		}

		c.frames = nil
		if err != nil {
			c.captureFailed(fmt.Errorf("can't export frames: %w", err))
		}
		return
	}
	if err := c.generateFrame(output); err != nil {
		c.captureFailed(err)
		return
	}
	frame := image.NewRGBA(c.capture.Image().Bounds())
	copy(frame.Pix, c.capture.Image().(*image.RGBA).Pix)
	c.frames = append(c.frames, frame)
}
