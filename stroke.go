package runal

// Background sets the background character and colors for the entire canvas.
func (c *Canvas) Background(text, fg, bg string) {
	c.BackgroundText(text)
	c.BackgroundBg(bg)
	c.BackgroundFg(fg)
}

// BackgroundText sets the character used for the background fill.
func (c *Canvas) BackgroundText(text string) {
	c.backgroundIndex = 0
	if len(text) == 0 {
		c.backgroundText = defaultBackgroundText
		return
	}
	c.backgroundText = text
}

// BackgroundFg sets the foreground (text) color used by the background fill.
func (c *Canvas) BackgroundFg(fg string) {
	c.backgroundFg = color(fg)
}

// BackgroundBg sets the background color used by the background fill.
func (c *Canvas) BackgroundBg(bg string) {
	c.backgroundBg = color(bg)
}

// Fill sets the fill character and its foreground and background colors.
func (c *Canvas) Fill(text, fg, bg string) {
	c.FillText(text)
	c.FillBg(bg)
	c.FillFg(fg)
}

// FillText sets the character used for fill operations.
func (c *Canvas) FillText(text string) {
	c.fill = true
	if len(text) == 0 {
		c.fillText = defaultFillText
		return
	}
	c.fillText = text
}

// FillFg sets the foreground color used for fill operations.
func (c *Canvas) FillFg(fg string) {
	c.fill = true
	c.fillFg = color(fg)
}

// FillBg sets the background color used for fill operations.
func (c *Canvas) FillBg(bg string) {
	c.fill = true
	c.fillBg = color(bg)
}

// Stroke sets the stroke (outline) character and colors.
func (c *Canvas) Stroke(text, fg, bg string) {
	c.StrokeText(text)
	c.StrokeBg(bg)
	c.StrokeFg(fg)
}

// StrokeText sets the character used for strokes.
func (c *Canvas) StrokeText(text string) {
	c.stroke = true
	c.strokeIndex = 0
	if len(text) == 0 {
		c.strokeText = defaultStrokeText
	}
	c.strokeText = text
}

// StrokeFg sets the foreground color for strokes.
func (c *Canvas) StrokeFg(fg string) {
	c.stroke = true
	c.strokeFg = color(fg)
}

// StrokeBg sets the background color for strokes.
func (c *Canvas) StrokeBg(bg string) {
	c.stroke = true
	c.strokeBg = color(bg)
}

// NoStroke disables stroke for subsequent shapes.
func (c *Canvas) NoStroke() {
	c.stroke = false
}

// NoFill disables fill for subsequent shapes.
func (c *Canvas) NoFill() {
	c.fill = false
}
