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

func (c *Canvas) loadImage(path string) *Image {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Errorf("can't load image: %v", err)
		c.DisableRendering()
		return nil
	}
	var img image.Image
	switch filepath.Ext(path) {
	case "jpg":
		img, err = jpeg.Decode(f)
	case "png":
		img, err = png.Decode(f)
	case "webp":
		img, err = webp.Decode(f)
	}

	if err != nil {
		log.Errorf("can't load image: %v", err)
		c.DisableRendering()
	}

	return &Image{
		file: img,
	}
}

func (c *Canvas) Image(img Image, x, y int) {
	m := mosaic.New().Width(80).Height(40)
	imageBuffer := m.RenderCells(img.file)

	for iy := range imageBuffer {
		for ix := range imageBuffer[x] {
			if c.outOfBounds(x+ix, y+iy) {
				continue
			}
			c.buffer[y+iy][x+ix] = imageBuffer[y][x]
		}
	}
}
