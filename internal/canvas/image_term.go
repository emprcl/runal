//go:build !js

package canvas

import (
	"image"
	col "image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/charmbracelet/x/ansi"
	"golang.org/x/image/webp"

	"github.com/emprcl/runal/internal/mosaic"
)

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

func (c *Canvas) LoadImage(path string) Image {
	f, err := os.Open(path)
	if err != nil {
		logErrorf("can't load image: %v", err)
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
		logErrorf("can't load image: %v", err)
		c.DisableRendering()
		return nil
	}

	return &imageFile{
		file: img,
	}
}

func colorFromImage(c col.Color) ansi.Color {
	rgba := col.RGBAModel.Convert(c).(col.RGBA)
	return ansi.RGBColor{
		R: rgba.R,
		G: rgba.G,
		B: rgba.B,
	}
}
