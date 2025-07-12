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

type Image struct {
	file image.Image
}

func (c *Canvas) LoadImage(path string) *Image {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Errorf("can't load image: %v", err)
		c.DisableRendering()
		return nil
	}
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

	return &Image{
		file: img,
	}
}

func (c *Canvas) Image(img *Image, x, y, w, h int) {
	if img == nil {
		log.Errorf("can't load empty image")
		c.DisableRendering()
		return
	}
	m := mosaic.New().Width(w).Height(h).Symbol(mosaic.Quarter)
	imageBuffer := m.RenderCells(img.file)

	c.toggleFill()

	for iy := range imageBuffer {
		for ix := range imageBuffer[iy] {
			if c.outOfBounds(x+ix, y+iy) {
				continue
			}
			c.write(c.formatStringCell(imageBuffer[iy][ix]), x+ix, y+iy, 2+int(c.scale))
		}
	}

	c.toggleFill()
}
