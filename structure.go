package runal

import "github.com/charmbracelet/lipgloss"

type saveState struct {
	strokeFg, strokeBg                   lipgloss.Color
	fillFg, fillBg                       lipgloss.Color
	backgroundFg, backgroundBg           lipgloss.Color
	strokeText, fillText, backgroundText string

	originX, originY int
	rotationAngle    float64
	scale            float64

	fill bool
}

func (c *Canvas) Push() {
	c.saveState = &saveState{
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
	}
}

func (c *Canvas) Pop() {
	if c.saveState == nil {
		return
	}

	c.strokeFg = c.saveState.strokeFg
	c.strokeBg = c.saveState.strokeBg
	c.fillFg = c.saveState.fillFg
	c.fillBg = c.saveState.fillBg
	c.backgroundFg = c.saveState.backgroundFg
	c.backgroundBg = c.saveState.backgroundBg
	c.strokeText = c.saveState.strokeText
	c.fillText = c.saveState.fillText
	c.backgroundText = c.saveState.backgroundText
	c.originX = c.saveState.originX
	c.originY = c.saveState.originY
	c.rotationAngle = c.saveState.rotationAngle
	c.scale = c.saveState.scale
	c.fill = c.saveState.fill

	c.saveState = nil
}
