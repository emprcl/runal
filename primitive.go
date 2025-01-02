package runal

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

func (c *Canvas) Size(w, h int) {
	c.autoResize = false
	if c.widthPadding {
		c.resize(w*2, h)
	} else {
		c.resize(w, h)
	}
}

func (c *Canvas) Flush() {
	c.flush = true
}

func (c *Canvas) WidthPadding(char string) {
	previousValue := c.widthPadding
	c.widthPadding = true
	c.widthPaddingChar = rune(char[0])

	if !previousValue {
		c.resize(c.Width, c.Height)
	}
}

func (c *Canvas) DisableWidthPadding() {
	previousValue := c.widthPadding
	c.widthPadding = false

	if c.autoResize && previousValue {
		c.resize(c.Width*2, c.Height)
	} else if previousValue {
		c.resize(c.Width, c.Height)
	} else {
		c.resize(c.Width/2, c.Height)
	}
}

func (c *Canvas) Fps(fps int) {
	c.bus <- newFPSEvent(fps)
}

func (c *Canvas) Text(text string, x, y int) {
	if x < 0 || y < 0 || x > c.Width-1 || y > c.Height-1 {
		return
	}
	for i, r := range text {
		if x+i < len(c.buffer[y])-1 {
			c.buffer[y][x+i] = c.formatCell(r)
		} else if x+i == len(c.buffer[y])-1 {
			c.buffer[y][x+i] = c.style(string(r))
		}
	}
}

func (c *Canvas) Distance(x1, y1, x2, y2 int) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func (c *Canvas) Background(color string) {
	c.fillColor = lipgloss.Color(color)
}

func (c *Canvas) Foreground(color string) {
	c.strokeColor = lipgloss.Color(color)
}
