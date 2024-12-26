package smbols

import (
	"fmt"
	"math"
	"smbols/internal/util"

	"github.com/charmbracelet/lipgloss"
)

const (
	padChar rune = ' '
)

type Canvas struct {
	width, height          int
	termWidth, termHeight  int
	buffer                 buffer
	strokeColor, fillColor lipgloss.Color
}

func newCanvas(width, height int) *Canvas {
	return &Canvas{
		width:       width / 2,
		height:      height,
		termWidth:   width,
		termHeight:  height,
		buffer:      newBuffer(width, height),
		strokeColor: lipgloss.Color("#ffffff"),
		fillColor:   lipgloss.Color("#000000"),
	}
}

func (c *Canvas) render() {
	output := ""
	for y := range c.buffer {
		line := ""
		for x := range c.buffer[y] {
			add := ""
			if c.buffer[y][x] == "" {
				add = c.Style("  ")
			} else {
				add = c.buffer[y][x]
			}
			if lipgloss.Width(line+add) <= c.termWidth {
				line += add
			}
			c.buffer[y][x] = ""
		}
		output += util.ForcePadding(line, c.termWidth, ' ')
	}
	fmt.Print(output)
}

func (c *Canvas) resize(width, height int) {
	newWidth := width / 2
	newHeight := height
	newBuffer := newBuffer(newWidth, newHeight)

	minWidth := c.width
	if newWidth < c.width {
		minWidth = newWidth
	}

	minHeight := c.height
	if newHeight < c.height {
		minHeight = newHeight
	}

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			newBuffer[y][x] = c.buffer[y][x]
		}
	}

	c.width = newWidth
	c.height = newHeight
	c.termWidth = width
	c.termHeight = height
	c.buffer = newBuffer
}

func (c *Canvas) Width() int {
	return c.width
}

func (c *Canvas) Height() int {
	return c.height
}

func (c *Canvas) Char(char rune, x, y int) {
	c.buffer[y][x] = c.Style(string([]rune{char, padChar}))
}

func (c *Canvas) Style(str string) string {
	return lipgloss.NewStyle().
		Background(c.fillColor).
		Foreground(c.strokeColor).
		Render(str)
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
