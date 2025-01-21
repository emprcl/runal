package runal

import (
	"github.com/charmbracelet/lipgloss"
)

func (c *Canvas) Background(color string) {
	c.backgroundColor = lipgloss.Color(color)
}

func (c *Canvas) Fill(color string) {
	c.fillColor = lipgloss.Color(color)
}

func (c *Canvas) Color(color string) {
	c.textColor = lipgloss.Color(color)
}
