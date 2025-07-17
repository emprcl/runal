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
	Cell(x, y int) Cell
	writable
}

type writable interface {
	write(c *Canvas, x, y, w, h int)
}

type imageFile struct {
	file  image.Image
	frame frame
}

func (i *imageFile) Cell(x, y int) Cell {
	if i.frame.outOfBounds(x, y) {
		return Cell{}
	}
	return i.frame[y][x].public()
}

func (i *imageFile) write(c *Canvas, x, y, w, h int) {
	m := mosaic.New().Symbol(mosaic.Half)
	if w > 0 {
		m = m.Width(w)
	}

	if h > 0 {
		m = m.Height(h)
	}
	imageBuffer := m.RenderCells(i.file)
	i.frame = newFrame(w, h)
	c.toggleFill()
	for iy := range imageBuffer {
		for ix := range imageBuffer[iy] {
			if c.outOfBounds(x+ix, y+iy) {
				continue
			}
			cell := cell{
				char:       imageBuffer[iy][ix].Char,
				background: colorFromImage(imageBuffer[iy][ix].Background),
				foreground: colorFromImage(imageBuffer[iy][ix].Foreground),
			}
			c.write(cell, x+ix, y+iy, 2)
			i.frame[iy][ix] = cell
		}
	}
	c.toggleFill()
}

type imageFrame struct {
	frame frame
}

func (i *imageFrame) Cell(x, y int) Cell {
	if i.frame.outOfBounds(x, y) {
		return Cell{}
	}
	return i.frame[y][x].public()
}

func (i *imageFrame) write(c *Canvas, x, y, w, h int) {
	if h == 0 {
		_, h = i.frame.size()
	}
	if w == 0 {
		w, _ = i.frame.size()
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
