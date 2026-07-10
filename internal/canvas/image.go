package canvas

type Image interface {
	Cell(x, y int) Cell
	writable
}

type writable interface {
	write(c *Canvas, x, y, w, h int)
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

func (c *Canvas) Image(img Image, x, y, w, h int) {
	if img == nil {
		logErrorf("can't load empty image")
		c.DisableRendering()
		return
	}
	img.write(c, x, y, w, h)
}
