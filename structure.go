package runal

import "github.com/charmbracelet/x/ansi"

type state struct {
	strokeFg, strokeBg                   ansi.Color
	fillFg, fillBg                       ansi.Color
	backgroundFg, backgroundBg           ansi.Color
	strokeText, fillText, backgroundText string

	originX, originY int
	rotationAngle    float64
	scale            float64

	fill   bool
	stroke bool
}

// Push saves the current transformation state.
func (c *Canvas) Push() {
	c.state = &state{
		strokeFg:       c.strokeFg,
		strokeBg:       c.strokeBg,
		fillFg:         c.fillFg,
		fillBg:         c.fillBg,
		backgroundFg:   c.backgroundFg,
		backgroundBg:   c.backgroundBg,
		strokeText:     c.strokeText,
		fillText:       c.fillText,
		backgroundText: c.backgroundText,
		originX:        c.originX,
		originY:        c.originY,
		rotationAngle:  c.rotationAngle,
		scale:          c.scale,
		fill:           c.fill,
		stroke:         c.stroke,
	}
}

// Pop restores the previous transformation state.
func (c *Canvas) Pop() {
	if c.state == nil {
		return
	}

	c.strokeFg = c.state.strokeFg
	c.strokeBg = c.state.strokeBg
	c.fillFg = c.state.fillFg
	c.fillBg = c.state.fillBg
	c.backgroundFg = c.state.backgroundFg
	c.backgroundBg = c.state.backgroundBg
	c.strokeText = c.state.strokeText
	c.fillText = c.state.fillText
	c.backgroundText = c.state.backgroundText
	c.originX = c.state.originX
	c.originY = c.state.originY
	c.rotationAngle = c.state.rotationAngle
	c.scale = c.state.scale
	c.fill = c.state.fill
	c.stroke = c.state.stroke

	c.state = nil
}
