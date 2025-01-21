package runal

import (
	"github.com/charmbracelet/lipgloss"
)

func (c *Canvas) Background(text, fg, bg string) {
	c.backgroundText = text
	c.backgroundIndex = 0
	c.backgroundBg = lipgloss.Color(bg)
	c.backgroundFg = lipgloss.Color(fg)

	if len(text) == 0 {
		c.backgroundText = defaultBackgroundText
	}
}

func (c *Canvas) Fill(text, fg, bg string) {
	c.fill = true
	c.fillText = text
	c.fillBg = lipgloss.Color(bg)
	c.fillFg = lipgloss.Color(fg)

	if len(text) == 0 {
		c.fillText = defaultFillText
	}
}

func (c *Canvas) NoFill() {
	c.fill = false
}

func (c *Canvas) Stroke(text, fg, bg string) {
	c.strokeText = text
	c.strokeBg = lipgloss.Color(bg)
	c.strokeFg = lipgloss.Color(fg)

	if len(text) == 0 {
		c.strokeText = defaultStrokeText
	}
}
