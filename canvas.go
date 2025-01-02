package runal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	defaultPaddingChar = ' '
)

type Canvas struct {
	buffer                 buffer
	strokeColor, fillColor lipgloss.Color
	bus                    chan event

	Width, Height int
	Framecount    int

	termWidth, termHeight int
	widthPaddingChar      rune
	widthPadding          bool
	flush                 bool
	autoResize            bool
}

func newCanvas(width, height int) *Canvas {
	return &Canvas{
		Width:            width,
		Height:           height,
		bus:              make(chan event),
		termWidth:        width,
		termHeight:       height,
		widthPaddingChar: defaultPaddingChar,
		widthPadding:     false,
		buffer:           newBuffer(width, height),
		strokeColor:      lipgloss.Color("#ffffff"),
		fillColor:        lipgloss.Color("#000000"),
		autoResize:       true,
	}
}

func (c *Canvas) render() {
	output := ""
	for y := range c.buffer {
		line := ""
		for x := range c.buffer[y] {
			add := ""
			if c.buffer[y][x] == "" {
				add = c.formatCell(defaultPaddingChar)
			} else {
				add = c.buffer[y][x]
			}
			if lipgloss.Width(line+add) < c.termWidth {
				line += add
			}
			if c.flush {
				c.buffer[y][x] = ""
			}
		}
		output += forcePadding(line, c.termWidth, ' ')
		if y < len(c.buffer)-1 {
			output += "\n"
		}
	}
	c.Framecount++
	c.flush = false
	fmt.Print(output)
}

func (c *Canvas) resize(width, height int) {
	newWidth := c.widthWithPadding(width)
	newHeight := height
	newBuffer := newBuffer(newWidth, newHeight)

	minWidth := c.Width
	if newWidth < c.Width {
		minWidth = newWidth
	}

	minHeight := c.Height
	if newHeight < c.Height {
		minHeight = newHeight
	}

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			newBuffer[y][x] = c.buffer[y][x]
		}
	}

	c.Width = newWidth
	c.Height = newHeight
	c.buffer = newBuffer
}

func (c *Canvas) style(str string) string {
	return lipgloss.NewStyle().
		Background(c.fillColor).
		Foreground(c.strokeColor).
		Render(str)
}

func (c *Canvas) formatCell(char rune) string {
	if c.widthPadding {
		return c.style(string([]rune{char, c.widthPaddingChar}))
	}
	return c.style(string(char))
}

func (c *Canvas) widthWithPadding(w int) int {
	if c.widthPadding {
		return w / 2
	}
	return w
}
