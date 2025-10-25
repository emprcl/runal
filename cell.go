package runal

import (
	"github.com/charmbracelet/x/ansi"
)

type Cell struct {
	Char       string
	Foreground string
	Background string
}

func (cll Cell) write(c *Canvas, x, y, _, _ int) {
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
	padChar    rune
	foreground ansi.Color
	background ansi.Color
}

func (c cell) public() Cell {
	return Cell{
		Char:       string(c.char),
		Foreground: colorToString(c.foreground),
		Background: colorToString(c.background),
	}
}
