package canvas

import "github.com/charmbracelet/x/ansi"

type state struct {
	strokeFg, strokeBg                      ansi.Color
	fillFg, fillBg                          ansi.Color
	backgroundFg, backgroundBg              ansi.Color
	strokeText, fillText, backgroundText    string
	strokeRunes, fillRunes, backgroundRunes []rune

	originX, originY int
	rotationAngle    float64
	scale            float64

	fill   bool
	stroke bool
}

// Push saves the current drawing state onto the stack. Every Push must be
// matched by a Pop; the stack is cleared at the end of each frame.
func (c *Canvas) Push() {
	c.stateStack = append(c.stateStack, c.currentState())
}

// Pop restores the most recently pushed drawing state.
func (c *Canvas) Pop() {
	if len(c.stateStack) == 0 {
		return
	}
	last := len(c.stateStack) - 1
	c.restoreState(c.stateStack[last])
	c.stateStack = c.stateStack[:last]
}

func (c *Canvas) currentState() state {
	return state{
		strokeFg:        c.strokeFg,
		strokeBg:        c.strokeBg,
		fillFg:          c.fillFg,
		fillBg:          c.fillBg,
		backgroundFg:    c.backgroundFg,
		backgroundBg:    c.backgroundBg,
		strokeText:      c.strokeText,
		strokeRunes:     c.strokeRunes,
		fillText:        c.fillText,
		fillRunes:       c.fillRunes,
		backgroundText:  c.backgroundText,
		backgroundRunes: c.backgroundRunes,
		originX:         c.originX,
		originY:         c.originY,
		rotationAngle:   c.rotationAngle,
		scale:           c.scale,
		fill:            c.fill,
		stroke:          c.stroke,
	}
}

func (c *Canvas) restoreState(s state) {
	c.strokeFg = s.strokeFg
	c.strokeBg = s.strokeBg
	c.fillFg = s.fillFg
	c.fillBg = s.fillBg
	c.backgroundFg = s.backgroundFg
	c.backgroundBg = s.backgroundBg
	c.strokeText = s.strokeText
	c.strokeRunes = s.strokeRunes
	c.fillText = s.fillText
	c.fillRunes = s.fillRunes
	c.backgroundText = s.backgroundText
	c.backgroundRunes = s.backgroundRunes
	c.originX = s.originX
	c.originY = s.originY
	c.rotationAngle = s.rotationAngle
	c.scale = s.scale
	c.fill = s.fill
	c.stroke = s.stroke
}
