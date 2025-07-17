package runal

import (
	"github.com/charmbracelet/lipgloss"
)

type Cell struct {
	Char       string
	Foreground string
	Background string
}

func (cll Cell) write(c *Canvas, x, y, w, h int) {
	if c.outOfBounds(x, y) {
		return
	}
	c.write(cll.private(), x, y, 1)
}

func (cll Cell) private() cell {
	return cell{
		char:       []rune(cll.Char)[0],
		foreground: color(cll.Foreground),
		background: color(cll.Background),
	}
}

type cell struct {
	char       rune
	foreground lipgloss.Color
	background lipgloss.Color
}

func (c cell) public() Cell {
	return Cell{
		Char:       string(c.char),
		Foreground: string(c.foreground),
		Background: string(c.background),
	}
}
