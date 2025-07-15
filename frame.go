package runal

import "github.com/charmbracelet/lipgloss"

type Cell struct {
	Char       rune
	Foreground lipgloss.Color
	Background lipgloss.Color
}

type Frame [][]Cell

func newFrame(width, height int) Frame {
	buff := make([][]Cell, height)
	for i := range buff {
		buff[i] = make([]Cell, width)
	}
	return buff
}

func (c *Canvas) Get(x, y, w, h int) Image {
	if w <= 0 {
		w = 1
	}
	if h <= 0 {
		h = 1
	}
	frame := newFrame(w, h)
	for fy := 0; fy < h; fy++ {
		for fx := 0; fx < w; fx++ {
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
