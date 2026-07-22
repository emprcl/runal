package canvas

import (
	"fmt"

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
		char:       firstRune(cll.Char, defaultPaddingRune),
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
	char := ""
	if c.char != 0 {
		char = string(c.char)
	}
	return Cell{
		Char:       char,
		Foreground: colorToString(c.foreground),
		Background: colorToString(c.background),
	}
}

func colorToString(c ansi.Color) string {
	if c == nil {
		return ""
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
