package runal

type frame [][]cell

func (f frame) size() (int, int) {
	return len(f[0]), len(f)
}

func (f frame) outOfBounds(x, y int) bool {
	w, h := f.size()
	return x < 0 || y < 0 || x >= w || y >= h
}

func newFrame(width, height int) frame {
	buff := make([][]cell, height)
	for i := range buff {
		buff[i] = make([]cell, width)
	}
	return buff
}

func (c *Canvas) Get(x, y, w, h int) Image {
	frame := newFrame(w, h)
	for fy := range frame {
		for fx := range frame[fy] {
			if c.outOfBounds(x+fx, y+fy) {
				continue
			}
			frame[fy][fx] = c.buffer[y+fy][x+fx]
		}
	}
	return &imageFrame{
		frame: frame,
	}
}

func (c *Canvas) Set(x, y int, cells writable) {
	cells.write(c, x, y, 0, 0)
}
