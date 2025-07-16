package runal

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"golang.org/x/image/webp"

	"github.com/emprcl/runal/pkg/mosaic"
)

type Image interface {
	write(c *Canvas, x, y, w, h int)
}

type imageFile struct {
	file image.Image
}

func (i *imageFile) write(c *Canvas, x, y, w, h int) {
	m := mosaic.New()
	if w > 0 {
		m = m.Width(w)
	}

	if h > 0 {
		m = m.Height(h)
	}
	m = m.Symbol(mosaic.Quarter)
	imageBuffer := m.RenderCells(i.file)
	c.toggleFill()
	for iy := range imageBuffer {
		for ix := range imageBuffer[iy] {
			if c.outOfBounds(x+ix, y+iy) {
				continue
			}
			c.write(Cell{
				Char:       imageBuffer[iy][ix].Char,
				Background: colorFromImage(imageBuffer[iy][ix].Background),
				Foreground: colorFromImage(imageBuffer[iy][ix].Foreground),
			}, x+ix, y+iy, 2)
		}
	}
	c.toggleFill()
}

type imageFrame struct {
	frame Frame
}

func (i *imageFrame) write(c *Canvas, x, y, w, h int) {
	if h == 0 {
		_, h = i.frame.Size()
	}
	if w == 0 {
		w, _ = i.frame.Size()
	}
	for iy := range i.frame {
		for ix := range i.frame[iy] {
			if c.outOfBounds(x+ix, y+iy) ||
				x+ix >= x+w || y+iy >= y+h {
				continue
			}
			c.write(i.frame[iy][ix], x+ix, y+iy, 1)
		}
	}
}

func (c *Canvas) LoadImage(path string) Image {
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("can't load image: %v", err)
		c.DisableRendering()
		return nil
	}
	defer f.Close()
	var img image.Image
	switch filepath.Ext(path) {
	case ".jpg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	case ".webp":
		img, err = webp.Decode(f)
	}

	if err != nil {
		log.Errorf("can't load image: %v", err)
		c.DisableRendering()
		return nil
	}

	return &imageFile{
		file: img,
	}
}

func (c *Canvas) Image(img Image, x, y, w, h int) {
	if img == nil {
		log.Errorf("can't load empty image")
		c.DisableRendering()
		return
	}
	img.write(c, x, y, w, h)
}
