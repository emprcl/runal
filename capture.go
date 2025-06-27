package runal

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"os"

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
	totalFrames := duration * c.fps
	c.frames = make([]image.Image, 0, totalFrames)
	c.saveFilename = filename
}

// SaveCanvasToMP4 exports the canvas to a mp4 (h264) video
// for a given duration (in seconds).
func (c *Canvas) SaveCanvasToMP4(filename string, duration int) {
	if !checkFFMPEG() {
		fmt.Println("ffmpeg is not installed")
		c.DisableRendering()
	}
	if c.frames != nil {
		return
	}
	c.videoFormat = videoFormatMp4
	totalFrames := duration * c.fps
	c.frames = make([]image.Image, 0, totalFrames)
	c.saveFilename = filename
}

// SavedCanvasFont sets a custom font (tff) file used for rendering text characters
// in exported images generated via SaveCanvasTo...().
func (c *Canvas) SavedCanvasFont(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	config := c.capture.Config()
	config.MonoRegularFontBytes = file
	c.capture, _ = ansitoimage.NewConverter(config)
}

// SavedCanvasFontSize sets the font size used for rendering text characters
// in exported images generated via SaveCanvas().
func (c *Canvas) SavedCanvasFontSize(size int) {
	config := c.capture.Config()
	config.MonoRegularFontPoints = float64(size)
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

func (c *Canvas) captureResize(width, height int) {
	config := c.capture.Config()
	config.PageCols = width - 2
	config.PageRows = height
	c.capture, _ = ansitoimage.NewConverter(config)
}

func (c *Canvas) generateFrame(frame string) {
	err := c.capture.Parse(frame)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Canvas) exportCanvasToPNG(frame string) {
	fmt.Println("Saving png...")
	c.generateFrame(frame)
	img, _ := c.capture.ToPNG()
	err := os.WriteFile(c.saveFilename, img, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Canvas) exportFramesToGIF() error {
	outGif := &gif.GIF{}

	delay := 100 / c.fps

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
			log.Fatal(err)
		}
		defer f.Close()

		if err := png.Encode(f, img); err != nil {
			log.Fatal(err)
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
			fmt.Println("Saving GIF...")
			err = c.exportFramesToGIF()
		case videoFormatMp4:
			fmt.Println("Saving MP4...")
			err = c.exportFramesToMP4()
		}

		if err != nil {
			log.Fatal(err)
		}
		c.frames = nil
		return
	}
	c.generateFrame(output)
	frame := image.NewRGBA(c.capture.Image().Bounds())
	copy(frame.Pix, c.capture.Image().(*image.RGBA).Pix)
	c.frames = append(c.frames, frame)
}
