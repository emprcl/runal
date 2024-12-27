package runal

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

func (c *Canvas) Width() int {
	return c.width
}

func (c *Canvas) Height() int {
	return c.height
}

func (c *Canvas) Framecount() uint64 {
	return c.framecount
}

func (c *Canvas) Text(text string, x, y int) {
	if x < 0 || y < 0 || x > c.width-1 || y > c.height-1 {
		return
	}
	for i, r := range text {
		if x+i < len(c.buffer[y])-1 {
			c.buffer[y][x+i] = c.style(string([]rune{r, padChar}))
		}
	}
}

func (c *Canvas) Dist(x1, y1, x2, y2 int) float64 {
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
