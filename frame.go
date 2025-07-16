package runal

import (
	"github.com/charmbracelet/lipgloss"
)

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
